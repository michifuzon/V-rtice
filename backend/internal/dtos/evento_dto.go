package dtos

import "time"

type EventoRequest struct {
	Titulo          string    `json:"titulo" binding:"required"`
	Descripcion     string    `json:"descripcion"`
	FechaHora       time.Time `json:"fecha_hora"`
	DuracionMinutos int       `json:"duracion_minutos"`
	Ubicacion       string    `json:"ubicacion"`
	Categoria       string    `json:"categoria"`
	CupoTotal       int       `json:"cupo_total" binding:"required,min=1"`
	Precio          float64   `json:"precio" binding:"required"`
	ImagenURL       string    `json:"imagen_url"`
}

type CompradorInfo struct {
	Nombre      string    `json:"nombre"`
	Apellido    string    `json:"apellido"`
	Email       string    `json:"email"`
	FechaCompra time.Time `json:"fecha_compra"`
}

type ReporteEvento struct {
	EventoID            uint          `json:"evento_id"`
	Titulo              string        `json:"titulo"`
	CupoTotal           int           `json:"cupo_total"`
	CupoDisponible      int           `json:"cupo_disponible"`
	EntradasVendidas    int           `json:"entradas_vendidas"`
	PorcentajeOcupacion float64       `json:"porcentaje_ocupacion"`
	Compradores         []CompradorInfo `json:"compradores"`
}
