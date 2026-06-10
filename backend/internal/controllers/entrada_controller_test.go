package controllers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/isaacunaa/ticketek-ds2026/backend/internal/domain"
	"github.com/isaacunaa/ticketek-ds2026/backend/internal/services"
	"github.com/stretchr/testify/assert"
)

// routerEntrada arma el router de test inyectando usuarioID en el contexto
// para simular el middleware JWT, sin necesitar un token real.
func routerEntrada(svc services.IEntradaService, usuarioID uint) *gin.Engine {
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("usuario_id", usuarioID)
		c.Next()
	})
	ctrl := NuevoEntradaController(svc)
	r.POST("/entradas", ctrl.Comprar)
	r.GET("/entradas/me", ctrl.MisEntradas)
	r.DELETE("/entradas/:id", ctrl.Cancelar)
	r.POST("/entradas/:id/transfer", ctrl.Traspasar)
	return r
}

func postJSON(r *gin.Engine, path string, body interface{}) *httptest.ResponseRecorder {
	data, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, path, bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func deleteReq(r *gin.Engine, path string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(http.MethodDelete, path, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

// ── Comprar ───────────────────────────────────────────────────────────────────

func TestComprar_Exito(t *testing.T) {
	mockSvc := new(MockEntradaService)
	entrada := domain.Entrada{ID: 1, EventoID: 5, UsuarioID: 10, Estado: "activa"}
	mockSvc.On("Comprar", uint(10), uint(5), 1).Return([]domain.Entrada{entrada}, nil)

	w := postJSON(routerEntrada(mockSvc, 10), "/entradas", map[string]interface{}{
		"evento_id": 5,
	})

	assert.Equal(t, http.StatusCreated, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestComprar_BodyInvalido(t *testing.T) {
	mockSvc := new(MockEntradaService)

	// Falta evento_id — binding requerido
	w := postJSON(routerEntrada(mockSvc, 10), "/entradas", map[string]interface{}{})

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockSvc.AssertNotCalled(t, "Comprar")
}

func TestComprar_EventoNoDisponible(t *testing.T) {
	mockSvc := new(MockEntradaService)
	mockSvc.On("Comprar", uint(10), uint(99), 1).Return(nil, services.ErrEventoNoDisponible)

	w := postJSON(routerEntrada(mockSvc, 10), "/entradas", map[string]interface{}{
		"evento_id": 99,
	})

	assert.Equal(t, http.StatusNotFound, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestComprar_SinCupo(t *testing.T) {
	mockSvc := new(MockEntradaService)
	mockSvc.On("Comprar", uint(10), uint(3), 1).Return(nil, services.ErrSinCupo)

	w := postJSON(routerEntrada(mockSvc, 10), "/entradas", map[string]interface{}{
		"evento_id": 3,
	})

	assert.Equal(t, http.StatusConflict, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestComprar_ErrorInterno(t *testing.T) {
	mockSvc := new(MockEntradaService)
	mockSvc.On("Comprar", uint(10), uint(1), 1).Return(nil, errors.New("fallo de BD"))

	w := postJSON(routerEntrada(mockSvc, 10), "/entradas", map[string]interface{}{
		"evento_id": 1,
	})

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockSvc.AssertExpectations(t)
}

// ── MisEntradas ───────────────────────────────────────────────────────────────

func TestMisEntradas_Exito(t *testing.T) {
	mockSvc := new(MockEntradaService)
	entradas := []domain.Entrada{
		{ID: 1, UsuarioID: 10, Estado: "activa"},
	}
	mockSvc.On("ListarPorUsuario", uint(10)).Return(entradas, nil)

	req := httptest.NewRequest(http.MethodGet, "/entradas/me", nil)
	w := httptest.NewRecorder()
	routerEntrada(mockSvc, 10).ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestMisEntradas_ErrorInterno(t *testing.T) {
	mockSvc := new(MockEntradaService)
	mockSvc.On("ListarPorUsuario", uint(10)).Return(nil, errors.New("error"))

	req := httptest.NewRequest(http.MethodGet, "/entradas/me", nil)
	w := httptest.NewRecorder()
	routerEntrada(mockSvc, 10).ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockSvc.AssertExpectations(t)
}

// ── Cancelar ──────────────────────────────────────────────────────────────────

func TestCancelar_Exito(t *testing.T) {
	mockSvc := new(MockEntradaService)
	mockSvc.On("Cancelar", uint(10), uint(1)).Return(nil)

	w := deleteReq(routerEntrada(mockSvc, 10), "/entradas/1")

	assert.Equal(t, http.StatusNoContent, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestCancelar_IDInvalido(t *testing.T) {
	mockSvc := new(MockEntradaService)

	w := deleteReq(routerEntrada(mockSvc, 10), "/entradas/abc")

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockSvc.AssertNotCalled(t, "Cancelar")
}

func TestCancelar_NoEncontrada(t *testing.T) {
	mockSvc := new(MockEntradaService)
	mockSvc.On("Cancelar", uint(10), uint(99)).Return(services.ErrEntradaNoEncontrada)

	w := deleteReq(routerEntrada(mockSvc, 10), "/entradas/99")

	assert.Equal(t, http.StatusNotFound, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestCancelar_NoEsDueno(t *testing.T) {
	mockSvc := new(MockEntradaService)
	mockSvc.On("Cancelar", uint(10), uint(2)).Return(services.ErrNoEsDueno)

	w := deleteReq(routerEntrada(mockSvc, 10), "/entradas/2")

	assert.Equal(t, http.StatusForbidden, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestCancelar_EntradaNoActiva(t *testing.T) {
	mockSvc := new(MockEntradaService)
	mockSvc.On("Cancelar", uint(10), uint(3)).Return(services.ErrEntradaNoActiva)

	w := deleteReq(routerEntrada(mockSvc, 10), "/entradas/3")

	assert.Equal(t, http.StatusConflict, w.Code)
	mockSvc.AssertExpectations(t)
}

// ── Traspasar ─────────────────────────────────────────────────────────────────

func TestTraspasar_Exito(t *testing.T) {
	mockSvc := new(MockEntradaService)
	entradaActualizada := &domain.Entrada{ID: 1, UsuarioID: 20, Estado: "activa"}
	mockSvc.On("Traspasar", uint(10), uint(1), "destino@test.com").
		Return(entradaActualizada, nil)

	w := postJSON(routerEntrada(mockSvc, 10), "/entradas/1/transfer", map[string]string{
		"email": "destino@test.com",
	})

	assert.Equal(t, http.StatusOK, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestTraspasar_IDInvalido(t *testing.T) {
	mockSvc := new(MockEntradaService)

	w := postJSON(routerEntrada(mockSvc, 10), "/entradas/abc/transfer", map[string]string{
		"email": "destino@test.com",
	})

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockSvc.AssertNotCalled(t, "Traspasar")
}

func TestTraspasar_BodyInvalido(t *testing.T) {
	mockSvc := new(MockEntradaService)

	// email vacío — el binding lo rechaza (binding:"required,email")
	w := postJSON(routerEntrada(mockSvc, 10), "/entradas/1/transfer", map[string]string{})

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockSvc.AssertNotCalled(t, "Traspasar")
}

func TestTraspasar_EntradaNoEncontrada(t *testing.T) {
	mockSvc := new(MockEntradaService)
	mockSvc.On("Traspasar", uint(10), uint(99), "dest@test.com").
		Return(nil, services.ErrEntradaNoEncontrada)

	w := postJSON(routerEntrada(mockSvc, 10), "/entradas/99/transfer", map[string]string{
		"email": "dest@test.com",
	})

	assert.Equal(t, http.StatusNotFound, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestTraspasar_NoEsDueno(t *testing.T) {
	mockSvc := new(MockEntradaService)
	mockSvc.On("Traspasar", uint(10), uint(1), "dest@test.com").
		Return(nil, services.ErrNoEsDueno)

	w := postJSON(routerEntrada(mockSvc, 10), "/entradas/1/transfer", map[string]string{
		"email": "dest@test.com",
	})

	assert.Equal(t, http.StatusForbidden, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestTraspasar_DestinatarioNoExiste(t *testing.T) {
	mockSvc := new(MockEntradaService)
	mockSvc.On("Traspasar", uint(10), uint(1), "fantasma@test.com").
		Return(nil, services.ErrDestinatarioNoExiste)

	w := postJSON(routerEntrada(mockSvc, 10), "/entradas/1/transfer", map[string]string{
		"email": "fantasma@test.com",
	})

	assert.Equal(t, http.StatusNotFound, w.Code)
	mockSvc.AssertExpectations(t)
}
