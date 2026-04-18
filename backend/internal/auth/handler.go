package auth

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
)

type Handler struct {
	service Service
	oauth   *oauth2.Config
}

func NewHandler(service Service, oauth *oauth2.Config) *Handler {
	return &Handler{
		service: service,
		oauth:   oauth,
	}
}

func (h *Handler) RegisterRoutes(router *gin.RouterGroup) {
	group := router.Group("/auth")
	{
		group.GET("/google/login", h.GoogleLogin)
		group.GET("/google/callback", h.GoogleCallback)
	}
}

func generateStateOauthCookie(c *gin.Context) string {
	var expiration = 365 * 24 * 60 * 60
	b := make([]byte, 16)
	rand.Read(b)
	state := base64.URLEncoding.EncodeToString(b)
	// Set an HttpOnly cookie
	c.SetCookie("oauthstate", state, expiration, "/", "", false, true)
	return state
}

func (h *Handler) GoogleLogin(c *gin.Context) {
	oauthState := generateStateOauthCookie(c)
	u := h.oauth.AuthCodeURL(oauthState)
	c.Redirect(http.StatusTemporaryRedirect, u)
}

func (h *Handler) GoogleCallback(c *gin.Context) {
	// Verify state
	oauthState, err := c.Cookie("oauthstate")
	if err != nil || c.Query("state") != oauthState {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid oauth state"})
		return
	}

	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Code not found in request"})
		return
	}

	user, token, err := h.service.ProcessGoogleCallback(c.Request.Context(), code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Set JWT as HttpOnly Cookie
	// Set secure to true if running on https (in production)
	// For dev, secure=false might be needed if localhost without https
	c.SetCookie("jwt_token", token, 3600*24*7, "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{
		"message": "Successfully logged in",
		"user": gin.H{
			"uuid":  user.Uuid.String(), // Note: verify user.Uuid string format
			"name":  user.Name,
			"email": user.Email,
		},
	})
}
