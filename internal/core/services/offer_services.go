package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/ICOMP-UNC/newworld-LaureanoOlocco/internal/core/domain"
	"github.com/ICOMP-UNC/newworld-LaureanoOlocco/internal/core/ports"
)

type OfferService struct {
	offerRepository ports.IOffersRepository
	UserService     *UserService // Add UserService to OfferService
}

type ServerResponse struct {
	Supplies map[string]map[string]int `json:"supplies"`
}

func NewOfferService(repository ports.IOffersRepository, userService *UserService) *OfferService {
	return &OfferService{
		offerRepository: repository,
		UserService:     userService,
	}
}

// func (s *OfferService) GetAllOffers() ([]domain.OfferWithPrice, error) {
// 	// Obtener las ofertas desde el repositorio
// 	offers, err := s.offerRepository.GetOffersData()
// 	if err != nil {
// 		return nil, err
// 	}

// 	// Hacer la solicitud HTTP al servidor cppServer para obtener la cantidad de ofertas

// 	var resp *http.Response

// 	if os.Getenv("RUN_LOCAL") == "true" {
// 		resp, err = http.Get("http://localhost:8083/supplies")
// 		if err != nil {
// 			return nil, err
// 		}
// 	} else {
// 		resp, err = http.Get("http://cppserver:8083/supplies")
// 		if err != nil {
// 			return nil, err
// 		}
// 	}

// 	defer resp.Body.Close()

// 	// Leer el cuerpo de la respuesta
// 	body, err := io.ReadAll(resp.Body)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// Crear el nuevo JSON con la información en el formato deseado
// 	var supplies map[string]map[string]int
// 	err = json.Unmarshal(body, &supplies)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// Crear las ofertas con cantidad
// 	var offersWithPrice []domain.OfferWithPrice
// 	for _, offer := range offers {
// 		quantity := 0
// 		// Normalizamos las categorías para asegurarnos de mapear correctamente
// 		var normalizedCategory string
// 		switch offer.Category {
// 		case "food":
// 			normalizedCategory = "food"
// 		case "drink":
// 			normalizedCategory = "food" // Asumiendo que "water" está bajo "food" en suministros
// 		default:
// 			normalizedCategory = offer.Category
// 		}

// 		if categorySupplies, ok := supplies[normalizedCategory]; ok {
// 			if qty, ok := categorySupplies[offer.Name]; ok {
// 				quantity = qty
// 			}
// 		}
// 		offersWithPrice = append(offersWithPrice, domain.OfferWithPrice{
// 			ID:       offer.ID,
// 			Name:     offer.Name,
// 			Quantity: quantity / 5,
// 			Price:    offer.Price,
// 			Category: offer.Category,
// 		})
// 	}

// 	return offersWithPrice, nil
// }

func (s *OfferService) GetAllOffers() ([]domain.OfferWithPrice, error) {
	// Obtener las ofertas desde el repositorio
	offers, err := s.offerRepository.GetOffersData()
	if err != nil {
		return nil, err
	}

	// Hacer la solicitud HTTP al servidor cppServer
	var resp *http.Response
	if os.Getenv("RUN_LOCAL") == "true" {
		resp, err = http.Get("http://localhost:8083/supplies")
	} else {
		resp, err = http.Get("http://cppserver:8083/supplies")
	}
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Leer el cuerpo de la respuesta
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Deserializar en la estructura completa
	var serverData ServerResponse
	err = json.Unmarshal(body, &serverData)
	if err != nil {
		return nil, err
	}

	// Extraer solo supplies
	supplies := serverData.Supplies

	// Crear las ofertas con cantidad
	var offersWithPrice []domain.OfferWithPrice
	for _, offer := range offers {
		quantity := 0
		var normalizedCategory string
		switch offer.Category {
		case "food":
			normalizedCategory = "food"
		case "drink":
			normalizedCategory = "food"
		default:
			normalizedCategory = offer.Category
		}

		if categorySupplies, ok := supplies[normalizedCategory]; ok {
			if qty, ok := categorySupplies[offer.Name]; ok {
				quantity = qty
			}
		}
		offersWithPrice = append(offersWithPrice, domain.OfferWithPrice{
			ID:       offer.ID,
			Name:     offer.Name,
			Quantity: quantity / 5,
			Price:    offer.Price,
			Category: offer.Category,
		})
	}

	return offersWithPrice, nil
}

func (s *OfferService) ProcessOrder(email string, order domain.OrderCheckout) (error, domain.Order) {

	// Check each item's quantity in the order
	for _, item := range order.Order {
		if item.Quantity <= 0 {
			return fmt.Errorf("Bad Request"), domain.Order{}
		}
	}

	// Obtain the id of the user to be used as a foreign key
	user, err := s.UserService.GetUserByEmail(email)
	if err != nil {
		return err, domain.Order{}
	}

	// Count the total of the order
	total := 0
	for _, item := range order.Order {
		total += item.Quantity
	}

	// Insert the order in the database
	err = s.offerRepository.InsertOrder(user.ID, total)
	if err != nil {
		return err, domain.Order{}
	}

	// Create the response
	orderResponse := domain.Order{
		Total:  total,
		Status: "pending",
	}

	return nil, orderResponse
}

func (s *OfferService) GetOrderById(id string) (error, domain.UserOrderStatus) {
	// get the order from the database
	idOrder, user, status, total, err := s.offerRepository.GetOrderById(id)
	if err != nil {
		return err, domain.UserOrderStatus{}
	}

	// create the response
	orderResponse := domain.UserOrderStatus{
		ID:     idOrder,
		User:   user,
		Total:  total,
		Status: status,
	}

	return nil, orderResponse
}

func (s *OfferService) GetAllOrders() ([]domain.UserOrderStatus, error) {
	// get the orders from the database
	orders, err := s.offerRepository.GetAllOrders()
	if err != nil {
		return nil, err
	}

	// create the response
	var ordersResponse []domain.UserOrderStatus
	for _, order := range orders {
		orderResponse := domain.UserOrderStatus{
			ID:     order.ID,
			User:   order.User,
			Total:  order.Total,
			Status: order.Status,
		}
		ordersResponse = append(ordersResponse, orderResponse)
	}

	return ordersResponse, nil
}
