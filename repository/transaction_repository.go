package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"kasir-api/model"
)

type TransactionRepository struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

func (repo *TransactionRepository) CheckoutItem(items []model.CheckoutItem) (*model.Transaction, error) {
    if len(items) == 0 {
        return nil, errors.New("checkout items cannot be empty")
    }
    
    tx, err := repo.db.Begin()
    if err != nil {
        return nil, err
    }
    defer tx.Rollback()

    totalAmount := 0
    details := make([]model.TransactionDetail, 0, len(items))

    for _, item := range items {
        var productPrice, stock int
        var productName string

        if item.Quantity <= 0 {
            return nil, errors.New("quantity must be greater than 0")
        }

        err := tx.QueryRow("SELECT name, price, stock FROM products WHERE id = $1", item.ProductID).Scan(&productName, &productPrice, &stock)
        if err == sql.ErrNoRows {
            return nil, fmt.Errorf("product with ID %d not found", item.ProductID)
        }
        if err != nil {
            return nil, fmt.Errorf("failed to get product: %w", err)
        }

        if stock < item.Quantity {
            return nil, fmt.Errorf("insufficient stock for product '%s': available %d, requested %d", 
                productName, stock, item.Quantity)
        }

        subtotal := productPrice * item.Quantity
        totalAmount += subtotal
        
        result, err := tx.Exec("UPDATE products SET stock = stock - $1 WHERE id = $2 AND stock >= $1", item.Quantity, item.ProductID)
        if err != nil {
            return nil, fmt.Errorf("failed to update stock: %w", err)
        }

        rowsAffected, err := result.RowsAffected()
        if err != nil {
            return nil, fmt.Errorf("failed to check affected rows: %w", err)
        }
        if rowsAffected == 0 {
            return nil, fmt.Errorf("failed to update stock for product '%s' (concurrent update?)", productName)
        }

        details = append(details, model.TransactionDetail{
            ProductID:   item.ProductID,
            ProductName: productName,
            Quantity:    item.Quantity,
            Subtotal:    subtotal,
        })
    }

    // Insert transaction
    var transactionID int
    err = tx.QueryRow("INSERT INTO transactions (total_amount) VALUES ($1) RETURNING id", totalAmount).Scan(&transactionID)
    if err != nil {
        return nil, fmt.Errorf("failed to create transaction: %w", err)
    }

    // ✅ ADD THIS: Insert transaction details
    for i := range details {
        details[i].TransactionID = transactionID
        _, err = tx.Exec(
            "INSERT INTO transaction_details (transaction_id, product_id, quantity, subtotal) VALUES ($1, $2, $3, $4)",
            transactionID, details[i].ProductID, details[i].Quantity, details[i].Subtotal,
        )
        if err != nil {
            return nil, fmt.Errorf("failed to insert transaction detail: %w", err)
        }
    }

    // ✅ Commit the transaction
    if err := tx.Commit(); err != nil {
        return nil, fmt.Errorf("failed to commit transaction: %w", err)
    }

    return &model.Transaction{
        ID:          transactionID,
        TotalAmount: totalAmount,
        Details:     details,
    }, nil
}