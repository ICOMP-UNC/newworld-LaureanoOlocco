package handlers

import (
	"github.com/ICOMP-UNC/newworld-LaureanoOlocco/internal/core/domain"
	"github.com/ICOMP-UNC/newworld-LaureanoOlocco/internal/core/services"

	"github.com/gofiber/fiber/v2"
)

// Register godoc
// @Summary Register
// @Description register a new user
// @Tags auth
// @Param Register body domain.Register true "User registration details"
// @Accept  json
// @Produce  json
// @Success 201 {string} string "User added"
// @Failure 500 {string} string "Bad server"
// @Failure 400 {string} string "Bad request"
// @Router /auth/register [post]
func Register(c *fiber.Ctx) error {
	var register domain.Register
	if err := c.BodyParser(&register); err != nil {
		return c.Status(400).SendString("Bad request")
	}

	if err := services.CreateUser(register); err != nil {
		return c.Status(500).SendString("Bad server")
	}

	return c.Status(201).SendString("User added")
}

// Login godoc
// @Summary Login
// @Description login a user
// @Tags auth
// @Accept  json
// @Produce  json
// @Param Login body domain.Login true "User login details"
// @Success 200 {string} string "JWT"
// @Failure 500 {string} string "Bad server"
// @Failure 400 {string} string "Bad request"
// @Router /auth/login [post]
func Login(c *fiber.Ctx) error {
	var login domain.Login
	if err := c.BodyParser(&login); err != nil {
		return c.Status(400).SendString("Bad request")
	}

	token, err := services.LoginUser(login)
	if err != nil {
		return c.Status(500).SendString("Bad server")
	}

	return c.SendString(token)
}
