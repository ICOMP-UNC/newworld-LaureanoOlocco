package server

import (
	"log"
	"os"

	"github.com/ICOMP-UNC/newworld-LaureanoOlocco/internal/core/ports"
	"github.com/ICOMP-UNC/newworld-LaureanoOlocco/internal/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/swagger"
)

type Server struct {
	userHandlers   ports.IUserHandlers
	offerHandlers  ports.IOffersHandlers
	addminHandlers ports.IAdminHandlers
}

func NewServer(userHandlers ports.IUserHandlers, OfferHandlers ports.IOffersHandlers, AdminHandlers ports.IAdminHandlers) *Server {
	return &Server{
		userHandlers:   userHandlers,
		offerHandlers:  OfferHandlers,
		addminHandlers: AdminHandlers,
	}
}

func (s *Server) Start() {
	app := fiber.New()

	// Use logger middleware
	app.Use(cors.New())
	app.Use(logger.New())

	// Endpoints related to everyone
	userRoutes := app.Group("/user")
	userRoutes.Post("/register", s.userHandlers.Register)
	userRoutes.Post("/login", s.userHandlers.Login)
	userRoutes.Get("/search", s.userHandlers.GetUserByEmail)

	// Endpoints related to offers
	authRoutes := app.Group("/auth")
	authRoutes.Use(utils.AuthToken)
	authRoutes.Get("/offers", s.offerHandlers.GetOffers)
	authRoutes.Post("/checkout", s.offerHandlers.Checkout)
	authRoutes.Get("/order/:id", s.offerHandlers.GetOrderById)

	// Endpoints related to admin
	adminRoutes := app.Group("/admin")
	// adminRoutes.Use(utils.AuthAdminToken)
	adminRoutes.Get("/dashboard", s.addminHandlers.GetDashboard)
	adminRoutes.Patch("/orders/:id", s.addminHandlers.UpdateOrderStatus)

	// Use swagger middleware
	app.Get("/swagger/*", swagger.HandlerDefault)

	// Start the server
	port := os.Getenv("API_PORT")
	log.Printf("Server running on port %s", port)
	if err := app.Listen(":" + port); err != nil {
		log.Fatalf("Error starting server: %q", err)
	}
}
