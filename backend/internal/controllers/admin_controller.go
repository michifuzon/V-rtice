package controllers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/isaacunaa/ticketek-ds2026/backend/internal/dtos"
	"github.com/isaacunaa/ticketek-ds2026/backend/internal/services"
)

type AdminController struct {
	eventoService services.IEventoService
}

func NuevoAdminController(eventoService services.IEventoService) *AdminController {
	return &AdminController{eventoService: eventoService}
}

// CrearEvento maneja POST /api/v1/admin/eventos
func (c *AdminController) CrearEvento(ctx *gin.Context) {
	var req dtos.EventoRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	evento, err := c.eventoService.Crear(req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error al crear el evento"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"evento": evento})
}

// ActualizarEvento maneja PUT /api/v1/admin/eventos/:id
func (c *AdminController) ActualizarEvento(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "id inválido"})
		return
	}

	var req dtos.EventoRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	evento, err := c.eventoService.Actualizar(uint(id), req)
	if err != nil {
		if errors.Is(err, services.ErrEventoNoEncontrado) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error al actualizar el evento"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"evento": evento})
}

// CancelarEvento maneja DELETE /api/v1/admin/eventos/:id
func (c *AdminController) CancelarEvento(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "id inválido"})
		return
	}

	if err := c.eventoService.Cancelar(uint(id)); err != nil {
		if errors.Is(err, services.ErrEventoNoEncontrado) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error al cancelar el evento"})
		return
	}

	ctx.Status(http.StatusNoContent)
}

// ReporteOcupacion maneja GET /api/v1/admin/eventos/:id/reporte
func (c *AdminController) ReporteOcupacion(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "id inválido"})
		return
	}

	reporte, err := c.eventoService.ObtenerReporte(uint(id))
	if err != nil {
		if errors.Is(err, services.ErrEventoNoEncontrado) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error al obtener el reporte"})
		return
	}

	ctx.JSON(http.StatusOK, reporte)
}
