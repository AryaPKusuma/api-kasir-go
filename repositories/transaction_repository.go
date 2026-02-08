package repositories

import (
	"database/sql"
	"fmt"
	"kasir/models"
	"strings"
)

type TransactionRepository struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

func (repo *TransactionRepository) CreateTransaction(items []models.CheckoutItem, useLock bool) (*models.Transaction, error) {
	tx, err := repo.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	totalAmount := 0
	details := make([]models.TransactionDetail, 0)

	for _, item := range items {
		if item.Quantity <= 0 {
			return nil, fmt.Errorf("quantity must be greater than 0 for product id %d", item.ProductID)
		}

		var productPrice, stock int
		var productName string

		query := "SELECT name, price, stock FROM products WHERE id = $1"
		if useLock {
			query += " FOR UPDATE"
		}

		err := tx.QueryRow(query, item.ProductID).Scan(&productName, &productPrice, &stock)
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("product id %d not found", item.ProductID)
		}
		if err != nil {
			return nil, err
		}

		if stock < item.Quantity {
			return nil, fmt.Errorf("insufficient stock for product id %d: requested %d, available %d", item.ProductID, item.Quantity, stock)
		}

		subtotal := productPrice * item.Quantity
		totalAmount += subtotal

		_, err = tx.Exec("UPDATE products SET stock = stock - $1 WHERE id = $2", item.Quantity, item.ProductID)
		if err != nil {
			return nil, err
		}

		details = append(details, models.TransactionDetail{
			ProductID:   item.ProductID,
			ProductName: productName,
			Quantity:    item.Quantity,
			Subtotal:    subtotal,
		})
	}

	var transactionID int
	err = tx.QueryRow("INSERT INTO transactions (total_amount) VALUES ($1) RETURNING id", totalAmount).Scan(&transactionID)
	if err != nil {
		return nil, err
	}

	if len(details) > 0 {
		// Prepare a single query for batch insert
		valueStrings := make([]string, 0, len(details))
		valueArgs := make([]interface{}, 0, len(details)*4)

		for i, detail := range details {
			idx := i * 4
			valueStrings = append(valueStrings, fmt.Sprintf("($%d, $%d, $%d, $%d)", idx+1, idx+2, idx+3, idx+4))
			valueArgs = append(valueArgs, transactionID, detail.ProductID, detail.Quantity, detail.Subtotal)
		}

		query := fmt.Sprintf("INSERT INTO transaction_details (transaction_id, product_id, quantity, subtotal) VALUES %s",
			strings.Join(valueStrings, ", "))

		_, err = tx.Exec(query, valueArgs...)
		if err != nil {
			return nil, err
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &models.Transaction{
		ID:          transactionID,
		TotalAmount: totalAmount,
		Details:     details,
	}, nil
}

func (repo *TransactionRepository) GetTransactionSummary(startDate, endDate string) (map[string]interface{}, error) {
	summary := make(map[string]interface{})

	// Build query conditions based on date parameters
	whereClause := "WHERE t.total_amount IS NOT NULL"
	params := []interface{}{}
	paramIndex := 1

	if startDate != "" {
		whereClause += fmt.Sprintf(" AND DATE(t.created_at) >= $%d", paramIndex)
		params = append(params, startDate)
		paramIndex++
	}
	if endDate != "" {
		whereClause += fmt.Sprintf(" AND DATE(t.created_at) <= $%d", paramIndex)
		params = append(params, endDate)
		paramIndex++
	}

	query := fmt.Sprintf(`
		SELECT
			COALESCE(SUM(t.total_amount), 0) as total_revenue,
			COUNT(DISTINCT t.id) as total_transaction
		FROM transactions t %s`, whereClause)

	stmt, err := repo.db.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var totalRevenue int64
	var totalTransaction int

	err = stmt.QueryRow(params...).Scan(&totalRevenue, &totalTransaction)
	if err != nil {
		return nil, err
	}

	summary["total_revenue"] = totalRevenue
	summary["total_transaction"] = totalTransaction

	// Get best selling product
	bestProductQuery := fmt.Sprintf(`
		SELECT p.name, SUM(td.quantity) as total_quantity
		FROM transaction_details td
		JOIN transactions t ON td.transaction_id = t.id
		JOIN products p ON td.product_id = p.id
		%s
		GROUP BY p.id, p.name
		ORDER BY total_quantity DESC
		LIMIT 1`, whereClause)

	bestStmt, err := repo.db.Prepare(bestProductQuery)
	if err != nil {
		return nil, err
	}
	defer bestStmt.Close()

	var bestProductName string
	var bestProductQuantity int
	err = bestStmt.QueryRow(params...).Scan(&bestProductName, &bestProductQuantity)
	if err != nil {
		if err == sql.ErrNoRows {
			// No transactions found, return empty best product
			summary["best_products"] = map[string]interface{}{"name": "", "quantity": 0}
		} else {
			return nil, err
		}
	} else {
		summary["best_products"] = map[string]interface{}{"name": bestProductName, "quantity": bestProductQuantity}
	}

	return summary, nil
}
