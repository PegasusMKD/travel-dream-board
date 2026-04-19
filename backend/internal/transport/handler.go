package transport

import "github.com/gin-gonic/gin"

type Handler struct {
	svc Service
}

func NewHandler(svc Service) *Handler {
	return &Handler{
		svc: svc,
	}
}

func (h *Handler) CreateTransport(ctx *gin.Context) {
	url := ctx.Query("url")
	if url == "" {
		ctx.JSON(400, gin.H{"error": "Url parameter is required"})
		return
	}

	boardUuid := ctx.Query("boardUuid")
	if boardUuid == "" {
		ctx.JSON(400, gin.H{"error": "Board UUID parameter is required"})
		return
	}

	userUuidRaw, exists := ctx.Get("user_uuid")
	if !exists {
		ctx.AbortWithStatusJSON(401, gin.H{"error": "Unauthorized"})
		return
	}
	userUuid := userUuidRaw.(string)

	data, err := h.svc.CreateTransport(ctx, url, boardUuid, userUuid)
	if err != nil {
		ctx.AbortWithError(500, err)
		return
	}

	ctx.JSON(200, data)
}

func (h *Handler) GetTransportById(ctx *gin.Context) {
	id := ctx.Param("uuid")
	if id == "" {
		ctx.JSON(400, gin.H{"error": "Uuid parameter is required"})
		return
	}

	data, err := h.svc.GetTransportById(ctx, id)
	if err != nil {
		ctx.AbortWithError(500, err)
		return
	}

	ctx.JSON(200, data)
}

func (h *Handler) UpdateTransportById(ctx *gin.Context) {
	id := ctx.Param("uuid")
	if id == "" {
		ctx.JSON(400, gin.H{"error": "Uuid parameter is required"})
		return
	}

	var body Transport
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.AbortWithError(500, err)
		return
	}

	err := h.svc.UpdateTransportById(ctx, id, &body)
	if err != nil {
		ctx.AbortWithError(500, err)
		return
	}

	ctx.Status(200)
}

func (h *Handler) DeleteTransportById(ctx *gin.Context) {
	uuid := ctx.Param("uuid")
	if uuid == "" {
		ctx.JSON(400, gin.H{"error": "uuid parameter is required"})
		return
	}

	err := h.svc.DeleteTransportById(ctx, uuid)
	if err != nil {
		ctx.AbortWithError(500, err)
		return
	}

	ctx.Status(200)
}

func (h *Handler) RegisterRoutes(router *gin.RouterGroup) {
	group := router.Group("/transport")
	{
		group.POST("/", h.CreateTransport)
		group.GET("/:uuid", h.GetTransportById)
		group.PATCH("/:uuid", h.UpdateTransportById)
		group.DELETE("/:uuid", h.DeleteTransportById)
	}
}
