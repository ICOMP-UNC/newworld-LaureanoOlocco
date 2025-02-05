package handlers

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/ICOMP-UNC/newworld-LaureanoOlocco/internal/core/domain"
	"github.com/ICOMP-UNC/newworld-LaureanoOlocco/internal/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var JwtKey []byte

func init() {
	JwtKey = []byte(os.Getenv("JWT_SECRET_KEY"))
}

// Mock del servicio de usuario
type mockUserService struct {
	mock.Mock
}

func (m *mockUserService) Login(email string, password string) error {
	args := m.Called(email, password)
	return args.Error(0)
}

func (m *mockUserService) GetUserByEmail(email string) (*domain.User, error) {
	args := m.Called(email)
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *mockUserService) Register(username, email, password, confirmPass string) error {
	args := m.Called(username, email, password, confirmPass)
	return args.Error(0)
}

type mockOfferService struct {
	mock.Mock
}

func (m *mockOfferService) GetAllOffers() ([]domain.OfferWithPrice, error) {
	args := m.Called()
	return args.Get(0).([]domain.OfferWithPrice), args.Error(1)
}

func (m *mockOfferService) GetOrderById(string) (error, domain.UserOrderStatus) {
	args := m.Called()
	return args.Error(0), args.Get(1).(domain.UserOrderStatus)
}

func (m *mockOfferService) GetAllOrders() ([]domain.UserOrderStatus, error) {
	args := m.Called()
	return args.Get(0).([]domain.UserOrderStatus), args.Error(1)
}

func (m *mockOfferService) ProcessOrder(email string, order domain.OrderCheckout) (error, domain.Order) {
	args := m.Called(email, order)
	return args.Error(0), args.Get(1).(domain.Order)
}

//-----------------------------------------------------------USER HANDLERS----------------------------------------------------------------//

// TestLogin tests the login function
func TestLogin(t *testing.T) {

	// Establece el valor de la variable de entorno para este test
	os.Setenv("JWT_SECRET_KEY", "clave_secreta")
	defer os.Unsetenv("JWT_SECRET_KEY") // Limpia la variable de entorno después del test

	// Inicializa la aplicación Fiber
	app := fiber.New()

	// Inicializa el mock del servicio de usuario
	mockService := new(mockUserService)

	// Inicializa los handlers de usuario
	handlers := &UserHandlers{
		userService: mockService,
	}

	// Mock para la generación de JWT
	originalGenerateJWT := utils.GenerateJWT
	utils.GenerateJWT = func(email, password string) (string, error) {
		return "token", nil
	}
	defer func() {
		// Restaura la función GenerateJWT original después del test
		utils.GenerateJWT = originalGenerateJWT
	}()

	// Configura las rutas
	app.Post("/user/login", handlers.Login)

	// Caso de prueba: Successful login
	t.Run("Successful login", func(t *testing.T) {
		mockService.On("Login", "Ubuntu@gmail.com", "Ubuntu").Return(nil)

		// Crea una solicitud
		reqBody := `{"email": "Ubuntu@gmail.com", "password": "Ubuntu"}`
		req := httptest.NewRequest(http.MethodPost, "/user/login", bytes.NewBuffer([]byte(reqBody)))
		req.Header.Set("Content-Type", "application/json")

		// Crea un nuevo response recorder
		resp, err := app.Test(req, -1) // -1 para desactivar el límite del cuerpo de la respuesta
		if err != nil {
			t.Fatalf("Error while testing the login function: %s", err)
		}

		// Verifica el código de estado de la respuesta
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		// Verifica el cuerpo de la respuesta
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("Error while reading the response body: %v", err)
		}
		expectedBody := `{"code":"200","token":"token"}`
		assert.JSONEq(t, expectedBody, string(body))

		// Verifica las expectativas del mock
		mockService.AssertExpectations(t)
	})

	// Caso de prueba: User not found
	t.Run("User not found", func(t *testing.T) {
		mockService.On("Login", "unknown@gmail.com", "password").Return(errors.New("user not found"))

		// Crea una solicitud
		reqBody := `{"email": "unknown@gmail.com", "password": "password"}`
		req := httptest.NewRequest(http.MethodPost, "/user/login", bytes.NewBuffer([]byte(reqBody)))
		req.Header.Set("Content-Type", "application/json")

		// Crea un nuevo response recorder
		resp, err := app.Test(req, -1) // -1 para desactivar el límite del cuerpo de la respuesta
		if err != nil {
			t.Fatalf("Error while testing the login function: %s", err)
		}

		// Verifica el código de estado de la respuesta
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)

		// Verifica el cuerpo de la respuesta
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("Error while reading the response body: %v", err)
		}
		expectedBody := `{"code":"404","message":"User not found"}`
		assert.JSONEq(t, expectedBody, string(body))

		// Verifica las expectativas del mock
		mockService.AssertExpectations(t)
	})
}

// TestRegister tests the register function
func TestRegister(t *testing.T) {
	// Inicializa la aplicación Fiber
	app := fiber.New()

	// Inicializa el mock del servicio de usuario
	mockService := new(mockUserService)

	// Inicializa los handlers de usuario
	handlers := &UserHandlers{
		userService: mockService,
	}

	// Configura las rutas
	app.Post("/user/register", handlers.Register)

	// Caso de prueba: Successful registration
	t.Run("Successful registration", func(t *testing.T) {
		mockService.On("Register", "testuser", "testuser@example.com", "password123", "password123").Return(nil)

		// Crea una solicitud
		reqBody := `{"username": "testuser", "email": "testuser@example.com", "password": "password123", "confirmPassword": "password123"}`
		req := httptest.NewRequest(http.MethodPost, "/user/register", bytes.NewBuffer([]byte(reqBody)))
		req.Header.Set("Content-Type", "application/json")

		// Crea un nuevo response recorder
		resp, err := app.Test(req, -1) // -1 para desactivar el límite del cuerpo de la respuesta
		if err != nil {
			t.Fatalf("Error while testing the register function: %s", err)
		}

		// Verifica el código de estado de la respuesta
		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		// Verifica el cuerpo de la respuesta
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("Error while reading the response body: %v", err)
		}
		expectedBody := `{"code":"201","message":"User added"}`
		assert.JSONEq(t, expectedBody, string(body))

		// Verifica las expectativas del mock
		mockService.AssertExpectations(t)
	})

	// Caso de prueba: Passwords are not equal
	t.Run("Passwords are not equal", func(t *testing.T) {
		mockService.On("Register", "testuser", "testuser@example.com", "password123", "password124").Return(errors.New("the passwords are not equal"))

		// Crea una solicitud
		reqBody := `{"username": "testuser", "email": "testuser@example.com", "password": "password123", "confirmPassword": "password124"}`
		req := httptest.NewRequest(http.MethodPost, "/user/register", bytes.NewBuffer([]byte(reqBody)))
		req.Header.Set("Content-Type", "application/json")

		// Crea un nuevo response recorder
		resp, err := app.Test(req, -1) // -1 para desactivar el límite del cuerpo de la respuesta
		if err != nil {
			t.Fatalf("Error while testing the register function: %s", err)
		}

		// Verifica el código de estado de la respuesta
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		// Verifica el cuerpo de la respuesta
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("Error while reading the response body: %v", err)
		}
		expectedBody := `{"code":"401","message":"the passwords are not equal"}`
		assert.JSONEq(t, expectedBody, string(body))

		// Verifica las expectativas del mock
		mockService.AssertExpectations(t)
	})
}

// -----------------------------------------------------------OFFER HANDLERS ---------------------------------------------------------------//
func TestGetOffers(t *testing.T) {
	// Inicializa la aplicación Fiber
	app := fiber.New()

	// Inicializa el mock del servicio de oferta
	mockService := new(mockOfferService)

	// Inicializa los handlers de oferta
	handlers := &OfferHandlers{
		offerService: mockService,
	}

	// Configura las rutas
	app.Get("/auth/offers", handlers.GetOffers)

	// Caso de prueba: Successful retrieval of offers
	t.Run("Successful retrieval of offers", func(t *testing.T) {
		mockOffers := []domain.OfferWithPrice{
			{ID: 1, Name: "meat", Quantity: 381, Price: 19, Category: "food"},
			{ID: 2, Name: "vegetables", Quantity: 400, Price: 18, Category: "food"},
		}
		mockService.On("GetAllOffers").Return(mockOffers, nil)

		// Crea una solicitud
		req := httptest.NewRequest(http.MethodGet, "/auth/offers", nil)
		req.Header.Set("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6ImV4YW1wbGVAZXhhbXBsZS5jb20iLCJyb2xlIjoidXNlciIsImV4cCI6MTcxODYzNzc4NH0.uLikpxnDlJ6Qu8nS2JjgE4mbcejQddIbgmm7Gn6LgJg")

		// Crea un nuevo response recorder
		resp, err := app.Test(req, -1) // -1 para desactivar el límite del cuerpo de la respuesta
		if err != nil {
			t.Fatalf("Error while testing the GetOffers function: %s", err)
		}

		// Verifica el código de estado de la respuesta
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		// Verifica el cuerpo de la respuesta
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("Error while reading the response body: %v", err)
		}
		expectedBody := `{"code":"200","message":[{"id":1,"name":"meat","quantity":381,"price":19,"category":"food"},{"id":2,"name":"vegetables","quantity":400,"price":18,"category":"food"}]}`
		assert.JSONEq(t, expectedBody, string(body))

		// Verifica las expectativas del mock
		mockService.AssertExpectations(t)
	})
}
