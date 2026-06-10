package controllers

import (
	"github.com/isaacunaa/ticketek-ds2026/backend/internal/domain"
	"github.com/isaacunaa/ticketek-ds2026/backend/internal/dtos"
	"github.com/stretchr/testify/mock"
)

// ── MockAuthService ───────────────────────────────────────────────────────────

type MockAuthService struct {
	mock.Mock
}

func (m *MockAuthService) Registrar(email, password, nombre, apellido string) (*domain.Usuario, string, error) {
	args := m.Called(email, password, nombre, apellido)
	if args.Get(0) == nil {
		return nil, args.String(1), args.Error(2)
	}
	return args.Get(0).(*domain.Usuario), args.String(1), args.Error(2)
}

func (m *MockAuthService) Login(email, password string) (*domain.Usuario, string, error) {
	args := m.Called(email, password)
	if args.Get(0) == nil {
		return nil, args.String(1), args.Error(2)
	}
	return args.Get(0).(*domain.Usuario), args.String(1), args.Error(2)
}

// ── MockEventoService ─────────────────────────────────────────────────────────

type MockEventoService struct {
	mock.Mock
}

func (m *MockEventoService) ListarEventos(categoria, search string) ([]domain.Evento, error) {
	args := m.Called(categoria, search)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Evento), args.Error(1)
}

func (m *MockEventoService) ObtenerEventoPorID(id uint) (*domain.Evento, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Evento), args.Error(1)
}

func (m *MockEventoService) Crear(req dtos.EventoRequest) (*domain.Evento, error) {
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Evento), args.Error(1)
}

func (m *MockEventoService) Actualizar(id uint, req dtos.EventoRequest) (*domain.Evento, error) {
	args := m.Called(id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Evento), args.Error(1)
}

func (m *MockEventoService) Cancelar(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockEventoService) ObtenerReporte(eventoID uint) (*dtos.ReporteEvento, error) {
	args := m.Called(eventoID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dtos.ReporteEvento), args.Error(1)
}

// ── MockFavoritoService ───────────────────────────────────────────────────────

type MockFavoritoService struct {
	mock.Mock
}

func (m *MockFavoritoService) Agregar(usuarioID, eventoID uint) error {
	args := m.Called(usuarioID, eventoID)
	return args.Error(0)
}

func (m *MockFavoritoService) Quitar(usuarioID, eventoID uint) error {
	args := m.Called(usuarioID, eventoID)
	return args.Error(0)
}

func (m *MockFavoritoService) Listar(usuarioID uint) ([]domain.Favorito, error) {
	args := m.Called(usuarioID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Favorito), args.Error(1)
}

// ── MockEntradaService ────────────────────────────────────────────────────────

type MockEntradaService struct {
	mock.Mock
}

func (m *MockEntradaService) Comprar(usuarioID, eventoID uint, cantidad int) ([]domain.Entrada, error) {
	args := m.Called(usuarioID, eventoID, cantidad)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Entrada), args.Error(1)
}

func (m *MockEntradaService) ListarPorUsuario(usuarioID uint) ([]domain.Entrada, error) {
	args := m.Called(usuarioID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Entrada), args.Error(1)
}

func (m *MockEntradaService) Cancelar(usuarioID, entradaID uint) error {
	args := m.Called(usuarioID, entradaID)
	return args.Error(0)
}

func (m *MockEntradaService) Traspasar(usuarioID, entradaID uint, emailDestinatario string) (*domain.Entrada, error) {
	args := m.Called(usuarioID, entradaID, emailDestinatario)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Entrada), args.Error(1)
}

