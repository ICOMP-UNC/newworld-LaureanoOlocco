package ports

import (
	"github.com/ICOMP-UNC/newworld-LaureanoOlocco/internal/core/domain"
	"github.com/gofiber/fiber/v2"
)

type IUserService interface {
	Login(email string, password string) error
	Register(username, email, password, passwordConfirmation string) error
	GetUserByEmail(email string) (*domain.User, error)
}

type IUserRepository interface {
	Login(email string, password string) error
	Register(username, email string, password string) error
	GetUserByEmail(email string) (*domain.User, error)
}

type IUserHandlers interface {
	Login(c *fiber.Ctx) error
	Register(c *fiber.Ctx) error
	GetUserByEmail(c *fiber.Ctx) error
}

type IOffersService interface {
	GetAllOffers() ([]domain.OfferWithPrice, error)
	ProcessOrder(email string, order domain.OrderCheckout) (error, domain.Order)
	GetOrderById(orderID string) (error, domain.UserOrderStatus)
	GetAllOrders() ([]domain.UserOrderStatus, error)
}

type IOffersRepository interface {
	GetOffersData() ([]domain.Offer, error)
	InsertOrder(userID, total int) error
	GetOrderById(id string) (int, string, string, int, error)
	GetAllOrders() ([]domain.UserOrderStatus, error)
}

type IOffersHandlers interface {
	GetOffers(c *fiber.Ctx) error
	Checkout(c *fiber.Ctx) error
	GetOrderById(c *fiber.Ctx) error
}

type IAdminHandlers interface {
	GetDashboard(c *fiber.Ctx) error
	UpdateOrderStatus(c *fiber.Ctx) error
}

type IAdminService interface {
	GetDashboardData() (domain.Dashboard, error)
	UpdateOrderStatus(orderID, Newstatus string) (domain.UserOrderStatus, error)
}

type IAdminRepository interface {
	UpdateOrderStatus(orderID, Newstatus string) (domain.UserOrderStatus, error)
}

type IServer interface {
	Initialize()
}
