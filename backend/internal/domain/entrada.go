package domain

import "time"

type Entrada struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	EventoID     uint      `gorm:"not null;index" json:"evento_id"`
	UsuarioID    uint      `gorm:"not null;index" json:"usuario_id"`
	Codigo       string    `gorm:"uniqueIndex;not null;size:36" json:"codigo"`
	Estado       string    `gorm:"not null;size:20;default:'activa'" json:"estado"`
	PrecioPagado float64   `gorm:"not null" json:"precio_pagado"`
	FechaCompra  time.Time `gorm:"not null" json:"fecha_compra"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`

	Evento  Evento  `gorm:"foreignKey:EventoID" json:"evento,omitempty"`
	Usuario Usuario `gorm:"foreignKey:UsuarioID" json:"usuario,omitempty"`
}
