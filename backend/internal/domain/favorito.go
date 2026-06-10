package domain

import "time"

type Favorito struct {
	UsuarioID uint      `gorm:"primaryKey" json:"usuario_id"`
	EventoID  uint      `gorm:"primaryKey" json:"evento_id"`
	CreatedAt time.Time `json:"created_at"`

	Evento Evento `gorm:"foreignKey:EventoID" json:"evento,omitempty"`
}
