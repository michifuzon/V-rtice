package services

import (
	"testing"
	"time"

	"github.com/isaacunaa/ticketek-ds2026/backend/internal/dao"
	"github.com/isaacunaa/ticketek-ds2026/backend/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

// helper: crea EntradaService con MockTransactor — evita repetir el wiring en cada test.
func nuevoEntradaServiceTest(entradaDAO dao.IEntradaDAO, eventoDAO dao.IEventoDAO, usuarioDAO dao.IUsuarioDAO) *EntradaService {
	return NuevoEntradaService(entradaDAO, eventoDAO, usuarioDAO, &MockTransactor{})
}

// ── Comprar ───────────────────────────────────────────────────────────────────

func TestComprar_Exito(t *testing.T) {
	mockEntradaDAO := new(MockEntradaDAO)
	mockEventoDAO := new(MockEventoDAO)
	mockUsuarioDAO := new(MockUsuarioDAO)

	evento := &domain.Evento{ID: 1, Estado: "activo", Precio: 1500.0, CupoDisponible: 10}
	mockEventoDAO.On("BuscarPorID", uint(1)).Return(evento, nil)
	mockUsuarioDAO.On("BuscarPorID", uint(10)).Return(&domain.Usuario{SuscripcionActiva: false}, nil)
	mockEntradaDAO.On("DescontarCupo", mock.Anything, uint(1)).Return(nil)
	mockEntradaDAO.On("Crear", mock.Anything, mock.AnythingOfType("*domain.Entrada")).Return(nil)

	svc := nuevoEntradaServiceTest(mockEntradaDAO, mockEventoDAO, mockUsuarioDAO)
	entradas, err := svc.Comprar(10, 1, 1)

	assert.NoError(t, err)
	assert.Len(t, entradas, 1)
	assert.Equal(t, uint(1), entradas[0].EventoID)
	assert.Equal(t, uint(10), entradas[0].UsuarioID)
	assert.Equal(t, "activa", entradas[0].Estado)
	assert.Equal(t, 1500.0, entradas[0].PrecioPagado)
	assert.NotEmpty(t, entradas[0].Codigo)
	mockEntradaDAO.AssertExpectations(t)
	mockEventoDAO.AssertExpectations(t)
}

func TestComprar_EventoNoExiste(t *testing.T) {
	mockEventoDAO := new(MockEventoDAO)
	mockEventoDAO.On("BuscarPorID", uint(99)).Return(nil, nil)

	svc := nuevoEntradaServiceTest(new(MockEntradaDAO), mockEventoDAO, new(MockUsuarioDAO))
	entradas, err := svc.Comprar(1, 99, 1)

	assert.ErrorIs(t, err, ErrEventoNoDisponible)
	assert.Nil(t, entradas)
	mockEventoDAO.AssertExpectations(t)
}

func TestComprar_EventoNoActivo(t *testing.T) {
	mockEventoDAO := new(MockEventoDAO)
	evento := &domain.Evento{ID: 2, Estado: "inactivo"}
	mockEventoDAO.On("BuscarPorID", uint(2)).Return(evento, nil)

	svc := nuevoEntradaServiceTest(new(MockEntradaDAO), mockEventoDAO, new(MockUsuarioDAO))
	entradas, err := svc.Comprar(1, 2, 1)

	assert.ErrorIs(t, err, ErrEventoNoDisponible)
	assert.Nil(t, entradas)
	mockEventoDAO.AssertExpectations(t)
}

func TestComprar_ConDescuentoClub(t *testing.T) {
	mockEntradaDAO := new(MockEntradaDAO)
	mockEventoDAO := new(MockEventoDAO)
	mockUsuarioDAO := new(MockUsuarioDAO)

	evento := &domain.Evento{ID: 1, Estado: "activo", Precio: 1000.0, CupoDisponible: 10}
	vence := time.Now().Add(24 * time.Hour)
	usuario := &domain.Usuario{SuscripcionActiva: true, SuscripcionVence: &vence}
	mockEventoDAO.On("BuscarPorID", uint(1)).Return(evento, nil)
	mockUsuarioDAO.On("BuscarPorID", uint(5)).Return(usuario, nil)
	mockEntradaDAO.On("DescontarCupo", mock.Anything, uint(1)).Return(nil)
	mockEntradaDAO.On("Crear", mock.Anything, mock.AnythingOfType("*domain.Entrada")).Return(nil)

	svc := nuevoEntradaServiceTest(mockEntradaDAO, mockEventoDAO, mockUsuarioDAO)
	entradas, err := svc.Comprar(5, 1, 1)

	assert.NoError(t, err)
	assert.Len(t, entradas, 1)
	assert.Equal(t, 900.0, entradas[0].PrecioPagado)
	mockEntradaDAO.AssertExpectations(t)
	mockEventoDAO.AssertExpectations(t)
}

func TestComprar_SuscripcionVencida(t *testing.T) {
	mockEntradaDAO := new(MockEntradaDAO)
	mockEventoDAO := new(MockEventoDAO)
	mockUsuarioDAO := new(MockUsuarioDAO)

	evento := &domain.Evento{ID: 2, Estado: "activo", Precio: 1000.0, CupoDisponible: 5}
	vencida := time.Now().Add(-24 * time.Hour)
	usuario := &domain.Usuario{SuscripcionActiva: true, SuscripcionVence: &vencida}
	mockEventoDAO.On("BuscarPorID", uint(2)).Return(evento, nil)
	mockUsuarioDAO.On("BuscarPorID", uint(5)).Return(usuario, nil)
	mockEntradaDAO.On("DescontarCupo", mock.Anything, uint(2)).Return(nil)
	mockEntradaDAO.On("Crear", mock.Anything, mock.AnythingOfType("*domain.Entrada")).Return(nil)

	svc := nuevoEntradaServiceTest(mockEntradaDAO, mockEventoDAO, mockUsuarioDAO)
	entradas, err := svc.Comprar(5, 2, 1)

	assert.NoError(t, err)
	assert.Len(t, entradas, 1)
	assert.Equal(t, 1000.0, entradas[0].PrecioPagado)
	mockEntradaDAO.AssertExpectations(t)
	mockEventoDAO.AssertExpectations(t)
}

func TestComprar_MultipleEntradas(t *testing.T) {
	mockEntradaDAO := new(MockEntradaDAO)
	mockEventoDAO := new(MockEventoDAO)
	mockUsuarioDAO := new(MockUsuarioDAO)

	evento := &domain.Evento{ID: 3, Estado: "activo", Precio: 500.0, CupoDisponible: 10}
	mockEventoDAO.On("BuscarPorID", uint(3)).Return(evento, nil)
	mockUsuarioDAO.On("BuscarPorID", uint(7)).Return(&domain.Usuario{SuscripcionActiva: false}, nil)
	mockEntradaDAO.On("DescontarCupo", mock.Anything, uint(3)).Return(nil).Times(3)
	mockEntradaDAO.On("Crear", mock.Anything, mock.AnythingOfType("*domain.Entrada")).Return(nil).Times(3)

	svc := nuevoEntradaServiceTest(mockEntradaDAO, mockEventoDAO, mockUsuarioDAO)
	entradas, err := svc.Comprar(7, 3, 3)

	assert.NoError(t, err)
	assert.Len(t, entradas, 3)
	for _, e := range entradas {
		assert.Equal(t, uint(3), e.EventoID)
		assert.Equal(t, 500.0, e.PrecioPagado)
		assert.NotEmpty(t, e.Codigo)
	}
	mockEntradaDAO.AssertExpectations(t)
	mockEventoDAO.AssertExpectations(t)
}

func TestComprar_SinCupo(t *testing.T) {
	mockEntradaDAO := new(MockEntradaDAO)
	mockEventoDAO := new(MockEventoDAO)

	mockUsuarioDAO := new(MockUsuarioDAO)
	evento := &domain.Evento{ID: 3, Estado: "activo", CupoDisponible: 0}
	mockEventoDAO.On("BuscarPorID", uint(3)).Return(evento, nil)
	mockUsuarioDAO.On("BuscarPorID", uint(1)).Return(&domain.Usuario{SuscripcionActiva: false}, nil)
	// DescontarCupo devuelve ErrRecordNotFound cuando cupo_disponible = 0
	mockEntradaDAO.On("DescontarCupo", mock.Anything, uint(3)).Return(gorm.ErrRecordNotFound)

	svc := nuevoEntradaServiceTest(mockEntradaDAO, mockEventoDAO, mockUsuarioDAO)
	entradas, err := svc.Comprar(1, 3, 1)

	assert.ErrorIs(t, err, ErrSinCupo)
	assert.Nil(t, entradas)
	mockEntradaDAO.AssertExpectations(t)
	mockEventoDAO.AssertExpectations(t)
}

// ── ListarPorUsuario ──────────────────────────────────────────────────────────

func TestListarPorUsuario_Exito(t *testing.T) {
	mockEntradaDAO := new(MockEntradaDAO)
	entradas := []domain.Entrada{
		{ID: 1, UsuarioID: 5, Estado: "activa"},
		{ID: 2, UsuarioID: 5, Estado: "cancelada"},
	}
	mockEntradaDAO.On("ListarPorUsuario", uint(5)).Return(entradas, nil)

	svc := nuevoEntradaServiceTest(mockEntradaDAO, new(MockEventoDAO), new(MockUsuarioDAO))
	resultado, err := svc.ListarPorUsuario(5)

	assert.NoError(t, err)
	assert.Len(t, resultado, 2)
	mockEntradaDAO.AssertExpectations(t)
}

// ── Cancelar ──────────────────────────────────────────────────────────────────

func TestCancelar_Exito(t *testing.T) {
	mockEntradaDAO := new(MockEntradaDAO)
	entrada := &domain.Entrada{ID: 1, UsuarioID: 10, Estado: "activa", EventoID: 3}
	mockEntradaDAO.On("BuscarPorID", uint(1)).Return(entrada, nil)
	mockEntradaDAO.On("CambiarEstado", mock.Anything, uint(1), "cancelada").Return(nil)
	mockEntradaDAO.On("DevolverCupo", mock.Anything, uint(3)).Return(nil)

	svc := nuevoEntradaServiceTest(mockEntradaDAO, new(MockEventoDAO), new(MockUsuarioDAO))
	err := svc.Cancelar(10, 1)

	assert.NoError(t, err)
	mockEntradaDAO.AssertExpectations(t)
}

func TestCancelar_EntradaNoEncontrada(t *testing.T) {
	mockEntradaDAO := new(MockEntradaDAO)
	mockEntradaDAO.On("BuscarPorID", uint(99)).Return(nil, nil)

	svc := nuevoEntradaServiceTest(mockEntradaDAO, new(MockEventoDAO), new(MockUsuarioDAO))
	err := svc.Cancelar(1, 99)

	assert.ErrorIs(t, err, ErrEntradaNoEncontrada)
	mockEntradaDAO.AssertExpectations(t)
}

func TestCancelar_NoEsDueno(t *testing.T) {
	mockEntradaDAO := new(MockEntradaDAO)
	// La entrada pertenece al usuario 99; el usuario 10 intenta cancelarla
	entrada := &domain.Entrada{ID: 1, UsuarioID: 99, Estado: "activa"}
	mockEntradaDAO.On("BuscarPorID", uint(1)).Return(entrada, nil)

	svc := nuevoEntradaServiceTest(mockEntradaDAO, new(MockEventoDAO), new(MockUsuarioDAO))
	err := svc.Cancelar(10, 1)

	assert.ErrorIs(t, err, ErrNoEsDueno)
	mockEntradaDAO.AssertExpectations(t)
}

func TestCancelar_EntradaNoActiva(t *testing.T) {
	mockEntradaDAO := new(MockEntradaDAO)
	entrada := &domain.Entrada{ID: 1, UsuarioID: 10, Estado: "cancelada"}
	mockEntradaDAO.On("BuscarPorID", uint(1)).Return(entrada, nil)

	svc := nuevoEntradaServiceTest(mockEntradaDAO, new(MockEventoDAO), new(MockUsuarioDAO))
	err := svc.Cancelar(10, 1)

	assert.ErrorIs(t, err, ErrEntradaNoActiva)
	mockEntradaDAO.AssertExpectations(t)
}

// ── Traspasar ─────────────────────────────────────────────────────────────────

func TestTraspasar_Exito(t *testing.T) {
	mockEntradaDAO := new(MockEntradaDAO)
	mockUsuarioDAO := new(MockUsuarioDAO)

	entradaOriginal := &domain.Entrada{ID: 1, UsuarioID: 10, Estado: "activa", EventoID: 5}
	entradaActualizada := &domain.Entrada{ID: 1, UsuarioID: 20, Estado: "activa", EventoID: 5}
	destinatario := &domain.Usuario{ID: 20, Email: "destino@test.com"}

	// Primera llamada: cargar la entrada para validar dueño/estado
	mockEntradaDAO.On("BuscarPorID", uint(1)).Return(entradaOriginal, nil).Once()
	mockUsuarioDAO.On("BuscarPorEmail", "destino@test.com").Return(destinatario, nil)
	mockEntradaDAO.On("CambiarDueno", mock.Anything, uint(1), uint(20)).Return(nil)
	// Segunda llamada: recargar la entrada actualizada para la respuesta
	mockEntradaDAO.On("BuscarPorID", uint(1)).Return(entradaActualizada, nil).Once()

	svc := nuevoEntradaServiceTest(mockEntradaDAO, new(MockEventoDAO), mockUsuarioDAO)
	resultado, err := svc.Traspasar(10, 1, "destino@test.com")

	assert.NoError(t, err)
	assert.NotNil(t, resultado)
	assert.Equal(t, uint(20), resultado.UsuarioID)
	mockEntradaDAO.AssertExpectations(t)
	mockUsuarioDAO.AssertExpectations(t)
}

func TestTraspasar_EntradaNoEncontrada(t *testing.T) {
	mockEntradaDAO := new(MockEntradaDAO)
	mockEntradaDAO.On("BuscarPorID", uint(99)).Return(nil, nil)

	svc := nuevoEntradaServiceTest(mockEntradaDAO, new(MockEventoDAO), new(MockUsuarioDAO))
	resultado, err := svc.Traspasar(1, 99, "otro@test.com")

	assert.ErrorIs(t, err, ErrEntradaNoEncontrada)
	assert.Nil(t, resultado)
	mockEntradaDAO.AssertExpectations(t)
}

func TestTraspasar_NoEsDueno(t *testing.T) {
	mockEntradaDAO := new(MockEntradaDAO)
	entrada := &domain.Entrada{ID: 1, UsuarioID: 99, Estado: "activa"}
	mockEntradaDAO.On("BuscarPorID", uint(1)).Return(entrada, nil)

	svc := nuevoEntradaServiceTest(mockEntradaDAO, new(MockEventoDAO), new(MockUsuarioDAO))
	resultado, err := svc.Traspasar(10, 1, "otro@test.com")

	assert.ErrorIs(t, err, ErrNoEsDueno)
	assert.Nil(t, resultado)
	mockEntradaDAO.AssertExpectations(t)
}

func TestTraspasar_EntradaNoActiva(t *testing.T) {
	mockEntradaDAO := new(MockEntradaDAO)
	entrada := &domain.Entrada{ID: 1, UsuarioID: 10, Estado: "cancelada"}
	mockEntradaDAO.On("BuscarPorID", uint(1)).Return(entrada, nil)

	svc := nuevoEntradaServiceTest(mockEntradaDAO, new(MockEventoDAO), new(MockUsuarioDAO))
	resultado, err := svc.Traspasar(10, 1, "otro@test.com")

	assert.ErrorIs(t, err, ErrEntradaNoActiva)
	assert.Nil(t, resultado)
	mockEntradaDAO.AssertExpectations(t)
}

func TestTraspasar_DestinatarioNoExiste(t *testing.T) {
	mockEntradaDAO := new(MockEntradaDAO)
	mockUsuarioDAO := new(MockUsuarioDAO)
	entrada := &domain.Entrada{ID: 1, UsuarioID: 10, Estado: "activa"}
	mockEntradaDAO.On("BuscarPorID", uint(1)).Return(entrada, nil)
	mockUsuarioDAO.On("BuscarPorEmail", "fantasma@test.com").Return(nil, nil)

	svc := nuevoEntradaServiceTest(mockEntradaDAO, new(MockEventoDAO), mockUsuarioDAO)
	resultado, err := svc.Traspasar(10, 1, "fantasma@test.com")

	assert.ErrorIs(t, err, ErrDestinatarioNoExiste)
	assert.Nil(t, resultado)
	mockEntradaDAO.AssertExpectations(t)
	mockUsuarioDAO.AssertExpectations(t)
}
