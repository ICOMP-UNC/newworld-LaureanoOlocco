package handlers

import (
	"github.com/ICOMP-UNC/newworld-LaureanoOlocco/internal/core/domain"
	"github.com/ICOMP-UNC/newworld-LaureanoOlocco/internal/core/ports"
	"github.com/ICOMP-UNC/newworld-LaureanoOlocco/internal/utils"
	"github.com/gofiber/fiber/v2"
)

type OfferHandlers struct {
	offerService ports.IOffersService
}

func NewOfferHandlers(offerService ports.IOffersService) *OfferHandlers {
	return &OfferHandlers{
		offerService: offerService,
	}
}

// @Summary Get all offers
// @Description Retrieve all offers using JWT token authentication
// @Tags Offers
// @Security ApiKeyAuth
// @Accept  json
// @Produce  json
// @Param Authorization header string true "JWT token"
// @Success 200 {object} domain.OffersRegister
// @Failure 400 {object} domain.BadResponse
// @Failure 500 "Bad server"
// @Router /auth/offers [get]
func (h *OfferHandlers) GetOffers(c *fiber.Ctx) error {
	offers, err := h.offerService.GetAllOffers()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(domain.OffersRegister{Code: "200", Message: offers})
}

// Checkout handles the checkout process
// @Summary Checkout an order
// @Description Checkout an order by providing the order details using JWT token authentication
// @Tags Offers
// @Security ApiKeyAuth
// @Accept  json
// @Produce  json
// @Param Authorization header string true "JWT token"
// @Param data body domain.OrderCheckout true "Order details"
// @Success 200 {object} domain.OrderResponse
// @Failure 400 {object} domain.BadResponse
// @Failure 401 "Invalid token"
// @Failure 500 "Bad server"
// @Router /auth/checkout [post]
func (h *OfferHandlers) Checkout(c *fiber.Ctx) error {
	var req domain.OrderCheckout
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(domain.BadResponse{Code: "400", Message: "Bad request"})
	}

	// Extract token using the proper utility function
	tokenString := utils.ExtractToken(c)
	if tokenString == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(domain.BadResponse{
			Code:    "401",
			Message: "Missing or malformed token",
		})
	}

	// Get email from token
	email, err := utils.ExtractEmailFromToken(tokenString)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(domain.BadResponse{
			Code:    "401",
			Message: "Invalid or expired token",
		})
	}

	// Call the service to process the order
	err, order := h.offerService.ProcessOrder(email, req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(domain.BadResponse{
			Code:    "500",
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(domain.OrderResponse{
		Code:    "200",
		Message: order,
	})
}

// Checkout handles the checkout process
// @Summary Get order by ID
// @Description Get an order by providing the order ID using JWT token authentication
// @Tags Offers
// @Security ApiKeyAuth
// @Accept  json
// @Produce  json
// @Param Authorization header string true "JWT token"
// @Param id path string true "Order ID"
// @Success 200 {object} domain.OrderStatusResponse
// @Failure 400 {object} domain.BadResponse
// @Failure 401 "Invalid token"
// @Failure 500 "Bad server"
// @Router /auth/order/{id} [get]
func (h *OfferHandlers) GetOrderById(c *fiber.Ctx) error {

	// Parse the order ID from the URL and trans
	orderID := c.Params("id")

	// Call the service to get the order by ID
	err, order := h.offerService.GetOrderById(orderID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	// build the response using domain.OrderStatus
	return c.Status(fiber.StatusOK).JSON(domain.OrderStatusResponse{Code: "200", Message: order})

}
