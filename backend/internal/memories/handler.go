package memories

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"

	"github.com/PegasusMKD/travel-dream-board/internal/utility"
)

type Handler struct {
	svc Service
}

func NewHandler(svc Service) *Handler {
	return &Handler{svc: svc}
}

// boardAccessAllowed reports whether the current request context is allowed to
// touch resources scoped to boardUuid. Authenticated owners pass through (the
// rest of the API trusts JWT cookies similarly); share-token holders must have
// a token bound to the same board.
func boardAccessAllowed(ctx *gin.Context, boardUuid string) bool {
	if _, ok := ctx.Get("user_uuid"); ok {
		return true
	}
	if v, ok := ctx.Get("share_token_board_uuid"); ok {
		if str, ok := v.(string); ok {
			return str == boardUuid
		}
	}
	return false
}

func (h *Handler) CreateMemory(ctx *gin.Context) {
	boardUuid := ctx.Query("boardUuid")
	if boardUuid == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "boardUuid is required"})
		return
	}
	if !boardAccessAllowed(ctx, boardUuid) {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
		return
	}

	file, err := ctx.FormFile("file")
	if err != nil || file == nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "file is required"})
		return
	}

	// Uploader attribution: prefer JWT user; otherwise fall back to a guest UUID
	// supplied by the share-token client (mirrors comment/vote attribution).
	userUuidRaw, hasUser := ctx.Get("user_uuid")
	var userUuid string
	if hasUser {
		userUuid = userUuidRaw.(string)
	} else {
		userUuid = ctx.Query("uploadedBy")
		if userUuid == "" {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "uploadedBy is required for guest uploads"})
			return
		}
	}

	path, err := utility.SaveMemoryUpload(file)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	mem, err := h.svc.CreateMemory(ctx, boardUuid, userUuid, path)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, mem)
}

func (h *Handler) GetMemoriesByBoard(ctx *gin.Context) {
	boardUuid := ctx.Query("boardUuid")
	if boardUuid == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "boardUuid is required"})
		return
	}
	if !boardAccessAllowed(ctx, boardUuid) {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
		return
	}
	mems, err := h.svc.GetMemoriesByBoardId(ctx, boardUuid)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusOK, mems)
}

func (h *Handler) GetMemoryImage(ctx *gin.Context) {
	uuid := ctx.Param("uuid")
	if uuid == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "uuid is required"})
		return
	}
	mem, err := h.svc.GetMemoryByUuid(ctx, uuid)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if mem == nil {
		ctx.AbortWithStatus(http.StatusNotFound)
		return
	}
	if !boardAccessAllowed(ctx, mem.BoardUuid) {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
		return
	}
	if _, err := os.Stat(mem.ImageUrl); err != nil {
		ctx.AbortWithStatus(http.StatusNotFound)
		return
	}
	ctx.File(mem.ImageUrl)
}

func (h *Handler) DeleteMemory(ctx *gin.Context) {
	uuid := ctx.Param("uuid")
	if uuid == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "uuid is required"})
		return
	}
	mem, err := h.svc.GetMemoryByUuid(ctx, uuid)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if mem == nil {
		ctx.AbortWithStatus(http.StatusNotFound)
		return
	}
	if !boardAccessAllowed(ctx, mem.BoardUuid) {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
		return
	}
	if err := h.svc.DeleteMemoryByUuid(ctx, uuid); err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	ctx.Status(http.StatusOK)
}

func (h *Handler) RegisterCollabRoutes(router *gin.RouterGroup) {
	group := router.Group("/memories")
	{
		group.POST("/", h.CreateMemory)
		group.GET("/", h.GetMemoriesByBoard)
		group.GET("/:uuid/image", h.GetMemoryImage)
	}
}

func (h *Handler) RegisterOwnerRoutes(router *gin.RouterGroup) {
	group := router.Group("/memories")
	{
		group.DELETE("/:uuid", h.DeleteMemory)
	}
}
