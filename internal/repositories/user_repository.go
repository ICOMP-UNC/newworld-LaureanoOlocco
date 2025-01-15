package repositories

import (
	"database/sql"
	"errors"

	"github.com/ICOMP-UNC/newworld-LaureanoOlocco/internal/core/domain"
	"github.com/ICOMP-UNC/newworld-LaureanoOlocco/internal/core/ports"
	"github.com/lib/pq"
)

type UserRepository struct {
	db *sql.DB
}

// Asure that UserRepository implements IUserRepository
var _ ports.IUserRepository = (*UserRepository)(nil)

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) Login(email string, password string) error {
	var storedPassword string
	query := `SELECT password FROM users WHERE email = $1`
	err := r.db.QueryRow(query, email).Scan(&storedPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("user not found")
		}
		return err
	}

	if storedPassword != password {
		return errors.New("invalid password")
	}

	return nil
}

func (r *UserRepository) Register(username, email, password string) error {
	// Insert the user into the database
	query := `INSERT INTO users (username, email, password) VALUES ($1, $2, $3)`
	_, err := r.db.Exec(query, username, email, password)
	if err != nil {

		// Check if the error is a unique violation
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code.Name() == "unique_violation" {
				return errors.New("email already registered")
			}
		}
		return err
	}

	return nil
}

// function to get user by email
func (r *UserRepository) GetUserByEmail(email string) (*domain.User, error) {
	var user domain.User
	query := `SELECT id, email, password FROM users WHERE email = $1`
	err := r.db.QueryRow(query, email).Scan(&user.ID, &user.Email, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return &user, nil
}
