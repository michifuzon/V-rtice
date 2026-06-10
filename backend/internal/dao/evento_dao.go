package dao

import (
	"errors"

	"github.com/isaacunaa/ticketek-ds2026/backend/internal/domain"
	"gorm.io/gorm"
)

type IEventoDAO interface {
	ListarActivos(categoria, search string) ([]domain.Evento, error)
	BuscarPorID(id uint) (*domain.Evento, error)
	Crear(evento *domain.Evento) error
	Actualizar(evento *domain.Evento) error
	CambiarEstado(id uint, estado string) error
}

type EventoDAO struct {
	db *gorm.DB
}

func NuevoEventoDAO(db *gorm.DB) *EventoDAO {
	return &EventoDAO{db: db}
}

// ListarActivos retorna todos los eventos con estado 'activo'.
// Acepta filtros opcionales: categoria exacta y búsqueda por título (LIKE).
func (d *EventoDAO) ListarActivos(categoria, search string) ([]domain.Evento, error) {
	var eventos []domain.Evento

	query := d.db.Where("estado = ?", "activo")

	if categoria != "" {
		query = query.Where("FIND_IN_SET(?, categoria) > 0", categoria)
	}

	if search != "" {
		query = query.Where("titulo LIKE ?", "%"+search+"%")
	}

	if err := query.Find(&eventos).Error; err != nil {
		return nil, err
	}

	return eventos, nil
}

func (d *EventoDAO) Crear(evento *domain.Evento) error {
	return d.db.Create(evento).Error
}

func (d *EventoDAO) Actualizar(evento *domain.Evento) error {
	return d.db.Save(evento).Error
}

func (d *EventoDAO) CambiarEstado(id uint, estado string) error {
	return d.db.Model(&domain.Evento{}).
		Where("id = ?", id).
		UpdateColumn("estado", estado).Error
}

// BuscarPorID retorna el evento con ese ID, o nil si no existe.
func (d *EventoDAO) BuscarPorID(id uint) (*domain.Evento, error) {
	var evento domain.Evento
	err := d.db.First(&evento, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &evento, nil
}
