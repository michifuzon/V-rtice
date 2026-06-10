package domain

import "time"

type Evento struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	Titulo          string    `gorm:"not null;size:200" json:"titulo"`
	Descripcion     string    `gorm:"type:text" json:"descripcion"`
	FechaHora       time.Time `gorm:"not null" json:"fecha_hora"`
	DuracionMinutos int       `json:"duracion_minutos"`
	Ubicacion       string    `gorm:"size:200" json:"ubicacion"`
	Categoria       string    `gorm:"size:50;index" json:"categoria"`
	CupoTotal       int       `gorm:"not null" json:"cupo_total"`
	CupoDisponible  int       `gorm:"not null" json:"cupo_disponible"`
	Precio          float64   `gorm:"not null" json:"precio"`
	ImagenURL       string    `gorm:"size:500" json:"imagen_url"`
	Estado          string    `gorm:"not null;size:20;default:'activo'" json:"estado"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}
