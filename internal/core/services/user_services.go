package services

import (
	"github.com/ICOMP-UNC/newworld-LaureanoOlocco/internal/core/domain"
	"github.com/ICOMP-UNC/newworld-LaureanoOlocco/internal/repositories"
)

// AddUser adds a new user to the database
func CreateUser(register domain.Register) error {
	return repositories.AddUser(register)
}

func LoginUser(login domain.Login) (string, error) {
	return repositories.Login(login)
}
