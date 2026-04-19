package votes

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

func (h *Handler) CreateVote(ctx *gin.Context) {
	var body Vote
	if err := ctx.ShouldBindJSON(&body); err != nil {
		log.Error("Failed parsing body", "error", err)
		ctx.AbortWithError(500, err)
		return
	}

	data, err := h.svc.CreateVote(ctx, &body)
	if err != nil {
		ctx.AbortWithError(500, err)
		return
	}

	ctx.JSON(200, data)
}

func (h *Handler) UpdateVoteByUuid(ctx *gin.Context) {
	id := ctx.Param("uuid")
	if id == "" {
		ctx.JSON(400, gin.H{"error": "Uuid parameter is required"})
		return
	}

	var body struct {
		Rank int32 `json:"rank"`
	}
	if err := ctx.ShouldBindJSON(&body); err != nil {
		log.Error("Failed parsing body", "error", err)
		ctx.AbortWithError(500, err)
		return
	}

	err := h.svc.UpdateVoteByUuid(ctx, id, body.Rank)
	if err != nil {
		ctx.AbortWithError(500, err)
		return
	}

	ctx.Status(200)
}

func (h *Handler) DeleteVoteByUuid(ctx *gin.Context) {
	uuid := ctx.Param("uuid")
	if uuid == "" {
		ctx.JSON(400, gin.H{"error": "uuid parameter is required"})
		return
	}

	err := h.svc.DeleteVoteByUuid(ctx, uuid)
	if err != nil {
		ctx.AbortWithError(500, err)
		return
	}

	ctx.Status(200)
}

func (h *Handler) RegisterRoutes(router *gin.RouterGroup) {
	group := router.Group("/votes")
	{
		group.POST("/", h.CreateVote)
		group.PATCH("/:uuid", h.UpdateVoteByUuid)
		group.DELETE("/:uuid", h.DeleteVoteByUuid)
	}
}
