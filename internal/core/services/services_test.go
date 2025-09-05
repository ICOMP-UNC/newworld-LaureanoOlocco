package services

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/ICOMP-UNC/newworld-LaureanoOlocco/internal/core/domain"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Login(email, password string) error {
	args := m.Called(email, password)
	return args.Error(0)
}

func (m *MockUserRepository) Register(username, email, password string) error {
	args := m.Called(username, email, password)
	return args.Error(0)
}

func (m *MockUserRepository) GetUserByEmail(email string) (*domain.User, error) {
	args := m.Called(email)
	user := args.Get(0)
	if user == nil {
		return nil, args.Error(1)
	}
	return user.(*domain.User), args.Error(1)
}

func TestUserService_Login(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)

	t.Run("successful login", func(t *testing.T) {
		mockRepo.On("Login", "test@example.com", "password").Return(nil)

		err := service.Login("test@example.com", "password")
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("login with wrong password", func(t *testing.T) {
		mockRepo.On("Login", "test@example.com", "wrongpassword").Return(errors.New("invalid password"))

		err := service.Login("test@example.com", "wrongpassword")
		assert.Error(t, err)
		assert.Equal(t, "invalid password", err.Error())
		mockRepo.AssertExpectations(t)
	})
}

func TestUserService_Register(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)

	t.Run("successful register", func(t *testing.T) {
		mockRepo.On("Register", "testuser", "test@example.com", "password").Return(nil).Once()

		err := service.Register("testuser", "test@example.com", "password", "password")
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("register with non-matching passwords", func(t *testing.T) {
		err := service.Register("testuser", "test@example.com", "password", "differentpassword")
		assert.Error(t, err)
		assert.Equal(t, "the passwords are not equal", err.Error())
	})

	t.Run("register with existing email", func(t *testing.T) {
		mockRepo.On("Register", "testuser", "test@example.com", "password").Return(errors.New("email already registered")).Once()

		err := service.Register("testuser", "test@example.com", "password", "password")
		assert.Error(t, err)
		assert.Equal(t, "email already registered", err.Error())
		mockRepo.AssertExpectations(t)
	})
}

func TestUserService_GetUserByEmail(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)

	t.Run("get existing user", func(t *testing.T) {
		mockUser := &domain.User{ID: 1, Email: "test@example.com", Password: "password"}
		mockRepo.On("GetUserByEmail", "test@example.com").Return(mockUser, nil)

		user, err := service.GetUserByEmail("test@example.com")
		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, "test@example.com", user.Email)
		mockRepo.AssertExpectations(t)
	})

	t.Run("get non-existing user", func(t *testing.T) {
		mockRepo.On("GetUserByEmail", "nonexistent@example.com").Return(nil, errors.New("user not found"))

		user, err := service.GetUserByEmail("nonexistent@example.com")
		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Equal(t, "user not found", err.Error())
		mockRepo.AssertExpectations(t)
	})
}

// Mock del repositorio de ofertas
type MockOfferRepository struct {
	mock.Mock
}

func (m *MockOfferRepository) GetOffersData() ([]domain.Offer, error) {
	args := m.Called()
	return args.Get(0).([]domain.Offer), args.Error(1)
}

func (m *MockOfferRepository) InsertOrder(userID, total int) error {
	args := m.Called(userID, total)
	return args.Error(0)
}

func (m *MockOfferRepository) GetOrderById(id string) (int, string, string, int, error) {
	args := m.Called(id)
	return args.Int(0), args.String(1), args.String(2), args.Int(3), args.Error(4)
}

func (m *MockOfferRepository) GetAllOrders() ([]domain.UserOrderStatus, error) {
	args := m.Called()
	return args.Get(0).([]domain.UserOrderStatus), args.Error(1)
}

// Mock del servicio de usuario
type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) Login(email, password string) error {
	args := m.Called(email, password)
	return args.Error(0)
}

func (m *MockUserService) Register(username, email, password string) error {
	args := m.Called(username, email, password)
	return args.Error(0)
}

func (m *MockUserService) GetUserByEmail(email string) (*domain.User, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func TestGetAllOffers(t *testing.T) {
	// Configurar el mock del repositorio
	mockRepo := new(MockOfferRepository)

	// Definir las ofertas mockeadas
	offers := []domain.Offer{
		{ID: 1, Name: "meat", Price: 19.0, Category: "food"},
		{ID: 2, Name: "vegetables", Price: 18.0, Category: "food"},
		{ID: 3, Name: "fruits", Price: 3.0, Category: "food"},
		{ID: 4, Name: "water", Price: 20.0, Category: "drink"},
		{ID: 5, Name: "antibiotics", Price: 9.0, Category: "medicine"},
		{ID: 6, Name: "analgesics", Price: 4.0, Category: "medicine"},
		{ID: 7, Name: "bandages", Price: 7.0, Category: "medicine"},
		{ID: 8, Name: "pistol ammo", Price: 3.0, Category: "ammo"},
		{ID: 9, Name: "rifle ammo", Price: 2.0, Category: "ammo"},
		{ID: 10, Name: "shotgun ammo", Price: 12.0, Category: "ammo"},
	}

	// Configurar el mock para que devuelva estas ofertas
	mockRepo.On("GetOffersData").Return(offers, nil)

	// Configurar httpmock
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// Configurar la respuesta simulada para la solicitud a "http://cppserver:9004/supplies"
	supplies := map[string]map[string]int{
		"food": {
			"meat":       1906,
			"vegetables": 2001,
			"fruits":     589,
			"water":      900,
		},
		"medicine": {
			"antibiotics": 304,
			"analgesics":  179,
			"bandages":    900,
		},
		"ammo": {
			"pistol ammo":  702,
			"rifle ammo":   1004,
			"shotgun ammo": 315,
		},
	}
	jsonData, _ := json.Marshal(supplies)

	httpmock.RegisterResponder("GET", "http://cppserver:8083/supplies",
		httpmock.NewStringResponder(200, string(jsonData)))

	// Crear instancia del mock del servicio de usuario
	mockUserService := new(MockUserService)

	// Crear instancia del servicio OfferService
	offerService := NewOfferService(mockRepo, userServiceFromMock(mockUserService))

	// Ejecutar la función que queremos probar
	offersWithPrice, err := offerService.GetAllOffers()
	if err != nil {
		t.Fatalf("error al obtener ofertas: %v", err)
	}

	// Verificar resultados
	assert.NotNil(t, offersWithPrice)
	assert.Equal(t, 10, len(offersWithPrice))

	// Asegurarnos de que se llame a GetOffersData() del mockRepo
	mockRepo.AssertExpectations(t)
}

func TestOfferService_ProcessOrder_Success(t *testing.T) {
	// Mock the OfferRepository
	mockRepo := new(MockOfferRepository)

	// Mock the UserService
	mockUserService := new(MockUserService)

	// Create the service with the mocks
	offerService := NewOfferService(mockRepo, userServiceFromMock(mockUserService))

	// Define the test email and order
	testEmail := "test@example.com"
	testOrder := domain.OrderCheckout{
		Order: []domain.Items{
			{ItemID: 1, Quantity: 2},
			{ItemID: 2, Quantity: 3},
		},
	}

	// Mock the expected calls
	mockUser := &domain.User{ID: 1, Email: testEmail}
	mockUserService.On("GetUserByEmail", testEmail).Return(mockUser, nil)
	mockRepo.On("InsertOrder", mockUser.ID, 5).Return(nil)

	// Call the function
	err, orderResponse := offerService.ProcessOrder(testEmail, testOrder)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, 5, orderResponse.Total)
	assert.Equal(t, "pending", orderResponse.Status)

	// Verify the expectations
	mockUserService.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
}

// Función para convertir el mock de UserService a UserService
func userServiceFromMock(mockUserService *MockUserService) *UserService {
	return &UserService{
		mockUserService,
	}
}
