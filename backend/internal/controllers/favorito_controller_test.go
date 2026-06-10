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

func routerFavorito(svc services.IFavoritoService, usuarioID uint) *gin.Engine {
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("usuario_id", usuarioID)
		c.Next()
	})
	ctrl := NuevoFavoritoController(svc)
	r.POST("/favoritos/:eventoId", ctrl.Agregar)
	r.DELETE("/favoritos/:eventoId", ctrl.Quitar)
	r.GET("/favoritos", ctrl.Listar)
	return r
}

// ── Agregar ───────────────────────────────────────────────────────────────────

func TestFavoritoAgregar_Exito(t *testing.T) {
	mockSvc := new(MockFavoritoService)
	mockSvc.On("Agregar", uint(10), uint(1)).Return(nil)

	req := httptest.NewRequest(http.MethodPost, "/favoritos/1", nil)
	w := httptest.NewRecorder()
	routerFavorito(mockSvc, 10).ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestFavoritoAgregar_IDInvalido(t *testing.T) {
	mockSvc := new(MockFavoritoService)

	req := httptest.NewRequest(http.MethodPost, "/favoritos/abc", nil)
	w := httptest.NewRecorder()
	routerFavorito(mockSvc, 10).ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockSvc.AssertNotCalled(t, "Agregar")
}

func TestFavoritoAgregar_EventoNoExiste(t *testing.T) {
	mockSvc := new(MockFavoritoService)
	mockSvc.On("Agregar", uint(10), uint(99)).Return(services.ErrEventoNoExiste)

	req := httptest.NewRequest(http.MethodPost, "/favoritos/99", nil)
	w := httptest.NewRecorder()
	routerFavorito(mockSvc, 10).ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestFavoritoAgregar_YaEsFavorito(t *testing.T) {
	mockSvc := new(MockFavoritoService)
	mockSvc.On("Agregar", uint(10), uint(1)).Return(services.ErrYaEsFavorito)

	req := httptest.NewRequest(http.MethodPost, "/favoritos/1", nil)
	w := httptest.NewRecorder()
	routerFavorito(mockSvc, 10).ServeHTTP(w, req)

	assert.Equal(t, http.StatusConflict, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestFavoritoAgregar_ErrorInterno(t *testing.T) {
	mockSvc := new(MockFavoritoService)
	mockSvc.On("Agregar", uint(10), uint(1)).Return(errors.New("fallo de BD"))

	req := httptest.NewRequest(http.MethodPost, "/favoritos/1", nil)
	w := httptest.NewRecorder()
	routerFavorito(mockSvc, 10).ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockSvc.AssertExpectations(t)
}

// ── Quitar ────────────────────────────────────────────────────────────────────

func TestFavoritoQuitar_Exito(t *testing.T) {
	mockSvc := new(MockFavoritoService)
	mockSvc.On("Quitar", uint(10), uint(1)).Return(nil)

	req := httptest.NewRequest(http.MethodDelete, "/favoritos/1", nil)
	w := httptest.NewRecorder()
	routerFavorito(mockSvc, 10).ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestFavoritoQuitar_IDInvalido(t *testing.T) {
	mockSvc := new(MockFavoritoService)

	req := httptest.NewRequest(http.MethodDelete, "/favoritos/abc", nil)
	w := httptest.NewRecorder()
	routerFavorito(mockSvc, 10).ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockSvc.AssertNotCalled(t, "Quitar")
}

func TestFavoritoQuitar_NoEncontrado(t *testing.T) {
	mockSvc := new(MockFavoritoService)
	mockSvc.On("Quitar", uint(10), uint(99)).Return(services.ErrFavoritoNoEncontrado)

	req := httptest.NewRequest(http.MethodDelete, "/favoritos/99", nil)
	w := httptest.NewRecorder()
	routerFavorito(mockSvc, 10).ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	mockSvc.AssertExpectations(t)
}

// ── Listar ────────────────────────────────────────────────────────────────────

func TestFavoritoListar_Exito(t *testing.T) {
	mockSvc := new(MockFavoritoService)
	favoritos := []domain.Favorito{
		{UsuarioID: 10, EventoID: 1, Evento: domain.Evento{ID: 1, Titulo: "Festival"}},
	}
	mockSvc.On("Listar", uint(10)).Return(favoritos, nil)

	req := httptest.NewRequest(http.MethodGet, "/favoritos", nil)
	w := httptest.NewRecorder()
	routerFavorito(mockSvc, 10).ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestFavoritoListar_ErrorInterno(t *testing.T) {
	mockSvc := new(MockFavoritoService)
	mockSvc.On("Listar", uint(10)).Return(nil, errors.New("fallo de BD"))

	req := httptest.NewRequest(http.MethodGet, "/favoritos", nil)
	w := httptest.NewRecorder()
	routerFavorito(mockSvc, 10).ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockSvc.AssertExpectations(t)
}
