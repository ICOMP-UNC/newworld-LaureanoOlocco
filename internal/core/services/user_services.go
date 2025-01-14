package services

import (
	"errors"
	"log"

	"github.com/ICOMP-UNC/newworld-LaureanoOlocco/internal/core/domain"
	"github.com/ICOMP-UNC/newworld-LaureanoOlocco/internal/core/ports"
)

type UserService struct {
	userRepository ports.IUserRepository
}

// Esta línea es para obtener feedback en caso de que no estemos implementando la interfaz correctamente
var _ ports.IUserService = (*UserService)(nil)

func NewUserService(repository ports.IUserRepository) *UserService {
	return &UserService{
		userRepository: repository,
	}
}

func (s *UserService) Login(email string, password string) error {
	err := s.userRepository.Login(email, password)
	if err != nil {
		return err
	}
	return nil
}

func (s *UserService) Register(username, email, password, confirmPass string) error {
	if password != confirmPass {
		return errors.New("the passwords are not equal")
	}
	err := s.userRepository.Register(username, email, password)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (s *UserService) GetUserByEmail(email string) (*domain.User, error) {
	user, err := s.userRepository.GetUserByEmail(email)
	if err != nil {
		return nil, err
	}
	return user, nil
}
