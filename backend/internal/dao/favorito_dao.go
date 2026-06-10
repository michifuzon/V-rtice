package dao

import (
	"errors"

	"github.com/isaacunaa/ticketek-ds2026/backend/internal/domain"
	"gorm.io/gorm"
)

type IFavoritoDAO interface {
	Agregar(favorito *domain.Favorito) error
	Quitar(usuarioID, eventoID uint) error
	Listar(usuarioID uint) ([]domain.Favorito, error)
	Existe(usuarioID, eventoID uint) (bool, error)
}

type FavoritoDAO struct {
	db *gorm.DB
}

func NuevoFavoritoDAO(db *gorm.DB) *FavoritoDAO {
	return &FavoritoDAO{db: db}
}

func (d *FavoritoDAO) Agregar(favorito *domain.Favorito) error {
	return d.db.Create(favorito).Error
}

func (d *FavoritoDAO) Quitar(usuarioID, eventoID uint) error {
	resultado := d.db.Where("usuario_id = ? AND evento_id = ?", usuarioID, eventoID).
		Delete(&domain.Favorito{})
	if resultado.Error != nil {
		return resultado.Error
	}
	if resultado.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (d *FavoritoDAO) Listar(usuarioID uint) ([]domain.Favorito, error) {
	var favoritos []domain.Favorito
	err := d.db.Preload("Evento").Where("usuario_id = ?", usuarioID).Find(&favoritos).Error
	return favoritos, err
}

func (d *FavoritoDAO) Existe(usuarioID, eventoID uint) (bool, error) {
	var f domain.Favorito
	err := d.db.Where("usuario_id = ? AND evento_id = ?", usuarioID, eventoID).First(&f).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
