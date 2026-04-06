package main

import (
	"depgraph/internal/config"
	"depgraph/internal/logger"
	"depgraph/internal/server"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	// Load configuration first to get log settings
	cfg, err := config.Load()
	if err != nil {
		// Use default logger if config fails
		defaultLog := logger.New("Main")
		defaultLog.Error("Failed to load configuration", "error", err)
		return
	}

	// Initialize logger with configuration
	if err := logger.Initialize(cfg.LogLevelStdout, cfg.LogLevelFile, cfg.LogFilePath, cfg.LogFilePrefix); err != nil {
		defaultLog := logger.New("Main")
		defaultLog.Error("Failed to initialize logger", "error", err)
		return
	}

	// Now create the main logger after initialization
	log := logger.New("Main")
	log.Info("Initializing DepGraph service...")
	server := server.NewServer()
	if server == nil {
		return
	}

	// Handle signals for graceful shutdown
	shutdownChan := make(chan struct{})
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan

		log.Info("Shutdown signal received")
		server.Stop()
		close(shutdownChan)
	}()

	server.Run()
	<-shutdownChan // Wait for shutdown to complete
}
