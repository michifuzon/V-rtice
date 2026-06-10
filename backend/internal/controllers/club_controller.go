package controllers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/isaacunaa/ticketek-ds2026/backend/internal/dao"
)

const PrecioSuscripcion = 5000.0

type ClubController struct {
	usuarioDAO dao.IUsuarioDAO
}

func NuevoClubController(usuarioDAO dao.IUsuarioDAO) *ClubController {
	return &ClubController{usuarioDAO: usuarioDAO}
}

// EstadoClub maneja GET /api/v1/club/estado
func (c *ClubController) EstadoClub(ctx *gin.Context) {
	usuarioID := ctx.GetUint("usuario_id")
	usuario, err := c.usuarioDAO.BuscarPorID(usuarioID)
	if err != nil || usuario == nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error al obtener usuario"})
		return
	}

	activa := usuario.SuscripcionActiva &&
		usuario.SuscripcionVence != nil &&
		time.Now().Before(*usuario.SuscripcionVence)

	ctx.JSON(http.StatusOK, gin.H{
		"activa": activa,
		"vence":  usuario.SuscripcionVence,
		"precio": PrecioSuscripcion,
	})
}

// CancelarSuscripcion maneja DELETE /api/v1/club/suscripcion
func (c *ClubController) CancelarSuscripcion(ctx *gin.Context) {
	usuarioID := ctx.GetUint("usuario_id")
	if err := c.usuarioDAO.ActualizarSuscripcion(usuarioID, false, nil); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error al cancelar la suscripción"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"activa":  false,
		"mensaje": "Membresía cancelada. Podés volver a suscribirte cuando quieras.",
	})
}

// Suscribirse maneja POST /api/v1/club/suscribir
func (c *ClubController) Suscribirse(ctx *gin.Context) {
	usuarioID := ctx.GetUint("usuario_id")
	usuario, err := c.usuarioDAO.BuscarPorID(usuarioID)
	if err != nil || usuario == nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error al obtener usuario"})
		return
	}

	// Si tiene suscripción activa, extender desde la fecha de vencimiento
	var base time.Time
	if usuario.SuscripcionActiva && usuario.SuscripcionVence != nil && time.Now().Before(*usuario.SuscripcionVence) {
		base = *usuario.SuscripcionVence
	} else {
		base = time.Now()
	}

	vence := base.AddDate(0, 1, 0) // +1 mes
	if err := c.usuarioDAO.ActualizarSuscripcion(usuarioID, true, &vence); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error al activar suscripción"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"activa":  true,
		"vence":   vence,
		"mensaje": "¡Bienvenido al Club Vórtice! Tu suscripción está activa.",
	})
}
