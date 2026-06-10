package controllers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/isaacunaa/ticketek-ds2026/backend/internal/services"
)

type EventoController struct {
	eventoService services.IEventoService
}

func NuevoEventoController(eventoService services.IEventoService) *EventoController {
	return &EventoController{eventoService: eventoService}
}

// Listar maneja GET /api/v1/eventos.
// Query params opcionales: ?categoria=string, ?search=string
func (c *EventoController) Listar(ctx *gin.Context) {
	categoria := ctx.Query("categoria")
	search := ctx.Query("search")

	eventos, err := c.eventoService.ListarEventos(categoria, search)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error al obtener eventos"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"eventos": eventos})
}

// ObtenerPorID maneja GET /api/v1/eventos/:id.
func (c *EventoController) ObtenerPorID(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "id inválido"})
		return
	}

	evento, err := c.eventoService.ObtenerEventoPorID(uint(id))
	if err != nil {
		if errors.Is(err, services.ErrEventoNoEncontrado) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error al obtener el evento"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"evento": evento})
}
