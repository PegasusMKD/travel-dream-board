package accomodations

import (
	"path/filepath"

	"github.com/PegasusMKD/travel-dream-board/internal/utility"
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

func (h *Handler) CreateAccomodation(ctx *gin.Context) {
	url := ctx.Query("url")
	file, _ := ctx.FormFile("file")

	if url == "" && file == nil {
		ctx.JSON(400, gin.H{"error": "url or file is required"})
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

	var imageBytes []byte
	var imageExt string
	if file != nil {
		localUrl, bytes, err := utility.SaveUpload(file, "accomodations")
		if err != nil {
			ctx.AbortWithError(500, err)
			return
		}
		url = localUrl
		imageBytes = bytes
		imageExt = filepath.Ext(file.Filename)
	}

	data, err := h.svc.CreateAccomodation(ctx, url, imageBytes, imageExt, boardUuid, userUuid)
	if err != nil {
		ctx.AbortWithError(500, err)
		return
	}

	ctx.JSON(200, data)
}

func (h *Handler) GetAccomodationById(ctx *gin.Context) {
	id := ctx.Param("uuid")
	if id == "" {
		ctx.JSON(400, gin.H{"error": "Uuid parameter is required"})
		return
	}

	data, err := h.svc.GetAccomodationById(ctx, id)
	if err != nil {
		ctx.AbortWithError(500, err)
		return
	}

	ctx.JSON(200, data)
}

func (h *Handler) UpdateAccomodationById(ctx *gin.Context) {
	id := ctx.Param("uuid")
	if id == "" {
		ctx.JSON(400, gin.H{"error": "Uuid parameter is required"})
		return
	}

	var body Accomodation
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.AbortWithError(500, err)
		return
	}

	err := h.svc.UpdateAccomodationById(ctx, id, &body)
	if err != nil {
		ctx.AbortWithError(500, err)
		return
	}

	ctx.Status(200)
}

func (h *Handler) DeleteAccomodationById(ctx *gin.Context) {
	uuid := ctx.Param("uuid")
	if uuid == "" {
		ctx.JSON(400, gin.H{"error": "uuid parameter is required"})
		return
	}

	err := h.svc.DeleteAccomodationById(ctx, uuid)
	if err != nil {
		ctx.AbortWithError(500, err)
		return
	}

	ctx.Status(200)
}

func (h *Handler) RegisterRoutes(router *gin.RouterGroup) {
	group := router.Group("/accomodations")
	{
		group.POST("/", h.CreateAccomodation)
		group.GET("/:uuid", h.GetAccomodationById)
		group.PATCH("/:uuid", h.UpdateAccomodationById)
		group.DELETE("/:uuid", h.DeleteAccomodationById)
	}
}

func (h *Handler) RegisterOwnerRoutes(router *gin.RouterGroup) {
	group := router.Group("/accomodations")
	{
		group.POST("/", h.CreateAccomodation)
		group.PATCH("/:uuid", h.UpdateAccomodationById)
		group.DELETE("/:uuid", h.DeleteAccomodationById)
	}
}

func (h *Handler) RegisterCollabRoutes(router *gin.RouterGroup) {
	router.GET("/accomodations/:uuid", h.GetAccomodationById)
}
