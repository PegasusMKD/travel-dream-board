package sharetokens

import "github.com/gin-gonic/gin"

type Handler struct {
	svc Service
}

func NewHandler(svc Service) *Handler {
	return &Handler{
		svc: svc,
	}
}

func (h *Handler) CreateShareToken(ctx *gin.Context) {
	boardUuid := ctx.Param("boardUuid")
	if boardUuid == "" {
		ctx.JSON(400, gin.H{"error": "boardUuid parameter is required"})
		return
	}

	token, err := h.svc.CreateShareToken(ctx, boardUuid)
	if err != nil {
		ctx.AbortWithError(500, err)
		return
	}

	ctx.JSON(200, token)
}

func (h *Handler) GetShareTokensForBoard(ctx *gin.Context) {
	boardUuid := ctx.Param("boardUuid")
	if boardUuid == "" {
		ctx.JSON(400, gin.H{"error": "boardUuid parameter is required"})
		return
	}

	tokens, err := h.svc.GetShareTokensForBoard(ctx, boardUuid)
	if err != nil {
		ctx.AbortWithError(500, err)
		return
	}

	ctx.JSON(200, tokens)
}

func (h *Handler) DeleteShareToken(ctx *gin.Context) {
	boardUuid := ctx.Param("boardUuid")
	token := ctx.Param("token")

	if boardUuid == "" || token == "" {
		ctx.JSON(400, gin.H{"error": "boardUuid and token parameters are required"})
		return
	}

	err := h.svc.DeleteShareToken(ctx, token, boardUuid)
	if err != nil {
		ctx.AbortWithError(500, err)
		return
	}

	ctx.Status(200)
}

func (h *Handler) RegisterRoutes(router *gin.RouterGroup) {
	// Let's attach these under /boards/:boardUuid/share-tokens
	group := router.Group("/boards/:boardUuid/share-tokens")
	{
		group.POST("/", h.CreateShareToken)
		group.GET("/", h.GetShareTokensForBoard)
		group.DELETE("/:token", h.DeleteShareToken)
	}
}
