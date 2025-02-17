package utils

import (
	"net/http/httptest"
	"os"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/valyala/fasthttp"
)

var mr *miniredis.Miniredis

func TestMain(m *testing.M) {
	var err error
	mr, err = miniredis.Run()
	if err != nil {
		panic(err)
	}

	// Configurar ambiente de prueba
	os.Setenv("REDIS_HOST", mr.Host())
	os.Setenv("REDIS_PORT", mr.Port())
	os.Setenv("JWT_SECRET_KEY", "test_secret_key")
	os.Setenv("ADMIN_EMAIL", "admin@test.com")
	os.Setenv("ADMIN_PASSWORD", "admin123")

	code := m.Run()

	mr.Close()
	os.Exit(code)
}

func TestGenerateJWT(t *testing.T) {
	tests := []struct {
		name      string
		email     string
		password  string
		wantErr   bool
		checkRole string
	}{
		{
			name:      "Valid user token",
			email:     "user@example.com",
			password:  "password123",
			wantErr:   false,
			checkRole: "user",
		},
		{
			name:      "Admin token",
			email:     os.Getenv("ADMIN_EMAIL"),
			password:  os.Getenv("ADMIN_PASSWORD"),
			wantErr:   false,
			checkRole: "admin",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := GenerateJWT(tt.email, tt.password)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Empty(t, token)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, token)

				// Verificar que el token está en Redis
				val, err := RedisClient.Get(ctx, token).Result()
				assert.NoError(t, err)
				assert.Equal(t, tt.email, val)
			}
		})
	}
}

func TestAuthToken(t *testing.T) {
	app := fiber.New()

	app.Get("/protected", AuthToken, func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	tests := []struct {
		name           string
		setup          func() (string, error)
		expectedStatus int
	}{
		{
			name: "Valid token",
			setup: func() (string, error) {
				return GenerateJWT("test@example.com", "password")
			},
			expectedStatus: fiber.StatusOK,
		},
		{
			name: "Missing token",
			setup: func() (string, error) {
				return "", nil
			},
			expectedStatus: fiber.StatusUnauthorized,
		},
		{
			name: "Invalid token",
			setup: func() (string, error) {
				return "invalid.token", nil
			},
			expectedStatus: fiber.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := tt.setup()
			if err != nil {
				t.Fatalf("Setup failed: %v", err)
			}

			req := httptest.NewRequest("GET", "/protected", nil)
			if token != "" {
				req.Header.Set("Authorization", "Bearer "+token)
			}

			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)
		})
	}
}

func TestExtractToken(t *testing.T) {
	tests := []struct {
		name          string
		authorization string
		expected      string
	}{
		{
			name:          "Valid Bearer token",
			authorization: "Bearer token123",
			expected:      "token123",
		},
		{
			name:          "Missing Bearer prefix",
			authorization: "token123",
			expected:      "",
		},
		{
			name:          "Empty authorization",
			authorization: "",
			expected:      "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Crear un contexto de Fiber para la prueba
			app := fiber.New()
			c := app.AcquireCtx(&fasthttp.RequestCtx{})
			defer app.ReleaseCtx(c)

			// Establecer el header de autorización
			c.Request().Header.Set("Authorization", tt.authorization)

			// Probar la función
			result := ExtractToken(c)
			assert.Equal(t, tt.expected, result)
		})
	}
}
