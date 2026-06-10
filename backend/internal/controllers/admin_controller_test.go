package controllers

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/isaacunaa/ticketek-ds2026/backend/internal/domain"
	"github.com/isaacunaa/ticketek-ds2026/backend/internal/dtos"
	"github.com/isaacunaa/ticketek-ds2026/backend/internal/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func routerAdmin(svc services.IEventoService) *gin.Engine {
	r := gin.New()
	ctrl := NuevoAdminController(svc)
	r.POST("/admin/eventos", ctrl.CrearEvento)
	r.PUT("/admin/eventos/:id", ctrl.ActualizarEvento)
	r.DELETE("/admin/eventos/:id", ctrl.CancelarEvento)
	r.GET("/admin/eventos/:id/reporte", ctrl.ReporteOcupacion)
	return r
}

func putJSON(r *gin.Engine, path string, body interface{}) *httptest.ResponseRecorder {
	return hacerRequest(r, "PUT", path, body)
}

// ── CrearEvento ───────────────────────────────────────────────────────────────

func TestCrearEvento_Exito(t *testing.T) {
	mockSvc := new(MockEventoService)
	evento := &domain.Evento{ID: 1, Titulo: "Festival", Estado: "activo"}
	mockSvc.On("Crear", mock.AnythingOfType("dtos.EventoRequest")).Return(evento, nil)

	w := postJSON(routerAdmin(mockSvc), "/admin/eventos", map[string]interface{}{
		"titulo": "Festival", "cupo_total": 100, "precio": 5000.0,
	})

	assert.Equal(t, http.StatusCreated, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestCrearEvento_BodyInvalido(t *testing.T) {
	mockSvc := new(MockEventoService)

	// falta titulo (required)
	w := postJSON(routerAdmin(mockSvc), "/admin/eventos", map[string]interface{}{
		"cupo_total": 100, "precio": 5000.0,
	})

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockSvc.AssertNotCalled(t, "Crear")
}

func TestCrearEvento_ErrorInterno(t *testing.T) {
	mockSvc := new(MockEventoService)
	mockSvc.On("Crear", mock.AnythingOfType("dtos.EventoRequest")).Return(nil, errors.New("fallo de BD"))

	w := postJSON(routerAdmin(mockSvc), "/admin/eventos", map[string]interface{}{
		"titulo": "Festival", "cupo_total": 100, "precio": 5000.0,
	})

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockSvc.AssertExpectations(t)
}

// ── ActualizarEvento ──────────────────────────────────────────────────────────

func TestActualizarEvento_Exito(t *testing.T) {
	mockSvc := new(MockEventoService)
	evento := &domain.Evento{ID: 1, Titulo: "Actualizado"}
	mockSvc.On("Actualizar", uint(1), mock.AnythingOfType("dtos.EventoRequest")).Return(evento, nil)

	w := putJSON(routerAdmin(mockSvc), "/admin/eventos/1", map[string]interface{}{
		"titulo": "Actualizado", "cupo_total": 200, "precio": 8000.0,
	})

	assert.Equal(t, http.StatusOK, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestActualizarEvento_IDInvalido(t *testing.T) {
	mockSvc := new(MockEventoService)

	w := putJSON(routerAdmin(mockSvc), "/admin/eventos/abc", map[string]interface{}{
		"titulo": "X", "cupo_total": 10, "precio": 100.0,
	})

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockSvc.AssertNotCalled(t, "Actualizar")
}

func TestActualizarEvento_NoEncontrado(t *testing.T) {
	mockSvc := new(MockEventoService)
	mockSvc.On("Actualizar", uint(99), mock.AnythingOfType("dtos.EventoRequest")).
		Return(nil, services.ErrEventoNoEncontrado)

	w := putJSON(routerAdmin(mockSvc), "/admin/eventos/99", map[string]interface{}{
		"titulo": "X", "cupo_total": 10, "precio": 100.0,
	})

	assert.Equal(t, http.StatusNotFound, w.Code)
	mockSvc.AssertExpectations(t)
}

// ── CancelarEvento ────────────────────────────────────────────────────────────

func TestCancelarEvento_Exito(t *testing.T) {
	mockSvc := new(MockEventoService)
	mockSvc.On("Cancelar", uint(1)).Return(nil)

	req := httptest.NewRequest(http.MethodDelete, "/admin/eventos/1", nil)
	w := httptest.NewRecorder()
	routerAdmin(mockSvc).ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestCancelarEvento_IDInvalido(t *testing.T) {
	mockSvc := new(MockEventoService)

	req := httptest.NewRequest(http.MethodDelete, "/admin/eventos/abc", nil)
	w := httptest.NewRecorder()
	routerAdmin(mockSvc).ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockSvc.AssertNotCalled(t, "Cancelar")
}

func TestCancelarEvento_NoEncontrado(t *testing.T) {
	mockSvc := new(MockEventoService)
	mockSvc.On("Cancelar", uint(99)).Return(services.ErrEventoNoEncontrado)

	req := httptest.NewRequest(http.MethodDelete, "/admin/eventos/99", nil)
	w := httptest.NewRecorder()
	routerAdmin(mockSvc).ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	mockSvc.AssertExpectations(t)
}

// ── ReporteOcupacion ──────────────────────────────────────────────────────────

func TestReporteOcupacion_Exito(t *testing.T) {
	mockSvc := new(MockEventoService)
	reporte := &dtos.ReporteEvento{
		EventoID: 1, Titulo: "Festival",
		CupoTotal: 100, CupoDisponible: 70,
		EntradasVendidas: 30, PorcentajeOcupacion: 30.0,
		Compradores: []dtos.CompradorInfo{},
	}
	mockSvc.On("ObtenerReporte", uint(1)).Return(reporte, nil)

	w := getRequest(routerAdmin(mockSvc), "/admin/eventos/1/reporte")

	assert.Equal(t, http.StatusOK, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestReporteOcupacion_IDInvalido(t *testing.T) {
	mockSvc := new(MockEventoService)

	w := getRequest(routerAdmin(mockSvc), "/admin/eventos/abc/reporte")

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockSvc.AssertNotCalled(t, "ObtenerReporte")
}

func TestReporteOcupacion_NoEncontrado(t *testing.T) {
	mockSvc := new(MockEventoService)
	mockSvc.On("ObtenerReporte", uint(99)).Return(nil, services.ErrEventoNoEncontrado)

	w := getRequest(routerAdmin(mockSvc), "/admin/eventos/99/reporte")

	assert.Equal(t, http.StatusNotFound, w.Code)
	mockSvc.AssertExpectations(t)
}
