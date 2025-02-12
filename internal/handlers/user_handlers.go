package handlers

import (
	"fmt"
	"log"
	"time"

	"github.com/ICOMP-UNC/newworld-LaureanoOlocco/internal/core/domain"
	"github.com/ICOMP-UNC/newworld-LaureanoOlocco/internal/core/ports"
	"github.com/ICOMP-UNC/newworld-LaureanoOlocco/internal/utils"
	"github.com/gofiber/fiber/v2"
)

// UserHandlers is a struct to represent the user handlers
type UserHandlers struct {
	userService ports.IUserService
}

// NewUserHandlers creates a new UserHandlers instance with the provided user service
func NewUserHandlers(service ports.IUserService) *UserHandlers {
	return &UserHandlers{userService: service}
}

// @Summary User login
// @Description Login with email and password
// @Tags Users
// @Accept  json
// @Produce  json
// @Param   userLogin     body    domain.UserLogin     true        "User login details"
// @Success 200 {object} domain.UserResponseLogin
// @Response 400 {object} domain.BadResponse
// @Response 404 "User not found"
// @Response 401 "Invalid password"
// @Response 500 "Bad server"
// @Router /user/login [post]
func (h *UserHandlers) Login(c *fiber.Ctx) error {
	var req domain.UserLogin
	ip := c.IP()

	// Parse the body into the req struct
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(domain.BadResponse{Code: "400", Message: "Bad request"})
	}

	// if the request is fine, we call the login function from the user service
	err := h.userService.Login(req.Email, req.Password)
	if err != nil {
		// Only check rate limit for failed login attempts
		if err.Error() == "invalid password" {
			allowed, err := utils.CheckLoginRateLimit(ip)
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(domain.BadResponse{
					Code:    "500",
					Message: "Internal server error",
				})
			}

			if !allowed {
				_, reset, _ := utils.GetRateLimitInfo(ip, "login")
				remainingMinutes := int((reset - time.Now().Unix()) / 60)
				return c.Status(fiber.StatusTooManyRequests).JSON(domain.BadResponse{
					Code:    "429",
					Message: fmt.Sprintf("Too many failed attempts. Try again in %d minutes", remainingMinutes),
				})
			}
		}

		if err.Error() == "user not found" {
			return c.Status(fiber.StatusNotFound).JSON(domain.BadResponse{Code: "404", Message: "User not found"})
		}
		if err.Error() == "invalid password" {
			return c.Status(fiber.StatusUnauthorized).JSON(domain.BadResponse{Code: "401", Message: "Invalid password"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(domain.BadResponse{Code: "500", Message: "Internal server error"})
	}

	// if the login is successful, we generate a token
	token, err := utils.GenerateJWT(req.Email, req.Password)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(domain.BadResponse{Code: "500", Message: "Internal server error"})
	}

	// In case of successful login, clear failed attempts counter
	if err := utils.ClearLoginRateLimit(ip); err != nil {
		log.Printf("Error clearing rate limit: %v", err)
	}

	// if the token is generated successfully, we return the token
	return c.JSON(domain.UserResponseLogin{Code: "200", Token: token})
}

// @Summary User registration
// @Description Register with username, email, password, and confirm password
// @Tags Users
// @Accept  json
// @Produce  json
// @Param   userRegister    body    domain.UserRegister     true        "User registration details"
// @Router /user/register [post]
// @Success 201 {object} domain.Response
// @Response 400 {object} domain.BadResponse
// @Response 401 "The passwords are not equal"
// @Response 402 "Mail already registered"
// @Response 500 "Bad server"
func (h *UserHandlers) Register(c *fiber.Ctx) error {
	var req domain.UserRegister
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(domain.BadResponse{Code: "400", Message: "Bad request"})
	}

	err := h.userService.Register(req.Username, req.Email, req.Password, req.ConfirmPassword)
	if err != nil {
		log.Println(err)
		if err.Error() == "the passwords are not equal" {
			return c.Status(fiber.StatusBadRequest).JSON(domain.BadResponse{Code: "401", Message: err.Error()})
		} else if err.Error() == "email already registered" {
			return c.Status(fiber.StatusBadRequest).JSON(domain.BadResponse{Code: "402", Message: err.Error()})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(domain.BadResponse{Code: "500", Message: "Bad server"})
	}

	return c.Status(fiber.StatusCreated).JSON(domain.Response{Code: "201", Message: "User added"})
}

// @Summary Get user by email
// @Description Get user details by email
// @Tags Users
// @Accept  json
// @Produce  json
// @Param   email     query    string     true        "Email"
// @Success 200 {object} domain.User
// @Response 404 "Not found"
// @Response 500 "Bad server"
// @Router /user/search [get]
func (h *UserHandlers) GetUserByEmail(c *fiber.Ctx) error {
	email := c.Query("email")

	user, err := h.userService.GetUserByEmail(email)

	if err != nil {
		return c.JSON(fiber.Map{"error": err.Error()})
	}

	if user == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "user not found"})
	}

	return c.JSON(domain.User{
		ID:       user.ID,
		Email:    user.Email,
		Password: user.Password,
	})
}
