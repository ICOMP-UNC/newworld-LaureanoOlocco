package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	//"github.com/ICOMP-UNC/newworld-LaureanoOlocco/internal/core/domain"
	"github.com/ICOMP-UNC/newworld-LaureanoOlocco/internal/core/domain"
	"github.com/golang-jwt/jwt/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

var dbPool *pgxpool.Pool

// SetDB sets the database pool for analytics package
func SetDB(pool *pgxpool.Pool) {
	dbPool = pool
}

/**
*---------------------------------------------------------------------------------|
*		USER REPOSITORY
*---------------------------------------------------------------------------------
 */

// AddUser adds a new user to the database
func AddUser(register domain.Register) error {

	fmt.Println("¡Hola, mundo!")

	if dbPool == nil {
		return errors.New("database pool is not initialized")
	}
	var role string
	if register.Username == "ubuntu" || register.Email == "ubuntu@ubuntu.com" || register.Password == "ubuntu" {
		role = "admin"
	} else {
		role = "user"
	}
	// GenerateJWT(register.Username)
	jwt, err0 := GenerateJWT(register.Username, role)
	if err0 != nil {
		return err0
	}

	_, err := dbPool.Exec(context.Background(), "INSERT INTO users (username, email, password, role, jwt) VALUES ($1, $2, $3, $4, $5)", register.Username, register.Email, register.Password, role, jwt)
	if err != nil {
		return err
	}

	return nil
}

func Login(login domain.Login) (string, error) {
	if dbPool == nil {
		return "", errors.New("database pool is not initialized")
	}

	var jwtKey string
	err := dbPool.QueryRow(context.Background(), "SELECT jwt FROM users WHERE email = $1 AND password = $2", login.Email, login.Password).Scan(&jwtKey)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", errors.New("user does not exist or invalid credentials")
		}
		return "", err
	}

	return jwtKey, nil
}

// ---------- J W T 	 M A N A G G I N G -----------

func GenerateJWT(username string, role string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["username"] = username
	claims["role"] = role

	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	tokenString, err := token.SignedString([]byte("your_secret_key"))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func CheckJWT(jwtString string) (bool, error) {
	// Parse the JWT
	token, err := jwt.Parse(jwtString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return false, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte("your_secret_key"), nil
	})

	if err != nil {
		return false, err
	}

	if !token.Valid {
		return false, errors.New("invalid JWT")
	}

	return true, nil
}

func CheckRole(tokenString string) error {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte("your_secret_key"), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		role := claims["role"].(string)
		if role != "admin" {
			return errors.New("user is not an admin")
		}
	} else {
		return err
	}
	return nil
}
