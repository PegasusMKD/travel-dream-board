package middleware

import (
	"net/http"

	"github.com/PegasusMKD/travel-dream-board/internal/auth"
	"github.com/PegasusMKD/travel-dream-board/internal/sharetokens"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
)

// RequireAuth requires a valid JWT in the jwt_token cookie
func RequireAuth(authService auth.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := c.Cookie("jwt_token")
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		userUuid, err := authService.ValidateToken(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		c.Set("user_uuid", userUuid)
		c.Next()
	}
}

// RequireBoardAccess requires either a valid JWT (owner) OR a valid Share Token (guest)
func RequireBoardAccess(authService auth.Service, shareTokenSvc sharetokens.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		boardUuidStr := c.Param("uuid") // assuming route matches /boards/:uuid/...

		// Try JWT first
		tokenString, err := c.Cookie("jwt_token")
		if err == nil {
			userUuid, err := authService.ValidateToken(tokenString)
			if err == nil {
				c.Set("user_uuid", userUuid)
				c.Next()
				return
			}
		}

		// Try Share Token via Header or Query
		shareToken := c.GetHeader("X-Share-Token")
		if shareToken == "" {
			shareToken = c.Query("token")
		}

		if shareToken != "" {
			tokenRecord, err := shareTokenSvc.GetShareToken(c.Request.Context(), shareToken)
			if err == nil {
				// Convert pgtype.UUID to string for comparison or format
				var boardUuid pgtype.UUID
				err = boardUuid.Scan(boardUuidStr) // basic pgtype scanning
				if err == nil && boardUuid.String() == tokenRecord.BoardUuid {
					c.Set("share_token", shareToken)
					c.Next()
					return
				}
			}
		}

		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Forbidden - Board Access Required"})
	}
}
