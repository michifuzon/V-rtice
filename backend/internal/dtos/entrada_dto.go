package dtos

type ComprarRequest struct {
	EventoID uint `json:"evento_id" binding:"required"`
	Cantidad int  `json:"cantidad"`
}

type TraspasoRequest struct {
	Email string `json:"email" binding:"required,email"`
}
