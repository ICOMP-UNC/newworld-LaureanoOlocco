package utils

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
)

// Añade estas constantes junto con las otras en utils.go
const (
	MaxLoginAttempts = 3 // Máximo número de intentos de login
	LoginWindow      = 2 // Ventana de tiempo en minutos para login
)

var (
	AdminEmail    = os.Getenv("ADMIN_EMAIL")
	AdminPassword = os.Getenv("ADMIN_PASSWORD")
)

var JwtKey = []byte(os.Getenv("JWT_SECRET_KEY"))
var RedisClient *redis.Client
var ctx = context.Background()

func init() {
	log.Printf("Admin email configured: %s", AdminEmail)
	log.Printf("Admin password configured: %s", AdminPassword)

	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")
	if redisPort == "" {
		redisPort = "6379"
	}

	log.Printf("Attempting to connect to Redis at %s:%s", redisHost, redisPort)

	// Configurar Redis
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     redisHost + ":" + redisPort,
		Password: "",
		DB:       0,
	})

	// Probar conexión con Redis
	_, err := RedisClient.Ping(ctx).Result()
	if err != nil {
		log.Printf("Failed to connect to Redis: %v", err)
		log.Fatalf("Redis connection failed: %v", err)
	}
	log.Printf("Successfully connected to Redis")
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

	expirationTime := time.Now().Add(99 * time.Minute)
	role := "user"

	if IsAdmin(email, password) {
		log.Printf("Admin role assigned")
		role = "admin"
	}

	claims := &Claims{
		Email: email,
		Role:  role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(JwtKey)
	if err != nil {
		log.Printf("Error signing token: %v", err)
		return "", errors.New("error generating token")
	}

	err = RedisClient.Set(ctx, tokenString, email, 99*time.Minute).Err()
	if err != nil {
		return "", errors.New("failed to store JWT in Redis")
	}

	return tokenString, nil
}

func AuthToken(c *fiber.Ctx) error {
	tokenString := ExtractToken(c)
	log.Printf("Token received: %s", tokenString) // Debug del token recibido

	if tokenString == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Missing or malformed JWT"})
	}

	// Debug Redis check
	redisValue, err := RedisClient.Get(ctx, tokenString).Result()
	if err == redis.Nil {
		log.Printf("Token not found in Redis") // Debug Redis
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Token not found in Redis"})
	} else if err != nil {
		log.Printf("Redis error: %v", err) // Debug Redis error
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Redis error: " + err.Error()})
	}
	log.Printf("Redis value found: %s", redisValue) // Debug Redis value

	// Debug JWT parsing
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return JwtKey, nil
	})
	if err != nil {
		log.Printf("JWT Parse error: %v", err) // Debug JWT parsing
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "JWT Parse error: " + err.Error()})
	}
	if !token.Valid {
		log.Printf("Token is invalid") // Debug token validity
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Token is invalid"})
	}

	log.Printf("Token validated successfully for email: %s", claims.Email) // Debug success
	c.Locals("claims", claims)
	return c.Next()
}

// cheking if the token is valid and if the role is admin
func AuthAdminToken(c *fiber.Ctx) error {
	tokenString := ExtractToken(c)
	if tokenString == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Missing or malformed JWT"})
	}

	// Verificar si el token está en Redis
	_, err := RedisClient.Get(ctx, tokenString).Result()
	if err == redis.Nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Token expired or not found"})
	} else if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Redis error"})
	}

	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return JwtKey, nil
	})

	if err != nil || !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid JWT"})
	}

	// Verificar el rol del usuario
	if claims.Role != "admin" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Insufficient permissions"})
	}

	c.Locals("claims", claims)
	return c.Next()
}

// ExtractEmailFromToken extracts the email from a JWT token
func ExtractEmailFromToken(tokenStr string) (string, error) {
	// Verificar si el token está en Redis
	_, err := RedisClient.Get(ctx, tokenStr).Result()
	if err == redis.Nil {
		return "", errors.New("invalid or expired jwt") // Cambié "Invalid" a minúscula
	} else if err != nil {
		return "", errors.New("redis error") // Cambié "Redis" a minúscula
	}

	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return JwtKey, nil
	})

	if err != nil || !token.Valid {
		return "", errors.New("invalid or expired jwt") // Cambié "Invalid" a minúscula
	}

	return claims.Email, nil
}

func Logout(c *fiber.Ctx) error {
	tokenString := ExtractToken(c)
	if tokenString == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Missing or malformed JWT"})
	}

	// Eliminar token de Redis
	err := RedisClient.Del(ctx, tokenString).Err()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to logout"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Logged out successfully"})
}

func ExtractToken(c *fiber.Ctx) string {
	bearerToken := c.Get("Authorization")
	log.Printf("Authorization header: %s", bearerToken) // Debug header completo
	if bearerToken == "" {
		return ""
	}
	parts := strings.Split(bearerToken, " ")
	if len(parts) == 2 {
		return parts[1]
	}
	return ""
}

// CheckLoginRateLimit verifica los intentos de login para una IP específica
func CheckLoginRateLimit(ip string) (bool, error) {
	key := fmt.Sprintf("rate_limit:login:%s", ip)

	// Incrementar el contador
	count, err := RedisClient.Incr(ctx, key).Result()
	if err != nil {
		return false, err
	}

	// Si es el primer intento, establecer la expiración
	if count == 1 {
		err = RedisClient.Expire(ctx, key, LoginWindow*time.Minute).Err()
		if err != nil {
			return false, err
		}
	}

	// Verificar si se excedió el límite
	return count <= MaxLoginAttempts, nil
}

// GetRateLimitInfo obtiene información sobre el rate limit actual
func GetRateLimitInfo(ip string, prefix string) (remaining int64, reset int64, err error) {
	key := fmt.Sprintf("rate_limit:%s:%s", prefix, ip)

	// Obtener el contador actual
	count, err := RedisClient.Get(ctx, key).Int64()
	if err == redis.Nil {
		return MaxLoginAttempts, 0, nil
	}
	if err != nil {
		return 0, 0, err
	}

	// Obtener el tiempo restante
	ttl, err := RedisClient.TTL(ctx, key).Result()
	if err != nil {
		return 0, 0, err
	}

	remaining = MaxLoginAttempts - count
	if remaining < 0 {
		remaining = 0
	}

	reset = time.Now().Add(ttl).Unix()

	return remaining, reset, nil
}

func ClearLoginRateLimit(ip string) error {
	key := fmt.Sprintf("rate_limit:login:%s", ip)
	return RedisClient.Del(ctx, key).Err()
}

func IsAdmin(email, password string) bool {
	return email == AdminEmail && password == AdminPassword
}
