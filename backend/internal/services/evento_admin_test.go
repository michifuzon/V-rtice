package services

import (
	"errors"
	"testing"
	"time"

	"github.com/isaacunaa/ticketek-ds2026/backend/internal/domain"
	"github.com/isaacunaa/ticketek-ds2026/backend/internal/dtos"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// nuevoEventoServiceTest crea un EventoService con entradaDAO y usuarioDAO vacíos.
// Usarlo en tests que solo ejercitan métodos que no los necesitan.
func nuevoEventoServiceTest(eventoDAO *MockEventoDAO) *EventoService {
	return NuevoEventoService(eventoDAO, new(MockEntradaDAO), new(MockUsuarioDAO))
}

// ── Crear ─────────────────────────────────────────────────────────────────────

func TestEventoCrear_Exito(t *testing.T) {
	mockDAO := new(MockEventoDAO)
	mockDAO.On("Crear", mock.AnythingOfType("*domain.Evento")).Return(nil)

	svc := nuevoEventoServiceTest(mockDAO)
	req := dtos.EventoRequest{Titulo: "Festival", CupoTotal: 100, Precio: 5000}
	evento, err := svc.Crear(req)

	assert.NoError(t, err)
	assert.NotNil(t, evento)
	assert.Equal(t, "Festival", evento.Titulo)
	assert.Equal(t, 100, evento.CupoDisponible)
	assert.Equal(t, "activo", evento.Estado)
	mockDAO.AssertExpectations(t)
}

func TestEventoCrear_ErrorDAO(t *testing.T) {
	mockDAO := new(MockEventoDAO)
	errDB := errors.New("fallo de BD")
	mockDAO.On("Crear", mock.AnythingOfType("*domain.Evento")).Return(errDB)

	svc := nuevoEventoServiceTest(mockDAO)
	evento, err := svc.Crear(dtos.EventoRequest{Titulo: "X", CupoTotal: 10, Precio: 100})

	assert.ErrorIs(t, err, errDB)
	assert.Nil(t, evento)
	mockDAO.AssertExpectations(t)
}

// ── Actualizar ────────────────────────────────────────────────────────────────

func TestEventoActualizar_Exito(t *testing.T) {
	mockDAO := new(MockEventoDAO)
	existente := &domain.Evento{ID: 1, Titulo: "Viejo", CupoTotal: 50}
	mockDAO.On("BuscarPorID", uint(1)).Return(existente, nil)
	mockDAO.On("Actualizar", mock.AnythingOfType("*domain.Evento")).Return(nil)

	svc := nuevoEventoServiceTest(mockDAO)
	req := dtos.EventoRequest{Titulo: "Nuevo", CupoTotal: 200, Precio: 9000}
	evento, err := svc.Actualizar(1, req)

	assert.NoError(t, err)
	assert.Equal(t, "Nuevo", evento.Titulo)
	assert.Equal(t, 200, evento.CupoTotal)
	mockDAO.AssertExpectations(t)
}

func TestEventoActualizar_NoEncontrado(t *testing.T) {
	mockDAO := new(MockEventoDAO)
	mockDAO.On("BuscarPorID", uint(99)).Return(nil, nil)

	svc := nuevoEventoServiceTest(mockDAO)
	evento, err := svc.Actualizar(99, dtos.EventoRequest{Titulo: "X", CupoTotal: 10, Precio: 100})

	assert.ErrorIs(t, err, ErrEventoNoEncontrado)
	assert.Nil(t, evento)
	mockDAO.AssertExpectations(t)
}

// ── Cancelar ──────────────────────────────────────────────────────────────────

func TestEventoCancelar_Exito(t *testing.T) {
	mockEventoDAO := new(MockEventoDAO)
	mockEntradaDAO := new(MockEntradaDAO)
	existente := &domain.Evento{ID: 1, Estado: "activo"}
	mockEventoDAO.On("BuscarPorID", uint(1)).Return(existente, nil)
	mockEventoDAO.On("CambiarEstado", uint(1), "cancelado").Return(nil)
	mockEntradaDAO.On("CancelarPorEvento", uint(1)).Return(nil)

	svc := NuevoEventoService(mockEventoDAO, mockEntradaDAO, new(MockUsuarioDAO))
	err := svc.Cancelar(1)

	assert.NoError(t, err)
	mockEventoDAO.AssertExpectations(t)
	mockEntradaDAO.AssertExpectations(t)
}

func TestEventoCancelar_NoEncontrado(t *testing.T) {
	mockDAO := new(MockEventoDAO)
	mockDAO.On("BuscarPorID", uint(99)).Return(nil, nil)

	svc := nuevoEventoServiceTest(mockDAO)
	err := svc.Cancelar(99)

	assert.ErrorIs(t, err, ErrEventoNoEncontrado)
	mockDAO.AssertNotCalled(t, "CambiarEstado")
	mockDAO.AssertExpectations(t)
}

// ── ObtenerReporte ────────────────────────────────────────────────────────────

func TestObtenerReporte_Exito(t *testing.T) {
	mockEventoDAO := new(MockEventoDAO)
	mockEntradaDAO := new(MockEntradaDAO)
	mockUsuarioDAO := new(MockUsuarioDAO)

	evento := &domain.Evento{ID: 1, Titulo: "Festival", CupoTotal: 100, CupoDisponible: 70}
	entradas := []domain.Entrada{
		{ID: 1, UsuarioID: 10, FechaCompra: time.Now()},
		{ID: 2, UsuarioID: 20, FechaCompra: time.Now()},
	}
	usuarios := []domain.Usuario{
		{ID: 10, Nombre: "Ana", Apellido: "García", Email: "ana@test.com"},
		{ID: 20, Nombre: "Luis", Apellido: "López", Email: "luis@test.com"},
	}

	mockEventoDAO.On("BuscarPorID", uint(1)).Return(evento, nil)
	mockEntradaDAO.On("ListarActivasPorEvento", uint(1)).Return(entradas, nil)
	mockUsuarioDAO.On("BuscarPorIDs", []uint{10, 20}).Return(usuarios, nil)

	svc := NuevoEventoService(mockEventoDAO, mockEntradaDAO, mockUsuarioDAO)
	reporte, err := svc.ObtenerReporte(1)

	assert.NoError(t, err)
	assert.NotNil(t, reporte)
	assert.Equal(t, 30, reporte.EntradasVendidas)
	assert.InDelta(t, 30.0, reporte.PorcentajeOcupacion, 0.01)
	assert.Len(t, reporte.Compradores, 2)
	mockEventoDAO.AssertExpectations(t)
	mockEntradaDAO.AssertExpectations(t)
	mockUsuarioDAO.AssertExpectations(t)
}

func TestObtenerReporte_NoEncontrado(t *testing.T) {
	mockDAO := new(MockEventoDAO)
	mockDAO.On("BuscarPorID", uint(99)).Return(nil, nil)

	svc := nuevoEventoServiceTest(mockDAO)
	reporte, err := svc.ObtenerReporte(99)

	assert.ErrorIs(t, err, ErrEventoNoEncontrado)
	assert.Nil(t, reporte)
	mockDAO.AssertExpectations(t)
}

func TestObtenerReporte_SinCompradores(t *testing.T) {
	mockEventoDAO := new(MockEventoDAO)
	mockEntradaDAO := new(MockEntradaDAO)

	evento := &domain.Evento{ID: 2, Titulo: "Vacío", CupoTotal: 50, CupoDisponible: 50}
	mockEventoDAO.On("BuscarPorID", uint(2)).Return(evento, nil)
	mockEntradaDAO.On("ListarActivasPorEvento", uint(2)).Return([]domain.Entrada{}, nil)

	svc := NuevoEventoService(mockEventoDAO, mockEntradaDAO, new(MockUsuarioDAO))
	reporte, err := svc.ObtenerReporte(2)

	assert.NoError(t, err)
	assert.Equal(t, 0, reporte.EntradasVendidas)
	assert.Equal(t, 0.0, reporte.PorcentajeOcupacion)
	assert.Empty(t, reporte.Compradores)
	mockEventoDAO.AssertExpectations(t)
	mockEntradaDAO.AssertExpectations(t)
}
