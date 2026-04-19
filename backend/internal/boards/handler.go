package boards

import "github.com/gin-gonic/gin"

type Handler struct {
	svc Service
}

func NewHandler(svc Service) *Handler {
	return &Handler{
		svc: svc,
	}
}

func (h *Handler) CreateBoard(ctx *gin.Context) {
	var body Board
	if err := ctx.ShouldBindJSON(&body); err != nil {
		log.Error("Failed parsing body", "error", err)
		ctx.AbortWithError(500, err)
		return
	}

	userUuidRaw, exists := ctx.Get("user_uuid")
	if !exists {
		ctx.AbortWithStatusJSON(401, gin.H{"error": "Unauthorized"})
		return
	}
	userUuid := userUuidRaw.(string)
	body.UserUuid = &userUuid

	board, err := h.svc.CreateBoard(ctx, &body)
	if err != nil {
		ctx.AbortWithError(500, err)
		return
	}

	ctx.JSON(200, board)
}

func (h *Handler) GetBoardById(ctx *gin.Context) {
	uuid := ctx.Param("uuid")
	if uuid == "" {
		ctx.JSON(400, gin.H{"error": "uuid parameter is required"})
		return
	}

	board, err := h.svc.GetBoardById(ctx, uuid)
	if err != nil {
		ctx.AbortWithError(500, err)
		return
	}

	ctx.JSON(200, board)
}

func (h *Handler) GetAllBoards(ctx *gin.Context) {
	userUuidRaw, exists := ctx.Get("user_uuid")
	if !exists {
		ctx.AbortWithStatusJSON(401, gin.H{"error": "Unauthorized"})
		return
	}
	userUuid := userUuidRaw.(string)

	boards, err := h.svc.GetAllBoards(ctx, userUuid)
	if err != nil {
		ctx.AbortWithError(500, err)
		return
	}

	ctx.JSON(200, boards)
}

func (h *Handler) UpdateBoardById(ctx *gin.Context) {
	uuid := ctx.Param("uuid")
	if uuid == "" {
		ctx.JSON(400, gin.H{"error": "uuid parameter is required"})
		return
	}

	var body Board
	if err := ctx.ShouldBindJSON(&body); err != nil {
		log.Error("Failed parsing body", "error", err)
		ctx.AbortWithError(500, err)
		return
	}

	err := h.svc.UpdateBoardById(ctx, uuid, &body)
	if err != nil {
		ctx.AbortWithError(500, err)
		return
	}

	ctx.Status(200)
}

func (h *Handler) DeleteBoardById(ctx *gin.Context) {
	uuid := ctx.Param("uuid")
	if uuid == "" {
		ctx.JSON(400, gin.H{"error": "uuid parameter is required"})
		return
	}

	err := h.svc.DeleteBoardById(ctx, uuid)
	if err != nil {
		ctx.AbortWithError(500, err)
		return
	}

	ctx.Status(200)
}

func (h *Handler) RegisterRoutes(router *gin.RouterGroup) {
	group := router.Group("/boards")
	{
		group.POST("/", h.CreateBoard)
		group.GET("/", h.GetAllBoards)
		group.GET("/:uuid", h.GetBoardById)
		group.PATCH("/:uuid", h.UpdateBoardById)
		group.DELETE("/:uuid", h.DeleteBoardById)
	}
}
