package services

import (
	"errors"

	"github.com/isaacunaa/ticketek-ds2026/backend/internal/dao"
	"github.com/isaacunaa/ticketek-ds2026/backend/internal/domain"
	"github.com/isaacunaa/ticketek-ds2026/backend/internal/dtos"
)

var ErrEventoNoEncontrado = errors.New("evento no encontrado")

type IEventoService interface {
	ListarEventos(categoria, search string) ([]domain.Evento, error)
	ObtenerEventoPorID(id uint) (*domain.Evento, error)
	Crear(req dtos.EventoRequest) (*domain.Evento, error)
	Actualizar(id uint, req dtos.EventoRequest) (*domain.Evento, error)
	Cancelar(id uint) error
	ObtenerReporte(eventoID uint) (*dtos.ReporteEvento, error)
}

type EventoService struct {
	eventoDAO  dao.IEventoDAO
	entradaDAO dao.IEntradaDAO
	usuarioDAO dao.IUsuarioDAO
}

func NuevoEventoService(eventoDAO dao.IEventoDAO, entradaDAO dao.IEntradaDAO, usuarioDAO dao.IUsuarioDAO) *EventoService {
	return &EventoService{
		eventoDAO:  eventoDAO,
		entradaDAO: entradaDAO,
		usuarioDAO: usuarioDAO,
	}
}

func (s *EventoService) ListarEventos(categoria, search string) ([]domain.Evento, error) {
	return s.eventoDAO.ListarActivos(categoria, search)
}

func (s *EventoService) ObtenerEventoPorID(id uint) (*domain.Evento, error) {
	evento, err := s.eventoDAO.BuscarPorID(id)
	if err != nil {
		return nil, err
	}
	if evento == nil {
		return nil, ErrEventoNoEncontrado
	}
	return evento, nil
}

func (s *EventoService) Crear(req dtos.EventoRequest) (*domain.Evento, error) {
	evento := &domain.Evento{
		Titulo:          req.Titulo,
		Descripcion:     req.Descripcion,
		FechaHora:       req.FechaHora,
		DuracionMinutos: req.DuracionMinutos,
		Ubicacion:       req.Ubicacion,
		Categoria:       req.Categoria,
		CupoTotal:       req.CupoTotal,
		CupoDisponible:  req.CupoTotal,
		Precio:          req.Precio,
		ImagenURL:       req.ImagenURL,
		Estado:          "activo",
	}
	if err := s.eventoDAO.Crear(evento); err != nil {
		return nil, err
	}
	return evento, nil
}

func (s *EventoService) Actualizar(id uint, req dtos.EventoRequest) (*domain.Evento, error) {
	evento, err := s.eventoDAO.BuscarPorID(id)
	if err != nil {
		return nil, err
	}
	if evento == nil {
		return nil, ErrEventoNoEncontrado
	}
	evento.Titulo          = req.Titulo
	evento.Descripcion     = req.Descripcion
	evento.FechaHora       = req.FechaHora
	evento.DuracionMinutos = req.DuracionMinutos
	evento.Ubicacion       = req.Ubicacion
	evento.Categoria       = req.Categoria
	evento.CupoTotal       = req.CupoTotal
	evento.Precio          = req.Precio
	evento.ImagenURL       = req.ImagenURL
	if err := s.eventoDAO.Actualizar(evento); err != nil {
		return nil, err
	}
	return evento, nil
}

func (s *EventoService) Cancelar(id uint) error {
	evento, err := s.eventoDAO.BuscarPorID(id)
	if err != nil {
		return err
	}
	if evento == nil {
		return ErrEventoNoEncontrado
	}
	if err := s.eventoDAO.CambiarEstado(id, "cancelado"); err != nil {
		return err
	}
	return s.entradaDAO.CancelarPorEvento(id)
}

func (s *EventoService) ObtenerReporte(eventoID uint) (*dtos.ReporteEvento, error) {
	evento, err := s.eventoDAO.BuscarPorID(eventoID)
	if err != nil {
		return nil, err
	}
	if evento == nil {
		return nil, ErrEventoNoEncontrado
	}

	entradas, err := s.entradaDAO.ListarActivasPorEvento(eventoID)
	if err != nil {
		return nil, err
	}

	compradores := make([]dtos.CompradorInfo, 0, len(entradas))

	if len(entradas) > 0 {
		ids := make([]uint, len(entradas))
		for i, e := range entradas {
			ids[i] = e.UsuarioID
		}

		usuarios, err := s.usuarioDAO.BuscarPorIDs(ids)
		if err != nil {
			return nil, err
		}

		userMap := make(map[uint]domain.Usuario, len(usuarios))
		for _, u := range usuarios {
			userMap[u.ID] = u
		}

		for _, e := range entradas {
			if u, ok := userMap[e.UsuarioID]; ok {
				compradores = append(compradores, dtos.CompradorInfo{
					Nombre:      u.Nombre,
					Apellido:    u.Apellido,
					Email:       u.Email,
					FechaCompra: e.FechaCompra,
				})
			}
		}
	}

	vendidas := evento.CupoTotal - evento.CupoDisponible
	var porcentaje float64
	if evento.CupoTotal > 0 {
		porcentaje = float64(vendidas) / float64(evento.CupoTotal) * 100
	}

	return &dtos.ReporteEvento{
		EventoID:            evento.ID,
		Titulo:              evento.Titulo,
		CupoTotal:           evento.CupoTotal,
		CupoDisponible:      evento.CupoDisponible,
		EntradasVendidas:    vendidas,
		PorcentajeOcupacion: porcentaje,
		Compradores:         compradores,
	}, nil
}
