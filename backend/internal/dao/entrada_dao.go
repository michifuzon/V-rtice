package dao

import (
	"errors"

	"github.com/isaacunaa/ticketek-ds2026/backend/internal/domain"
	"gorm.io/gorm"
)

type IEntradaDAO interface {
	Crear(tx *gorm.DB, entrada *domain.Entrada) error
	DescontarCupo(tx *gorm.DB, eventoID uint) error
	ListarPorUsuario(usuarioID uint) ([]domain.Entrada, error)
	BuscarPorID(id uint) (*domain.Entrada, error)
	CambiarEstado(tx *gorm.DB, entradaID uint, nuevoEstado string) error
	DevolverCupo(tx *gorm.DB, eventoID uint) error
	CambiarDueno(tx *gorm.DB, entradaID uint, nuevoUsuarioID uint) error
	CancelarPorEvento(eventoID uint) error
	ListarActivasPorEvento(eventoID uint) ([]domain.Entrada, error)
}

type EntradaDAO struct {
	db *gorm.DB
}

func NuevoEntradaDAO(db *gorm.DB) *EntradaDAO {
	return &EntradaDAO{db: db}
}

// Crear inserta una nueva entrada dentro de una transacción.
func (d *EntradaDAO) Crear(tx *gorm.DB, entrada *domain.Entrada) error {
	return tx.Create(entrada).Error
}

// DescontarCupo resta 1 al cupo_disponible del evento dentro de una transacción.
// Solo descuenta si cupo_disponible > 0 (previene negativos).
func (d *EntradaDAO) DescontarCupo(tx *gorm.DB, eventoID uint) error {
	resultado := tx.Model(&domain.Evento{}).
		Where("id = ? AND cupo_disponible > 0", eventoID).
		UpdateColumn("cupo_disponible", gorm.Expr("cupo_disponible - 1"))
	if resultado.Error != nil {
		return resultado.Error
	}
	if resultado.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// ListarPorUsuario retorna todas las entradas de un usuario con el evento incluido.
func (d *EntradaDAO) ListarPorUsuario(usuarioID uint) ([]domain.Entrada, error) {
	var entradas []domain.Entrada
	err := d.db.Preload("Evento").
		Where("usuario_id = ?", usuarioID).
		Find(&entradas).Error
	return entradas, err
}

// BuscarPorID retorna la entrada con ese ID (con Evento precargado), o nil si no existe.
func (d *EntradaDAO) BuscarPorID(id uint) (*domain.Entrada, error) {
	var entrada domain.Entrada
	err := d.db.Preload("Evento").First(&entrada, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &entrada, nil
}

// CambiarEstado actualiza el campo estado de una entrada dentro de una transacción.
func (d *EntradaDAO) CambiarEstado(tx *gorm.DB, entradaID uint, nuevoEstado string) error {
	return tx.Model(&domain.Entrada{}).
		Where("id = ?", entradaID).
		UpdateColumn("estado", nuevoEstado).Error
}

// DevolverCupo suma 1 a cupo_disponible del evento dentro de una transacción.
// Es el inverso de DescontarCupo y se usa al cancelar una entrada.
func (d *EntradaDAO) DevolverCupo(tx *gorm.DB, eventoID uint) error {
	return tx.Model(&domain.Evento{}).
		Where("id = ?", eventoID).
		UpdateColumn("cupo_disponible", gorm.Expr("cupo_disponible + 1")).Error
}

func (d *EntradaDAO) CancelarPorEvento(eventoID uint) error {
	return d.db.Model(&domain.Entrada{}).
		Where("evento_id = ? AND estado = ?", eventoID, "activa").
		UpdateColumn("estado", "cancelada").Error
}

func (d *EntradaDAO) ListarActivasPorEvento(eventoID uint) ([]domain.Entrada, error) {
	var entradas []domain.Entrada
	err := d.db.Where("evento_id = ? AND estado = ?", eventoID, "activa").Find(&entradas).Error
	return entradas, err
}

// CambiarDueno actualiza el usuario_id de la entrada dentro de una transacción.
func (d *EntradaDAO) CambiarDueno(tx *gorm.DB, entradaID uint, nuevoUsuarioID uint) error {
	return tx.Model(&domain.Entrada{}).
		Where("id = ?", entradaID).
		UpdateColumn("usuario_id", nuevoUsuarioID).Error
}

