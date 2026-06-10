package controllers

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/isaacunaa/ticketek-ds2026/backend/internal/domain"
	"github.com/isaacunaa/ticketek-ds2026/backend/internal/services"
	"github.com/stretchr/testify/assert"
)

func routerEvento(svc services.IEventoService) *gin.Engine {
	r := gin.New()
	ctrl := NuevoEventoController(svc)
	r.GET("/eventos", ctrl.Listar)
	r.GET("/eventos/:id", ctrl.ObtenerPorID)
	return r
}

func getRequest(r *gin.Engine, path string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(http.MethodGet, path, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

// ── Listar ────────────────────────────────────────────────────────────────────

func TestListar_Exito(t *testing.T) {
	mockSvc := new(MockEventoService)
	eventos := []domain.Evento{
		{ID: 1, Titulo: "Festival A"},
		{ID: 2, Titulo: "Festival B"},
	}
	mockSvc.On("ListarEventos", "", "").Return(eventos, nil)

	w := getRequest(routerEvento(mockSvc), "/eventos")

	assert.Equal(t, http.StatusOK, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestListar_ConQueryParams(t *testing.T) {
	mockSvc := new(MockEventoService)
	mockSvc.On("ListarEventos", "musica", "rock").Return([]domain.Evento{}, nil)

	w := getRequest(routerEvento(mockSvc), "/eventos?categoria=musica&search=rock")

	assert.Equal(t, http.StatusOK, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestListar_ErrorInterno(t *testing.T) {
	mockSvc := new(MockEventoService)
	mockSvc.On("ListarEventos", "", "").Return(nil, errors.New("fallo de BD"))

	w := getRequest(routerEvento(mockSvc), "/eventos")

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockSvc.AssertExpectations(t)
}

// ── ObtenerPorID ──────────────────────────────────────────────────────────────

func TestObtenerPorID_Exito(t *testing.T) {
	mockSvc := new(MockEventoService)
	evento := &domain.Evento{ID: 5, Titulo: "Gran Concierto", Estado: "activo"}
	mockSvc.On("ObtenerEventoPorID", uint(5)).Return(evento, nil)

	w := getRequest(routerEvento(mockSvc), "/eventos/5")

	assert.Equal(t, http.StatusOK, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestObtenerPorID_IDInvalido(t *testing.T) {
	mockSvc := new(MockEventoService)

	w := getRequest(routerEvento(mockSvc), "/eventos/abc")

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockSvc.AssertNotCalled(t, "ObtenerEventoPorID")
}

func TestObtenerPorID_NoEncontrado(t *testing.T) {
	mockSvc := new(MockEventoService)
	mockSvc.On("ObtenerEventoPorID", uint(99)).Return(nil, services.ErrEventoNoEncontrado)

	w := getRequest(routerEvento(mockSvc), "/eventos/99")

	assert.Equal(t, http.StatusNotFound, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestObtenerPorID_ErrorInterno(t *testing.T) {
	mockSvc := new(MockEventoService)
	mockSvc.On("ObtenerEventoPorID", uint(1)).Return(nil, errors.New("timeout"))

	w := getRequest(routerEvento(mockSvc), "/eventos/1")

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockSvc.AssertExpectations(t)
}
