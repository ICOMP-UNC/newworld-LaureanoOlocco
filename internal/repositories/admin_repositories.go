package repositories

import (
	"database/sql"
	"log"

	"github.com/ICOMP-UNC/newworld-LaureanoOlocco/internal/core/domain"
)

type AdminRepository struct {
	db *sql.DB
}

func NewAdminRepository(db *sql.DB) *AdminRepository {
	return &AdminRepository{
		db: db,
	}
}

func (r *AdminRepository) UpdateOrderStatus(orderID, newStatus string) (domain.UserOrderStatus, error) {
	// Begin a transaction
	tx, err := r.db.Begin()
	if err != nil {
		return domain.UserOrderStatus{}, err
	}
	defer func() {
		if rollbackErr := tx.Rollback(); rollbackErr != nil && rollbackErr != sql.ErrTxDone {
			log.Printf("tx.Rollback failed: %v", rollbackErr)
		}
	}()

	// Update the order status
	updateQuery := `UPDATE orders SET status = $1 WHERE id = $2`
	_, err = tx.Exec(updateQuery, newStatus, orderID)
	if err != nil {
		return domain.UserOrderStatus{}, err
	}

	// Retrieve the updated order and user information
	var orderStatus domain.UserOrderStatus
	query := `
		SELECT 
			o.id, u.username, o.total, o.status 
		FROM 
			orders o
		INNER JOIN 
			users u ON o.user_id = u.id
		WHERE 
			o.id = $1
	`
	row := tx.QueryRow(query, orderID)
	err = row.Scan(&orderStatus.ID, &orderStatus.User, &orderStatus.Total, &orderStatus.Status)
	if err != nil {
		return domain.UserOrderStatus{}, err
	}

	// Commit the transaction for asure the ACID
	if err := tx.Commit(); err != nil {
		return domain.UserOrderStatus{}, err
	}

	return orderStatus, nil
}
