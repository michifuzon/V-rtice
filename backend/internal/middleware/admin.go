package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func AutorizarAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.GetString("rol") != "admin" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "acceso denegado",
			})
			return
		}
		c.Next()
	}
}
