package services

import (
	"errors"
	"testing"

	"github.com/isaacunaa/ticketek-ds2026/backend/internal/domain"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

// ── Agregar ───────────────────────────────────────────────────────────────────

func TestFavoritoAgregar_Exito(t *testing.T) {
	mockFavoritoDAO := new(MockFavoritoDAO)
	mockEventoDAO := new(MockEventoDAO)

	evento := &domain.Evento{ID: 1, Estado: "activo"}
	mockEventoDAO.On("BuscarPorID", uint(1)).Return(evento, nil)
	mockFavoritoDAO.On("Existe", uint(10), uint(1)).Return(false, nil)
	mockFavoritoDAO.On("Agregar", &domain.Favorito{UsuarioID: 10, EventoID: 1}).Return(nil)

	svc := NuevoFavoritoService(mockFavoritoDAO, mockEventoDAO)
	err := svc.Agregar(10, 1)

	assert.NoError(t, err)
	mockEventoDAO.AssertExpectations(t)
	mockFavoritoDAO.AssertExpectations(t)
}

func TestFavoritoAgregar_EventoNoExiste(t *testing.T) {
	mockFavoritoDAO := new(MockFavoritoDAO)
	mockEventoDAO := new(MockEventoDAO)
	mockEventoDAO.On("BuscarPorID", uint(99)).Return(nil, nil)

	svc := NuevoFavoritoService(mockFavoritoDAO, mockEventoDAO)
	err := svc.Agregar(10, 99)

	assert.ErrorIs(t, err, ErrEventoNoExiste)
	mockEventoDAO.AssertExpectations(t)
	mockFavoritoDAO.AssertNotCalled(t, "Existe")
	mockFavoritoDAO.AssertNotCalled(t, "Agregar")
}

func TestFavoritoAgregar_YaEsFavorito(t *testing.T) {
	mockFavoritoDAO := new(MockFavoritoDAO)
	mockEventoDAO := new(MockEventoDAO)

	evento := &domain.Evento{ID: 2, Estado: "activo"}
	mockEventoDAO.On("BuscarPorID", uint(2)).Return(evento, nil)
	mockFavoritoDAO.On("Existe", uint(10), uint(2)).Return(true, nil)

	svc := NuevoFavoritoService(mockFavoritoDAO, mockEventoDAO)
	err := svc.Agregar(10, 2)

	assert.ErrorIs(t, err, ErrYaEsFavorito)
	mockEventoDAO.AssertExpectations(t)
	mockFavoritoDAO.AssertExpectations(t)
	mockFavoritoDAO.AssertNotCalled(t, "Agregar")
}

func TestFavoritoAgregar_ErrorDAO(t *testing.T) {
	mockFavoritoDAO := new(MockFavoritoDAO)
	mockEventoDAO := new(MockEventoDAO)
	errDB := errors.New("fallo de BD")

	evento := &domain.Evento{ID: 3}
	mockEventoDAO.On("BuscarPorID", uint(3)).Return(evento, nil)
	mockFavoritoDAO.On("Existe", uint(10), uint(3)).Return(false, nil)
	mockFavoritoDAO.On("Agregar", &domain.Favorito{UsuarioID: 10, EventoID: 3}).Return(errDB)

	svc := NuevoFavoritoService(mockFavoritoDAO, mockEventoDAO)
	err := svc.Agregar(10, 3)

	assert.ErrorIs(t, err, errDB)
	mockFavoritoDAO.AssertExpectations(t)
}

// ── Quitar ────────────────────────────────────────────────────────────────────

func TestFavoritoQuitar_Exito(t *testing.T) {
	mockFavoritoDAO := new(MockFavoritoDAO)
	mockFavoritoDAO.On("Quitar", uint(10), uint(1)).Return(nil)

	svc := NuevoFavoritoService(mockFavoritoDAO, new(MockEventoDAO))
	err := svc.Quitar(10, 1)

	assert.NoError(t, err)
	mockFavoritoDAO.AssertExpectations(t)
}

func TestFavoritoQuitar_FavoritoNoEncontrado(t *testing.T) {
	mockFavoritoDAO := new(MockFavoritoDAO)
	mockFavoritoDAO.On("Quitar", uint(10), uint(99)).Return(gorm.ErrRecordNotFound)

	svc := NuevoFavoritoService(mockFavoritoDAO, new(MockEventoDAO))
	err := svc.Quitar(10, 99)

	assert.ErrorIs(t, err, ErrFavoritoNoEncontrado)
	mockFavoritoDAO.AssertExpectations(t)
}

// ── Listar ────────────────────────────────────────────────────────────────────

func TestFavoritoListar_Exito(t *testing.T) {
	mockFavoritoDAO := new(MockFavoritoDAO)
	favoritos := []domain.Favorito{
		{UsuarioID: 10, EventoID: 1, Evento: domain.Evento{ID: 1, Titulo: "Festival"}},
		{UsuarioID: 10, EventoID: 2, Evento: domain.Evento{ID: 2, Titulo: "Concierto"}},
	}
	mockFavoritoDAO.On("Listar", uint(10)).Return(favoritos, nil)

	svc := NuevoFavoritoService(mockFavoritoDAO, new(MockEventoDAO))
	resultado, err := svc.Listar(10)

	assert.NoError(t, err)
	assert.Len(t, resultado, 2)
	assert.Equal(t, "Festival", resultado[0].Evento.Titulo)
	mockFavoritoDAO.AssertExpectations(t)
}

func TestFavoritoListar_Vacio(t *testing.T) {
	mockFavoritoDAO := new(MockFavoritoDAO)
	mockFavoritoDAO.On("Listar", uint(10)).Return([]domain.Favorito{}, nil)

	svc := NuevoFavoritoService(mockFavoritoDAO, new(MockEventoDAO))
	resultado, err := svc.Listar(10)

	assert.NoError(t, err)
	assert.Empty(t, resultado)
	mockFavoritoDAO.AssertExpectations(t)
}
