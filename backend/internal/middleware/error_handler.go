package middleware

import "github.com/gin-gonic/gin"

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			err := c.Errors.Last()
			c.JSON(c.Writer.Status(), gin.H{
				"error": err.Error(),
			})
		}
	}
}
