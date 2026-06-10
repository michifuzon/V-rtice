package services

import (
	"errors"
	"testing"

	"github.com/isaacunaa/ticketek-ds2026/backend/internal/domain"
	"github.com/stretchr/testify/assert"
)

// ── ListarEventos ─────────────────────────────────────────────────────────────

func TestListarEventos_Exito(t *testing.T) {
	mockDAO := new(MockEventoDAO)
	eventos := []domain.Evento{
		{ID: 1, Titulo: "Evento A", Estado: "activo"},
		{ID: 2, Titulo: "Evento B", Estado: "activo"},
	}
	mockDAO.On("ListarActivos", "", "").Return(eventos, nil)

	svc := NuevoEventoService(mockDAO, new(MockEntradaDAO), new(MockUsuarioDAO))
	resultado, err := svc.ListarEventos("", "")

	assert.NoError(t, err)
	assert.Len(t, resultado, 2)
	assert.Equal(t, "Evento A", resultado[0].Titulo)
	mockDAO.AssertExpectations(t)
}

func TestListarEventos_ConFiltros(t *testing.T) {
	mockDAO := new(MockEventoDAO)
	eventos := []domain.Evento{{ID: 1, Titulo: "Rock Fest", Categoria: "musica"}}
	mockDAO.On("ListarActivos", "musica", "rock").Return(eventos, nil)

	svc := NuevoEventoService(mockDAO, new(MockEntradaDAO), new(MockUsuarioDAO))
	resultado, err := svc.ListarEventos("musica", "rock")

	assert.NoError(t, err)
	assert.Len(t, resultado, 1)
	assert.Equal(t, "musica", resultado[0].Categoria)
	mockDAO.AssertExpectations(t)
}

func TestListarEventos_Error(t *testing.T) {
	mockDAO := new(MockEventoDAO)
	mockDAO.On("ListarActivos", "", "").Return(nil, errors.New("error de BD"))

	svc := NuevoEventoService(mockDAO, new(MockEntradaDAO), new(MockUsuarioDAO))
	resultado, err := svc.ListarEventos("", "")

	assert.Error(t, err)
	assert.Nil(t, resultado)
	mockDAO.AssertExpectations(t)
}

// ── ObtenerEventoPorID ────────────────────────────────────────────────────────

func TestObtenerEventoPorID_Exito(t *testing.T) {
	mockDAO := new(MockEventoDAO)
	evento := &domain.Evento{ID: 5, Titulo: "Festival", Estado: "activo"}
	mockDAO.On("BuscarPorID", uint(5)).Return(evento, nil)

	svc := NuevoEventoService(mockDAO, new(MockEntradaDAO), new(MockUsuarioDAO))
	resultado, err := svc.ObtenerEventoPorID(5)

	assert.NoError(t, err)
	assert.NotNil(t, resultado)
	assert.Equal(t, uint(5), resultado.ID)
	assert.Equal(t, "Festival", resultado.Titulo)
	mockDAO.AssertExpectations(t)
}

func TestObtenerEventoPorID_NoEncontrado(t *testing.T) {
	mockDAO := new(MockEventoDAO)
	mockDAO.On("BuscarPorID", uint(99)).Return(nil, nil)

	svc := NuevoEventoService(mockDAO, new(MockEntradaDAO), new(MockUsuarioDAO))
	resultado, err := svc.ObtenerEventoPorID(99)

	assert.ErrorIs(t, err, ErrEventoNoEncontrado)
	assert.Nil(t, resultado)
	mockDAO.AssertExpectations(t)
}

func TestObtenerEventoPorID_ErrorDAO(t *testing.T) {
	mockDAO := new(MockEventoDAO)
	mockDAO.On("BuscarPorID", uint(1)).Return(nil, errors.New("timeout"))

	svc := NuevoEventoService(mockDAO, new(MockEntradaDAO), new(MockUsuarioDAO))
	resultado, err := svc.ObtenerEventoPorID(1)

	assert.Error(t, err)
	assert.Nil(t, resultado)
	mockDAO.AssertExpectations(t)
}
