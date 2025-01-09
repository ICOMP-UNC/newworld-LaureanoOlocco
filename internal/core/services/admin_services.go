package services

import (
	"github.com/ICOMP-UNC/newworld-LaureanoOlocco/internal/core/domain"
	"github.com/ICOMP-UNC/newworld-LaureanoOlocco/internal/core/ports"
)

type AdminService struct {
	adminRepository ports.IAdminRepository
	offerService    *OfferService
}

func NewAdminService(repository ports.IAdminRepository, offerService *OfferService) *AdminService {
	return &AdminService{
		adminRepository: repository,
		offerService:    offerService,
	}
}

func (s *AdminService) GetDashboardData() (domain.Dashboard, error) {

	// get all the offers from the offer service
	offers, err := s.offerService.GetAllOffers()
	if err != nil {
		return domain.Dashboard{}, err
	}

	// get all the orders from the offer service
	orders, err := s.offerService.GetAllOrders()
	if err != nil {
		return domain.Dashboard{}, err
	}

	// get the balance from all the orders
	balance := 0
	for _, order := range orders {
		balance += order.Total
	}

	// create the dashboard data
	dashboard := domain.Dashboard{
		Offers:  offers,
		Orders:  orders,
		Balance: balance,
	}

	return dashboard, nil
}

func (s *AdminService) UpdateOrderStatus(orderID, Newstatus string) (domain.UserOrderStatus, error) {

	// update the order status
	order, err := s.adminRepository.UpdateOrderStatus(orderID, Newstatus)
	if err != nil {
		return domain.UserOrderStatus{}, err
	}

	return order, nil
}
