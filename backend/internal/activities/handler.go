package activities

import "github.com/gin-gonic/gin"

type Handler struct {
	svc Service
}

func NewHandler(svc Service) *Handler {
	return &Handler{
		svc: svc,
	}
}

func (h *Handler) CreateActivity(ctx *gin.Context) {
	url := ctx.Query("url")
	if url == "" {
		ctx.JSON(400, gin.H{"error": "Url parameter is required"})
		return
	}

	data, err := h.svc.CreateActivity(ctx, url)
	if err != nil {
		ctx.AbortWithError(500, err)
		return
	}

	ctx.JSON(200, data)
}

func (h *Handler) GetActivityById(ctx *gin.Context) {
	id := ctx.Param("uuid")
	if id == "" {
		ctx.JSON(400, gin.H{"error": "Uuid parameter is required"})
		return
	}

	data, err := h.svc.GetActivityById(ctx, id)
	if err != nil {
		ctx.AbortWithError(500, err)
		return
	}

	ctx.JSON(200, data)
}

func (h *Handler) UpdateActivityById(ctx *gin.Context) {
	id := ctx.Param("uuid")
	if id == "" {
		ctx.JSON(400, gin.H{"error": "Uuid parameter is required"})
		return
	}

	var body Activity
	if err := ctx.ShouldBindJSON(body); err != nil {
		log.Error("Failed parsing body", "error", err)
		ctx.AbortWithError(500, err)
		return
	}

	err := h.svc.UpdateActivityById(ctx, id, &body)
	if err != nil {
		ctx.AbortWithError(500, err)
		return
	}

	ctx.Status(200)
}

func (h *Handler) DeleteActivityById(ctx *gin.Context) {
	uuid := ctx.Param("uuid")
	if uuid == "" {
		ctx.JSON(400, gin.H{"error": "uuid parameter is required"})
		return
	}

	err := h.svc.DeleteActivityById(ctx, uuid)
	if err != nil {
		ctx.AbortWithError(500, err)
		return
	}

	ctx.Status(200)
}

func (h *Handler) RegisterRoutes(router *gin.RouterGroup) {
	group := router.Group("/accomodations")
	{
		group.POST("/", h.CreateActivity)
		group.GET("/:uuid", h.GetActivityById)
		group.PATCH("/:uuid", h.UpdateActivityById)
		group.DELETE("/:uuid", h.DeleteActivityById)
	}
}
