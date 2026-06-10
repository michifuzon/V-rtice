package controllers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/isaacunaa/ticketek-ds2026/backend/internal/dtos"
	"github.com/isaacunaa/ticketek-ds2026/backend/internal/services"
)

type EntradaController struct {
	entradaService services.IEntradaService
}

func NuevoEntradaController(entradaService services.IEntradaService) *EntradaController {
	return &EntradaController{entradaService: entradaService}
}

// Comprar maneja POST /api/v1/entradas
func (c *EntradaController) Comprar(ctx *gin.Context) {
	var req dtos.ComprarRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Cantidad < 1 {
		req.Cantidad = 1
	}
	if req.Cantidad > 10 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "el máximo de entradas por compra es 10"})
		return
	}

	usuarioID := ctx.GetUint("usuario_id")

	entradas, err := c.entradaService.Comprar(usuarioID, req.EventoID, req.Cantidad)
	if err != nil {
		if errors.Is(err, services.ErrEventoNoDisponible) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, services.ErrSinCupo) {
			ctx.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error al procesar la compra"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"entradas": entradas, "cantidad": len(entradas)})
}

// MisEntradas maneja GET /api/v1/entradas/me
func (c *EntradaController) MisEntradas(ctx *gin.Context) {
	usuarioID := ctx.GetUint("usuario_id")

	entradas, err := c.entradaService.ListarPorUsuario(usuarioID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error al obtener entradas"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"entradas": entradas})
}

// Cancelar maneja DELETE /api/v1/entradas/:id
func (c *EntradaController) Cancelar(ctx *gin.Context) {
	entradaID, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "id de entrada inválido"})
		return
	}

	usuarioID := ctx.GetUint("usuario_id")

	if err := c.entradaService.Cancelar(usuarioID, uint(entradaID)); err != nil {
		switch {
		case errors.Is(err, services.ErrEntradaNoEncontrada):
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case errors.Is(err, services.ErrNoEsDueno):
			ctx.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		case errors.Is(err, services.ErrEntradaNoActiva):
			ctx.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error al cancelar la entrada"})
		}
		return
	}

	ctx.Status(http.StatusNoContent)
}

// Traspasar maneja POST /api/v1/entradas/:id/transfer
func (c *EntradaController) Traspasar(ctx *gin.Context) {
	entradaID, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "id de entrada inválido"})
		return
	}

	var req dtos.TraspasoRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	usuarioID := ctx.GetUint("usuario_id")

	entrada, err := c.entradaService.Traspasar(usuarioID, uint(entradaID), req.Email)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrEntradaNoEncontrada):
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case errors.Is(err, services.ErrNoEsDueno):
			ctx.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		case errors.Is(err, services.ErrEntradaNoActiva):
			ctx.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		case errors.Is(err, services.ErrDestinatarioNoExiste):
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error al traspasar la entrada"})
		}
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"entrada": entrada})
}
