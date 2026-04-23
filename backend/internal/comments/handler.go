package comments

import (
	"github.com/gin-gonic/gin"
)

type Handler struct {
	svc Service
}

func NewHandler(svc Service) *Handler {
	return &Handler{
		svc: svc,
	}
}

func (h *Handler) CreateComment(ctx *gin.Context) {
	var body Comment
	if err := ctx.ShouldBindJSON(&body); err != nil {
		log.Error("Failed parsing body", "error", err)
		ctx.AbortWithError(500, err)
		return
	}

	data, err := h.svc.CreateComment(ctx, &body)
	if err != nil {
		ctx.AbortWithError(500, err)
		return
	}

	ctx.JSON(200, data)
}

func (h *Handler) UpdateCommentByUuid(ctx *gin.Context) {
	id := ctx.Param("uuid")
	if id == "" {
		ctx.JSON(400, gin.H{"error": "Uuid parameter is required"})
		return
	}

	var body struct {
		Content string `json:"content"`
	}
	if err := ctx.ShouldBindJSON(&body); err != nil {
		log.Error("Failed parsing body", "error", err)
		ctx.AbortWithError(500, err)
		return
	}

	err := h.svc.UpdateCommentByUuid(ctx, id, body.Content)
	if err != nil {
		ctx.AbortWithError(500, err)
		return
	}

	ctx.Status(200)
}

func (h *Handler) DeleteCommentByUuid(ctx *gin.Context) {
	uuid := ctx.Param("uuid")
	if uuid == "" {
		ctx.JSON(400, gin.H{"error": "uuid parameter is required"})
		return
	}

	err := h.svc.DeleteCommentByUuid(ctx, uuid)
	if err != nil {
		ctx.AbortWithError(500, err)
		return
	}

	ctx.Status(200)
}

func (h *Handler) RegisterRoutes(router *gin.RouterGroup) {
	group := router.Group("/comments")
	{
		group.POST("/", h.CreateComment)
		group.PATCH("/:uuid", h.UpdateCommentByUuid)
		group.DELETE("/:uuid", h.DeleteCommentByUuid)
	}
}

func (h *Handler) RegisterCollabRoutes(router *gin.RouterGroup) {
	group := router.Group("/comments")
	{
		group.POST("/", h.CreateComment)
		group.PATCH("/:uuid", h.UpdateCommentByUuid)
		group.DELETE("/:uuid", h.DeleteCommentByUuid)
	}
}
