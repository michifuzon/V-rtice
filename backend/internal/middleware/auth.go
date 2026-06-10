package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/isaacunaa/ticketek-ds2026/backend/internal/utils"
)

func AutenticacionJWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "falta el header Authorization",
			})
			return
		}

		partes := strings.Split(header, " ")
		if len(partes) != 2 || partes[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "formato de Authorization invalido, se espera 'Bearer <token>'",
			})
			return
		}

		claims, err := utils.ValidarToken(partes[1])
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "token invalido o expirado",
			})
			return
		}

		// jwt.MapClaims decodifica números como float64; convertir a uint explícitamente
		// para que ctx.GetUint("usuario_id") funcione correctamente en los controllers.
		usuarioID := uint(claims["usuario_id"].(float64))
		email := claims["email"].(string)
		rol := claims["rol"].(string)

		c.Set("usuario_id", usuarioID)
		c.Set("email", email)
		c.Set("rol", rol)

		c.Next()
	}
}
