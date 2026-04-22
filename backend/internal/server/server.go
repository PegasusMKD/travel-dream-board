package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/PegasusMKD/travel-dream-board/internal/accomodations"
	"github.com/PegasusMKD/travel-dream-board/internal/activities"
	"github.com/PegasusMKD/travel-dream-board/internal/auth"
	"github.com/PegasusMKD/travel-dream-board/internal/boards"
	"github.com/PegasusMKD/travel-dream-board/internal/comments"
	"github.com/PegasusMKD/travel-dream-board/internal/config"
	"github.com/PegasusMKD/travel-dream-board/internal/database"
	"github.com/PegasusMKD/travel-dream-board/internal/db"
	"github.com/PegasusMKD/travel-dream-board/internal/middleware"
	scrapeaudit "github.com/PegasusMKD/travel-dream-board/internal/scrape_audit"
	scrapeprocess "github.com/PegasusMKD/travel-dream-board/internal/scrape_process"
	"github.com/PegasusMKD/travel-dream-board/internal/sharetokens"
	"github.com/PegasusMKD/travel-dream-board/internal/transport"
	"github.com/PegasusMKD/travel-dream-board/internal/votes"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type GinServer struct {
	addr   string
	engine *gin.Engine
	server *http.Server
}

func (srv *GinServer) Run() {
	log.Info("Starting server on", "addr", srv.addr)
	if err := srv.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Error("Server failed", "error", err)
	}
}

func (srv *GinServer) Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	log.Info("Shutting down server...")
	if err := srv.server.Shutdown(ctx); err != nil {
		log.Error("Server forced to shutdown", "error", err)
	}

	log.Info("Stopping scheduler...")
}

func NewServer() *GinServer {
	cfg, err := config.Load()
	if err != nil {
		log.Error("Failed loading config for github.com/PegasusMKD/travel-dream-board!", "error", err)
	}

	addr := fmt.Sprintf(":%s", cfg.Port)

	queries, err := setupSqlc(cfg)
	if err != nil {
		log.Error("Failed to initialize database - cannot start server", "error", err)
		return nil
	}

	router := gin.Default()
	router.RemoveExtraSlash = true

	srv := &GinServer{
		engine: router,
		addr:   addr,
	}

	srv.server = &http.Server{
		Addr:    srv.addr,
		Handler: srv.engine,
	}

	srv.setupMiddleware(router)
	srv.setupRoutes(router, queries, cfg)
	srv.setupFrontend(router, cfg.FrontendDir)

	return srv
}

func setupSqlc(cfg *config.Config) (*db.Queries, error) {
	dbConfig := database.GetConfig(cfg.DatabaseURL, cfg.DatabaseMaxConns, cfg.DatabaseMaxIdleConns, cfg.DatabaseConnLifetime)

	log.Info("Running database migrations...")
	if err := database.RunMigrations(dbConfig.URL); err != nil {
		log.Error("Failed to run migrations", "error", err)
		return nil, err
	}
	log.Info("Migrations completed successfully")

	log.Info("Setting up database connection pool...")
	conn, err := database.SetupDatabasePool(dbConfig)
	if err != nil {
		log.Error("Failed initializing database pool", "error", err)
		return nil, err
	}
	log.Info("Database connection pool initialized successfully")

	return db.New(conn), nil
}

func (srv *GinServer) setupMiddleware(router *gin.Engine) {
	router.Use(gin.Logger())
	router.Use(middleware.ErrorHandler())
}

func (srv *GinServer) setupRoutes(router *gin.Engine, queries *db.Queries, cfg *config.Config) {
	router.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"status": "healthy"})
	})

	oauthConfig := &oauth2.Config{
		ClientID:     cfg.GoogleClientID,
		ClientSecret: cfg.GoogleClientSecret,
		RedirectURL:  cfg.GoogleRedirectURL,
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
		Endpoint:     google.Endpoint,
	}

	authRepo := auth.NewRepository(queries)
	authService := auth.NewService(authRepo, oauthConfig, cfg.JWTSecret)
	authHandler := auth.NewHandler(authService, oauthConfig)

	shareTokensRepository := sharetokens.NewRepository(queries)
	shareTokensService := sharetokens.NewService(shareTokensRepository)
	shareTokensHandler := sharetokens.NewHandler(shareTokensService)

	commentsRepository := comments.NewRepository(queries)
	commentsService := comments.NewService(commentsRepository)
	commentsHandler := comments.NewHandler(commentsService)

	votesRepository := votes.NewRepository(queries)
	votesService := votes.NewService(votesRepository)
	votesHandler := votes.NewHandler(votesService)

	scrapeAuditRepository := scrapeaudit.NewRepository(queries)
	scrapeAuditService := scrapeaudit.NewService(scrapeAuditRepository)

	scrapeService := scrapeprocess.NewService("TODO: Anthropic Key", scrapeAuditService)

	accomodationsRepository := accomodations.NewRepository(queries)
	accomodationsService := accomodations.NewService(accomodationsRepository, commentsService, votesService, scrapeService)
	accomodationsHandler := accomodations.NewHandler(accomodationsService)

	activitiesRepository := activities.NewRepository(queries)
	activitiesService := activities.NewService(activitiesRepository, commentsService, votesService, scrapeService)
	activitiesHandler := activities.NewHandler(activitiesService)

	transportRepository := transport.NewRepository(queries)
	transportService := transport.NewService(transportRepository, commentsService, votesService, scrapeService)
	transportHandler := transport.NewHandler(transportService)

	boardsRepository := boards.NewRepository(queries)
	boardsService := boards.NewService(boardsRepository, accomodationsService, activitiesService, transportService)
	boardsHandler := boards.NewHandler(boardsService)

	// Public routes (no auth required)
	v1Public := router.Group("/api/v1")
	{
		authHandler.RegisterRoutes(v1Public)
	}

	// Authenticated routes (require valid JWT)
	v1Auth := router.Group("/api/v1")
	v1Auth.Use(middleware.RequireAuth(authService))
	{
		authHandler.RegisterAuthenticatedRoutes(v1Auth)
		boardsHandler.RegisterRoutes(v1Auth)
		shareTokensHandler.RegisterRoutes(v1Auth)
		accomodationsHandler.RegisterRoutes(v1Auth)
		activitiesHandler.RegisterRoutes(v1Auth)
		transportHandler.RegisterRoutes(v1Auth)
		commentsHandler.RegisterRoutes(v1Auth)
		votesHandler.RegisterRoutes(v1Auth)
	}
}

func (srv *GinServer) setupFrontend(router *gin.Engine, frontendDir string) {
	absDir, err := filepath.Abs(frontendDir)
	if err != nil {
		log.Warn("Could not resolve frontend directory path", "dir", frontendDir, "error", err)
		return
	}

	if _, err := os.Stat(absDir); os.IsNotExist(err) {
		log.Warn("Frontend directory not found, skipping static file serving", "dir", absDir)
		return
	}

	log.Info("Serving frontend static files", "dir", absDir)

	router.Static("/assets", filepath.Join(absDir, "assets"))
	router.StaticFile("/vite.svg", filepath.Join(absDir, "vite.svg"))

	indexPath := filepath.Join(absDir, "index.html")
	router.NoRoute(func(c *gin.Context) {
		if strings.HasPrefix(c.Request.URL.Path, "/api/") {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}
		c.File(indexPath)
	})
}
