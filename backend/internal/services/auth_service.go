package services

import (
	"errors"

	"github.com/isaacunaa/ticketek-ds2026/backend/internal/dao"
	"github.com/isaacunaa/ticketek-ds2026/backend/internal/domain"
	"github.com/isaacunaa/ticketek-ds2026/backend/internal/utils"
)

var (
	ErrEmailYaRegistrado     = errors.New("el email ya está registrado")
	ErrCredencialesInvalidas = errors.New("credenciales inválidas")
)

type IAuthService interface {
	Registrar(email, password, nombre, apellido string) (*domain.Usuario, string, error)
	Login(email, password string) (*domain.Usuario, string, error)
}

type AuthService struct {
	usuarioDAO dao.IUsuarioDAO
}

func NuevoAuthService(usuarioDAO dao.IUsuarioDAO) *AuthService {
	return &AuthService{usuarioDAO: usuarioDAO}
}

func (s *AuthService) Registrar(email, password, nombre, apellido string) (*domain.Usuario, string, error) {
	existente, err := s.usuarioDAO.BuscarPorEmail(email)
	if err != nil {
		return nil, "", err
	}
	if existente != nil {
		return nil, "", ErrEmailYaRegistrado
	}

	hash, err := utils.HashearPassword(password)
	if err != nil {
		return nil, "", err
	}

	usuario := &domain.Usuario{
		Email:        email,
		PasswordHash: hash,
		Nombre:       nombre,
		Apellido:     apellido,
		Rol:          "cliente",
	}

	if err := s.usuarioDAO.Crear(usuario); err != nil {
		return nil, "", err
	}

	token, err := utils.GenerarToken(usuario.ID, usuario.Email, usuario.Rol)
	if err != nil {
		return nil, "", err
	}

	return usuario, token, nil
}

func (s *AuthService) Login(email, password string) (*domain.Usuario, string, error) {
	usuario, err := s.usuarioDAO.BuscarPorEmail(email)
	if err != nil {
		return nil, "", err
	}
	if usuario == nil {
		return nil, "", ErrCredencialesInvalidas
	}

	if err := utils.VerificarPassword(password, usuario.PasswordHash); err != nil {
		return nil, "", ErrCredencialesInvalidas
	}

	token, err := utils.GenerarToken(usuario.ID, usuario.Email, usuario.Rol)
	if err != nil {
		return nil, "", err
	}

	return usuario, token, nil
}
