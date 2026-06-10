package dao

import (
	"errors"
	"time"

	"github.com/isaacunaa/ticketek-ds2026/backend/internal/domain"
	"gorm.io/gorm"
)

type IUsuarioDAO interface {
	BuscarPorEmail(email string) (*domain.Usuario, error)
	Crear(usuario *domain.Usuario) error
	BuscarPorID(id uint) (*domain.Usuario, error)
	BuscarPorIDs(ids []uint) ([]domain.Usuario, error)
	ActualizarSuscripcion(usuarioID uint, activa bool, vence *time.Time) error
}

type UsuarioDAO struct {
	db *gorm.DB
}

func NuevoUsuarioDAO(db *gorm.DB) *UsuarioDAO {
	return &UsuarioDAO{db: db}
}

// Crear inserta un nuevo usuario en la base.
func (d *UsuarioDAO) Crear(usuario *domain.Usuario) error {
	return d.db.Create(usuario).Error
}

// BuscarPorEmail retorna el usuario con ese email, o nil si no existe.
func (d *UsuarioDAO) BuscarPorEmail(email string) (*domain.Usuario, error) {
	var usuario domain.Usuario
	err := d.db.Where("email = ?", email).First(&usuario).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &usuario, nil
}

func (d *UsuarioDAO) BuscarPorIDs(ids []uint) ([]domain.Usuario, error) {
	var usuarios []domain.Usuario
	err := d.db.Where("id IN ?", ids).Find(&usuarios).Error
	return usuarios, err
}

// BuscarPorID retorna el usuario con ese ID, o nil si no existe.
func (d *UsuarioDAO) BuscarPorID(id uint) (*domain.Usuario, error) {
	var usuario domain.Usuario
	err := d.db.First(&usuario, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &usuario, nil
}

// ActualizarSuscripcion actualiza el estado de suscripción del usuario.
func (d *UsuarioDAO) ActualizarSuscripcion(usuarioID uint, activa bool, vence *time.Time) error {
	return d.db.Model(&domain.Usuario{}).
		Where("id = ?", usuarioID).
		Updates(map[string]interface{}{
			"suscripcion_activa": activa,
			"suscripcion_vence":  vence,
		}).Error
}
