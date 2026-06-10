package services

import (
	"errors"
	"testing"

	"github.com/isaacunaa/ticketek-ds2026/backend/internal/domain"
	"github.com/isaacunaa/ticketek-ds2026/backend/internal/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// ── Registrar ─────────────────────────────────────────────────────────────────

func TestRegistrar_Exito(t *testing.T) {
	t.Setenv("JWT_SECRET", "secreto-test")

	mockDAO := new(MockUsuarioDAO)
	mockDAO.On("BuscarPorEmail", "nuevo@test.com").Return(nil, nil)
	mockDAO.On("Crear", mock.AnythingOfType("*domain.Usuario")).Return(nil)

	svc := NuevoAuthService(mockDAO)
	usuario, token, err := svc.Registrar("nuevo@test.com", "password123", "Juan", "Pérez")

	assert.NoError(t, err)
	assert.NotNil(t, usuario)
	assert.NotEmpty(t, token)
	assert.Equal(t, "nuevo@test.com", usuario.Email)
	assert.Equal(t, "Juan", usuario.Nombre)
	assert.Equal(t, "cliente", usuario.Rol)
	mockDAO.AssertExpectations(t)
}

func TestRegistrar_EmailDuplicado(t *testing.T) {
	mockDAO := new(MockUsuarioDAO)
	existente := &domain.Usuario{ID: 1, Email: "existente@test.com"}
	mockDAO.On("BuscarPorEmail", "existente@test.com").Return(existente, nil)

	svc := NuevoAuthService(mockDAO)
	usuario, token, err := svc.Registrar("existente@test.com", "pass123", "Ana", "López")

	assert.ErrorIs(t, err, ErrEmailYaRegistrado)
	assert.Nil(t, usuario)
	assert.Empty(t, token)
	mockDAO.AssertExpectations(t)
}

func TestRegistrar_ErrorBuscandoEmail(t *testing.T) {
	mockDAO := new(MockUsuarioDAO)
	errDB := errors.New("conexión perdida")
	mockDAO.On("BuscarPorEmail", "test@test.com").Return(nil, errDB)

	svc := NuevoAuthService(mockDAO)
	usuario, token, err := svc.Registrar("test@test.com", "pass123", "Pedro", "García")

	assert.ErrorIs(t, err, errDB)
	assert.Nil(t, usuario)
	assert.Empty(t, token)
	mockDAO.AssertExpectations(t)
}

func TestRegistrar_ErrorCreandoUsuario(t *testing.T) {
	t.Setenv("JWT_SECRET", "secreto-test")

	mockDAO := new(MockUsuarioDAO)
	errDB := errors.New("fallo al insertar")
	mockDAO.On("BuscarPorEmail", "nuevo@test.com").Return(nil, nil)
	mockDAO.On("Crear", mock.AnythingOfType("*domain.Usuario")).Return(errDB)

	svc := NuevoAuthService(mockDAO)
	usuario, token, err := svc.Registrar("nuevo@test.com", "pass123", "Maria", "Gomez")

	assert.ErrorIs(t, err, errDB)
	assert.Nil(t, usuario)
	assert.Empty(t, token)
	mockDAO.AssertExpectations(t)
}

// ── Login ─────────────────────────────────────────────────────────────────────

func TestLogin_Exito(t *testing.T) {
	t.Setenv("JWT_SECRET", "secreto-test")

	hash, err := utils.HashearPassword("correcta")
	assert.NoError(t, err)

	mockDAO := new(MockUsuarioDAO)
	usuarioMock := &domain.Usuario{
		ID:           1,
		Email:        "login@test.com",
		PasswordHash: hash,
		Rol:          "cliente",
	}
	mockDAO.On("BuscarPorEmail", "login@test.com").Return(usuarioMock, nil)

	svc := NuevoAuthService(mockDAO)
	usuario, token, err := svc.Login("login@test.com", "correcta")

	assert.NoError(t, err)
	assert.NotNil(t, usuario)
	assert.NotEmpty(t, token)
	assert.Equal(t, uint(1), usuario.ID)
	mockDAO.AssertExpectations(t)
}

func TestLogin_UsuarioNoExiste(t *testing.T) {
	mockDAO := new(MockUsuarioDAO)
	mockDAO.On("BuscarPorEmail", "noexiste@test.com").Return(nil, nil)

	svc := NuevoAuthService(mockDAO)
	usuario, token, err := svc.Login("noexiste@test.com", "cualquiera")

	assert.ErrorIs(t, err, ErrCredencialesInvalidas)
	assert.Nil(t, usuario)
	assert.Empty(t, token)
	mockDAO.AssertExpectations(t)
}

func TestLogin_PasswordIncorrecta(t *testing.T) {
	hash, _ := utils.HashearPassword("correcta")

	mockDAO := new(MockUsuarioDAO)
	usuarioMock := &domain.Usuario{
		ID:           2,
		Email:        "pass@test.com",
		PasswordHash: hash,
		Rol:          "cliente",
	}
	mockDAO.On("BuscarPorEmail", "pass@test.com").Return(usuarioMock, nil)

	svc := NuevoAuthService(mockDAO)
	usuario, token, err := svc.Login("pass@test.com", "incorrecta")

	assert.ErrorIs(t, err, ErrCredencialesInvalidas)
	assert.Nil(t, usuario)
	assert.Empty(t, token)
	mockDAO.AssertExpectations(t)
}

func TestLogin_ErrorDAO(t *testing.T) {
	mockDAO := new(MockUsuarioDAO)
	errDB := errors.New("timeout de BD")
	mockDAO.On("BuscarPorEmail", "error@test.com").Return(nil, errDB)

	svc := NuevoAuthService(mockDAO)
	usuario, token, err := svc.Login("error@test.com", "pass")

	assert.ErrorIs(t, err, errDB)
	assert.Nil(t, usuario)
	assert.Empty(t, token)
	mockDAO.AssertExpectations(t)
}
