package handlers

import (
	"github.com/ICOMP-UNC/newworld-LaureanoOlocco/internal/core/domain"
	"github.com/ICOMP-UNC/newworld-LaureanoOlocco/internal/core/ports"
	"github.com/gofiber/fiber/v2"
)

type AdminHandlers struct {
	adminService ports.IAdminService
}

func NewAdminHandlers(adminService ports.IAdminService) *AdminHandlers {
	return &AdminHandlers{
		adminService: adminService,
	}
}

// @Summary Get the dashboard data
// @Description Retrieve the dashboard data using JWT token authentication
// @Tags Admin
// @Security ApiKeyAuth
// @Accept  json
// @Produce  json
// @Param Authorization header string true "JWT token"
// @Success 200 {object} domain.DashboardResponse
// @Failure 400 {object} domain.BadResponse
// @Failure 500 "Bad server"
// @Router /admin/dashboard [get]
func (h *AdminHandlers) GetDashboard(c *fiber.Ctx) error {

	// call the service to get the dashboard data
	dashboard, err := h.adminService.GetDashboardData()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(domain.BadResponse{Code: "500", Message: "Bad server error"})
	}

	return c.Status(fiber.StatusOK).JSON(domain.DashboardResponse{Code: "200", Message: dashboard})

}

// @Summary Update the order status
// @Description Update the status of an order using JWT token authentication
// @Tags Admin
// @Security ApiKeyAuth
// @Accept  json
// @Produce  json
// @Param Authorization header string true "JWT token"
// @Param id path int true "Order ID"
// @Param data body domain.OrderStatusUpdate true "Order Status"
// @Success 200 {object} domain.UserOrderStatus
// @Failure 400 {object} domain.BadResponse
// @Response 404 "Order not found"
// @Failure 500 "Bad server"
// @Router /admin/orders/{id} [patch]
func (h *AdminHandlers) UpdateOrderStatus(c *fiber.Ctx) error {

	//get the order id from the url
	orderID := c.Params("id")
	if orderID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(domain.BadResponse{Code: "400", Message: "Bad request"})
	}

	// get the status from the body
	var status domain.OrderStatusUpdate
	if err := c.BodyParser(&status); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(domain.BadResponse{Code: "400", Message: "Bad request"})
	}

	//call the service to update the order status
	orderUpdated, err := h.adminService.UpdateOrderStatus(orderID, status.Status)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(domain.BadResponse{Code: "500", Message: "Bad server error"})
	}

	return c.Status(fiber.StatusOK).JSON(domain.OrderStatusResponse{Code: "200", Message: orderUpdated})

}
