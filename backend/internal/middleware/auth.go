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
// scoped to the board referenced by the :uuid path param.
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

// RequireAuthOrShareToken requires either a valid JWT or any valid (non-revoked)
// share token. Unlike RequireBoardAccess, it does not bind the token to the
// request's board UUID — useful for routes where the board UUID is not in the
// path (e.g. /votes/:uuid, /comments/:uuid, /accomodations/:uuid where :uuid is
// an item, not a board). Ownership enforcement for mutating endpoints lives in
// the handler.
func RequireAuthOrShareToken(authService auth.Service, shareTokenSvc sharetokens.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := c.Cookie("jwt_token")
		if err == nil {
			userUuid, err := authService.ValidateToken(tokenString)
			if err == nil {
				c.Set("user_uuid", userUuid)
				c.Next()
				return
			}
		}

		shareToken := c.GetHeader("X-Share-Token")
		if shareToken == "" {
			shareToken = c.Query("token")
		}

		if shareToken != "" {
			tokenRecord, err := shareTokenSvc.GetShareToken(c.Request.Context(), shareToken)
			if err == nil {
				c.Set("share_token", shareToken)
				c.Set("share_token_board_uuid", tokenRecord.BoardUuid)
				c.Next()
				return
			}
		}

		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
	}
}
