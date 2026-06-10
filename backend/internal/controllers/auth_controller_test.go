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

func init() {
	gin.SetMode(gin.TestMode)
}

func routerAuth(svc services.IAuthService) *gin.Engine {
	r := gin.New()
	ctrl := NuevoAuthController(svc)
	r.POST("/auth/register", ctrl.Registrar)
	r.POST("/auth/login", ctrl.Login)
	return r
}

func hacerRequest(r *gin.Engine, method, path string, body interface{}) *httptest.ResponseRecorder {
	data, _ := json.Marshal(body)
	req := httptest.NewRequest(method, path, bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

// ── Registrar ─────────────────────────────────────────────────────────────────

func TestRegistrar_Exito(t *testing.T) {
	mockSvc := new(MockAuthService)
	usuario := &domain.Usuario{ID: 1, Email: "nuevo@test.com", Nombre: "Juan", Rol: "cliente"}
	mockSvc.On("Registrar", "nuevo@test.com", "pass1234", "Juan", "Pérez").
		Return(usuario, "jwt-token", nil)

	w := hacerRequest(routerAuth(mockSvc), http.MethodPost, "/auth/register", map[string]string{
		"email": "nuevo@test.com", "password": "pass1234",
		"nombre": "Juan", "apellido": "Pérez",
	})

	assert.Equal(t, http.StatusCreated, w.Code)
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, "jwt-token", resp["token"])
	mockSvc.AssertExpectations(t)
}

func TestRegistrar_BodyInvalido(t *testing.T) {
	mockSvc := new(MockAuthService)

	// Falta el campo "nombre" — el binding de gin debe rechazarlo
	w := hacerRequest(routerAuth(mockSvc), http.MethodPost, "/auth/register", map[string]string{
		"email": "test@test.com", "password": "pass1234",
	})

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockSvc.AssertNotCalled(t, "Registrar")
}

func TestRegistrar_EmailDuplicado(t *testing.T) {
	mockSvc := new(MockAuthService)
	mockSvc.On("Registrar", "dup@test.com", "pass1234", "Ana", "López").
		Return(nil, "", services.ErrEmailYaRegistrado)

	w := hacerRequest(routerAuth(mockSvc), http.MethodPost, "/auth/register", map[string]string{
		"email": "dup@test.com", "password": "pass1234",
		"nombre": "Ana", "apellido": "López",
	})

	assert.Equal(t, http.StatusConflict, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestRegistrar_ErrorInterno(t *testing.T) {
	mockSvc := new(MockAuthService)
	mockSvc.On("Registrar", "err@test.com", "pass1234", "Pedro", "García").
		Return(nil, "", errors.New("fallo de BD"))

	w := hacerRequest(routerAuth(mockSvc), http.MethodPost, "/auth/register", map[string]string{
		"email": "err@test.com", "password": "pass1234",
		"nombre": "Pedro", "apellido": "García",
	})

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockSvc.AssertExpectations(t)
}

// ── Login ─────────────────────────────────────────────────────────────────────

func TestLogin_Exito(t *testing.T) {
	mockSvc := new(MockAuthService)
	usuario := &domain.Usuario{ID: 2, Email: "login@test.com", Rol: "cliente"}
	mockSvc.On("Login", "login@test.com", "mipass").
		Return(usuario, "jwt-token-login", nil)

	w := hacerRequest(routerAuth(mockSvc), http.MethodPost, "/auth/login", map[string]string{
		"email": "login@test.com", "password": "mipass",
	})

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, "jwt-token-login", resp["token"])
	mockSvc.AssertExpectations(t)
}

func TestLogin_BodyInvalido(t *testing.T) {
	mockSvc := new(MockAuthService)

	// Falta el campo "password"
	w := hacerRequest(routerAuth(mockSvc), http.MethodPost, "/auth/login", map[string]string{
		"email": "login@test.com",
	})

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockSvc.AssertNotCalled(t, "Login")
}

func TestLogin_CredencialesInvalidas(t *testing.T) {
	mockSvc := new(MockAuthService)
	mockSvc.On("Login", "user@test.com", "mala").
		Return(nil, "", services.ErrCredencialesInvalidas)

	w := hacerRequest(routerAuth(mockSvc), http.MethodPost, "/auth/login", map[string]string{
		"email": "user@test.com", "password": "mala",
	})

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestLogin_ErrorInterno(t *testing.T) {
	mockSvc := new(MockAuthService)
	mockSvc.On("Login", "fallo@test.com", "pass").
		Return(nil, "", errors.New("error de BD"))

	w := hacerRequest(routerAuth(mockSvc), http.MethodPost, "/auth/login", map[string]string{
		"email": "fallo@test.com", "password": "pass",
	})

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockSvc.AssertExpectations(t)
}
