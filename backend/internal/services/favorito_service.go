package services

import (
	"errors"

	"github.com/isaacunaa/ticketek-ds2026/backend/internal/dao"
	"github.com/isaacunaa/ticketek-ds2026/backend/internal/domain"
	"gorm.io/gorm"
)

var (
	ErrEventoNoExiste       = errors.New("el evento no existe")
	ErrYaEsFavorito         = errors.New("el evento ya está en favoritos")
	ErrFavoritoNoEncontrado = errors.New("el favorito no existe")
)

type IFavoritoService interface {
	Agregar(usuarioID, eventoID uint) error
	Quitar(usuarioID, eventoID uint) error
	Listar(usuarioID uint) ([]domain.Favorito, error)
}

type FavoritoService struct {
	favoritoDAO dao.IFavoritoDAO
	eventoDAO   dao.IEventoDAO
}

func NuevoFavoritoService(favoritoDAO dao.IFavoritoDAO, eventoDAO dao.IEventoDAO) *FavoritoService {
	return &FavoritoService{favoritoDAO: favoritoDAO, eventoDAO: eventoDAO}
}

func (s *FavoritoService) Agregar(usuarioID, eventoID uint) error {
	evento, err := s.eventoDAO.BuscarPorID(eventoID)
	if err != nil {
		return err
	}
	if evento == nil {
		return ErrEventoNoExiste
	}

	existe, err := s.favoritoDAO.Existe(usuarioID, eventoID)
	if err != nil {
		return err
	}
	if existe {
		return ErrYaEsFavorito
	}

	return s.favoritoDAO.Agregar(&domain.Favorito{
		UsuarioID: usuarioID,
		EventoID:  eventoID,
	})
}

func (s *FavoritoService) Quitar(usuarioID, eventoID uint) error {
	err := s.favoritoDAO.Quitar(usuarioID, eventoID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrFavoritoNoEncontrado
	}
	return err
}

func (s *FavoritoService) Listar(usuarioID uint) ([]domain.Favorito, error) {
	return s.favoritoDAO.Listar(usuarioID)
}
