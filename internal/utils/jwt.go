package utils

import (
	"errors"
	"log"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
)

// var JwtKey = []byte("yJ42jhCACeBRZsXiRi22qZSTnn1xbVev6ybirXaoYS8=")
var JwtKey = []byte(os.Getenv("JWT_SECRET_KEY"))

func init() {
	// Validar que la clave JWT esté configurada
	if len(JwtKey) == 0 {
		log.Fatalf("JWT_SECRET_KEY is not set in the environment")
	}
}

type Claims struct {
	Email string `json:"email"`
	Role  string `json:"role"`
	jwt.StandardClaims
}

// GenerateJWTFunc is a function type for generating JWT tokens
type GenerateJWTFunc func(email, password string) (string, error)

// GenerateJWT generates a JWT token with the provided email and role
var GenerateJWT GenerateJWTFunc = func(email, password string) (string, error) {

	// Set the expiration time of the token
	expirationTime := time.Now().Add(99 * time.Minute)

	// set role based on username and password
	role := "user"

	// the email is unique, so we can use it as a username and only one admin can have the email
	if email == "Ubuntu@gmail.com" && password == "Ubuntu" {
		role = "admin"
	}

	// Create the JWT claims, which includes the email and expiry time
	claims := &Claims{
		Email: email,
		Role:  role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	// Declare the token with the algorithm used for signing, and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(JwtKey)
	if err != nil {
		return "", errors.New("error generating token")
	}

	return tokenString, nil

}

// VerifyJWT verifies the provided token string and returns the claims if the token is valid for use in the handler
func AuthToken(c *fiber.Ctx) error {
	tokenString := c.Get("Authorization")
	if tokenString == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Missing or malformed JWT"})
	}

	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return JwtKey, nil
	})

	if err != nil || !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid or expired JWT"})
	}

	c.Locals("claims", claims)
	return c.Next()
}

// cheking if the token is valid and if the role is admin
func AuthAdminToken(c *fiber.Ctx) error {
	tokenString := c.Get("Authorization")
	if tokenString == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Missing or malformed JWT"})
	}

	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return JwtKey, nil
	})

	if err != nil || !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid or expired JWT"})
	}

	if claims.Role != "admin" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Insufficient permissions"})
	}

	c.Locals("claims", claims)
	return c.Next()
}

// ExtractEmailFromToken extracts the email from a JWT token
func ExtractEmailFromToken(tokenStr string) (string, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return JwtKey, nil
	})

	if err != nil || !token.Valid {
		return "", errors.New("invalid or expired JWT")
	}

	return claims.Email, nil
}
