package repositories

import (
	"database/sql"
	"log"

	"github.com/ICOMP-UNC/newworld-LaureanoOlocco/internal/core/domain"
)

type OfferRepository struct {
	db *sql.DB
}

func NewOfferRepository(db *sql.DB) *OfferRepository {
	return &OfferRepository{
		db: db,
	}
}

func (r *OfferRepository) GetOffersData() ([]domain.Offer, error) {

	// Query to get all offers from the database
	query := `SELECT id, name, price, category FROM offers`
	rows, err := r.db.Query(query)
	if err != nil {
		log.Printf("Error querying offers: %q", err)
		return nil, err
	}
	defer rows.Close()

	// Iterate over the rows and create the offers
	var offers []domain.Offer
	for rows.Next() {
		var offer domain.Offer
		err := rows.Scan(&offer.ID, &offer.Name, &offer.Price, &offer.Category)
		if err != nil {
			log.Printf("Error scanning offer: %q", err)
			return nil, err
		}
		offers = append(offers, offer)
	}

	// check for errors during the iteration
	if err := rows.Err(); err != nil {
		log.Printf("Error iterating over offers: %q", err)
		return nil, err
	}

	return offers, nil
}

func (r *OfferRepository) InsertOrder(userID, total int) error {
	query := `INSERT INTO orders (status, user_id, total) VALUES ($1, $2, $3) RETURNING id`
	var orderID int
	err := r.db.QueryRow(query, "pending", userID, total).Scan(&orderID)
	if err != nil {
		return err
	}
	return nil
}

func (r *OfferRepository) GetOrderById(id string) (int, string, string, int, error) {
	query := `
        SELECT o.id, o.status, o.total, u.username
        FROM orders o
        INNER JOIN users u ON o.user_id = u.id
        WHERE o.id = $1
    `
	row := r.db.QueryRow(query, id)

	var orderID int
	var status string
	var total int
	var username string

	err := row.Scan(&orderID, &status, &total, &username)
	if err != nil {
		return 0, "", "", 0, err
	}

	return orderID, username, status, total, nil
}

func (r *OfferRepository) GetAllOrders() ([]domain.UserOrderStatus, error) {
	query := `
		SELECT o.id, o.status, o.total, u.username
		FROM orders o
		INNER JOIN users u ON o.user_id = u.id
	`
	rows, err := r.db.Query(query)
	if err != nil {
		log.Printf("Error querying orders: %q", err)
		return nil, err
	}
	defer rows.Close()

	var orders []domain.UserOrderStatus
	for rows.Next() {
		var order domain.UserOrderStatus
		err := rows.Scan(&order.ID, &order.Status, &order.Total, &order.User)
		if err != nil {
			log.Printf("Error scanning order: %q", err)
			return nil, err
		}
		orders = append(orders, order)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Error iterating over orders: %q", err)
		return nil, err
	}

	return orders, nil
}
