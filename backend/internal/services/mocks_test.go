package services

import (
	"database/sql"
	"time"

	"github.com/isaacunaa/ticketek-ds2026/backend/internal/domain"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

// ── MockUsuarioDAO ────────────────────────────────────────────────────────────

type MockUsuarioDAO struct {
	mock.Mock
}

func (m *MockUsuarioDAO) BuscarPorEmail(email string) (*domain.Usuario, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Usuario), args.Error(1)
}

func (m *MockUsuarioDAO) Crear(usuario *domain.Usuario) error {
	args := m.Called(usuario)
	return args.Error(0)
}

func (m *MockUsuarioDAO) BuscarPorID(id uint) (*domain.Usuario, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Usuario), args.Error(1)
}

func (m *MockUsuarioDAO) BuscarPorIDs(ids []uint) ([]domain.Usuario, error) {
	args := m.Called(ids)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Usuario), args.Error(1)
}

func (m *MockUsuarioDAO) ActualizarSuscripcion(usuarioID uint, activa bool, vence *time.Time) error {
	args := m.Called(usuarioID, activa, vence)
	return args.Error(0)
}

// ── MockEventoDAO ─────────────────────────────────────────────────────────────

type MockEventoDAO struct {
	mock.Mock
}

func (m *MockEventoDAO) ListarActivos(categoria, search string) ([]domain.Evento, error) {
	args := m.Called(categoria, search)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Evento), args.Error(1)
}

func (m *MockEventoDAO) BuscarPorID(id uint) (*domain.Evento, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Evento), args.Error(1)
}

func (m *MockEventoDAO) Crear(evento *domain.Evento) error {
	args := m.Called(evento)
	return args.Error(0)
}

func (m *MockEventoDAO) Actualizar(evento *domain.Evento) error {
	args := m.Called(evento)
	return args.Error(0)
}

func (m *MockEventoDAO) CambiarEstado(id uint, estado string) error {
	args := m.Called(id, estado)
	return args.Error(0)
}

// ── MockEntradaDAO ────────────────────────────────────────────────────────────

type MockEntradaDAO struct {
	mock.Mock
}

func (m *MockEntradaDAO) Crear(tx *gorm.DB, entrada *domain.Entrada) error {
	args := m.Called(tx, entrada)
	return args.Error(0)
}

func (m *MockEntradaDAO) DescontarCupo(tx *gorm.DB, eventoID uint) error {
	args := m.Called(tx, eventoID)
	return args.Error(0)
}

func (m *MockEntradaDAO) ListarPorUsuario(usuarioID uint) ([]domain.Entrada, error) {
	args := m.Called(usuarioID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Entrada), args.Error(1)
}

func (m *MockEntradaDAO) BuscarPorID(id uint) (*domain.Entrada, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Entrada), args.Error(1)
}

func (m *MockEntradaDAO) CambiarEstado(tx *gorm.DB, entradaID uint, nuevoEstado string) error {
	args := m.Called(tx, entradaID, nuevoEstado)
	return args.Error(0)
}

func (m *MockEntradaDAO) DevolverCupo(tx *gorm.DB, eventoID uint) error {
	args := m.Called(tx, eventoID)
	return args.Error(0)
}

func (m *MockEntradaDAO) CambiarDueno(tx *gorm.DB, entradaID uint, nuevoUsuarioID uint) error {
	args := m.Called(tx, entradaID, nuevoUsuarioID)
	return args.Error(0)
}

func (m *MockEntradaDAO) CancelarPorEvento(eventoID uint) error {
	args := m.Called(eventoID)
	return args.Error(0)
}

func (m *MockEntradaDAO) ListarActivasPorEvento(eventoID uint) ([]domain.Entrada, error) {
	args := m.Called(eventoID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Entrada), args.Error(1)
}


// ── MockFavoritoDAO ───────────────────────────────────────────────────────────

type MockFavoritoDAO struct {
	mock.Mock
}

func (m *MockFavoritoDAO) Agregar(favorito *domain.Favorito) error {
	args := m.Called(favorito)
	return args.Error(0)
}

func (m *MockFavoritoDAO) Quitar(usuarioID, eventoID uint) error {
	args := m.Called(usuarioID, eventoID)
	return args.Error(0)
}

func (m *MockFavoritoDAO) Listar(usuarioID uint) ([]domain.Favorito, error) {
	args := m.Called(usuarioID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Favorito), args.Error(1)
}

func (m *MockFavoritoDAO) Existe(usuarioID, eventoID uint) (bool, error) {
	args := m.Called(usuarioID, eventoID)
	return args.Bool(0), args.Error(1)
}

// ── MockTransactor ────────────────────────────────────────────────────────────

// MockTransactor implementa ITransactor ejecutando el callback con tx=nil.
// Los MockDAO ignoran el parámetro tx, por lo que esto es suficiente para tests unitarios.
type MockTransactor struct{}

func (m *MockTransactor) Transaction(fc func(tx *gorm.DB) error, _ ...*sql.TxOptions) error {
	return fc(nil)
}
