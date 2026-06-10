package controllers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/isaacunaa/ticketek-ds2026/backend/internal/services"
)

type FavoritoController struct {
	favoritoService services.IFavoritoService
}

func NuevoFavoritoController(favoritoService services.IFavoritoService) *FavoritoController {
	return &FavoritoController{favoritoService: favoritoService}
}

// Agregar maneja POST /api/v1/favoritos/:eventoId
func (c *FavoritoController) Agregar(ctx *gin.Context) {
	eventoID, err := strconv.ParseUint(ctx.Param("eventoId"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "id de evento inválido"})
		return
	}

	usuarioID := ctx.GetUint("usuario_id")

	if err := c.favoritoService.Agregar(usuarioID, uint(eventoID)); err != nil {
		switch {
		case errors.Is(err, services.ErrEventoNoExiste):
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case errors.Is(err, services.ErrYaEsFavorito):
			ctx.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error al agregar favorito"})
		}
		return
	}

	ctx.Status(http.StatusCreated)
}

// Quitar maneja DELETE /api/v1/favoritos/:eventoId
func (c *FavoritoController) Quitar(ctx *gin.Context) {
	eventoID, err := strconv.ParseUint(ctx.Param("eventoId"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "id de evento inválido"})
		return
	}

	usuarioID := ctx.GetUint("usuario_id")

	if err := c.favoritoService.Quitar(usuarioID, uint(eventoID)); err != nil {
		if errors.Is(err, services.ErrFavoritoNoEncontrado) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error al quitar favorito"})
		return
	}

	ctx.Status(http.StatusNoContent)
}

// Listar maneja GET /api/v1/favoritos
func (c *FavoritoController) Listar(ctx *gin.Context) {
	usuarioID := ctx.GetUint("usuario_id")

	favoritos, err := c.favoritoService.Listar(usuarioID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error al obtener favoritos"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"favoritos": favoritos})
}
