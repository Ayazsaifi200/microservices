package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	
	_ "github.com/lib/pq"
	pb "github.com/GoogleCloudPlatform/microservices-demo/src/checkoutservice/genproto"
)

var db *sql.DB

// InitDatabase initializes PostgreSQL connection
func InitDatabase() error {
	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		dbHost = "postgres"
	}
	
	dbPort := os.Getenv("DB_PORT")
	if dbPort == "" {
		dbPort = "5432"
	}
	
	dbUser := os.Getenv("DB_USER")
	if dbUser == "" {
		dbUser = "postgres"
	}
	
	dbPassword := os.Getenv("DB_PASSWORD")
	if dbPassword == "" {
		dbPassword = "postgres123"
	}
	
	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = "ordersdb"
	}
	
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)
	
	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}
	
	// Test connection
	err = db.Ping()
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	
	log.Info("Successfully connected to PostgreSQL database")
	return nil
}

// SaveOrder saves order details to PostgreSQL
func SaveOrder(ctx context.Context, orderResult *pb.OrderResult, req *pb.PlaceOrderRequest) error {
	if db == nil {
		log.Warn("Database not initialized, skipping order save")
		return nil
	}
	
	// Start transaction
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()
	
	// Insert order
	orderQuery := `
		INSERT INTO orders (
			order_id, user_id, user_email, user_currency,
			shipping_address_street, shipping_address_city, shipping_address_state,
			shipping_address_country, shipping_address_zip_code,
			shipping_tracking_id, shipping_cost_currency, shipping_cost_units, shipping_cost_nanos,
			total_items
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
	`
	
	_, err = tx.ExecContext(ctx, orderQuery,
		orderResult.OrderId,
		req.UserId,
		req.Email,
		req.UserCurrency,
		req.Address.StreetAddress,
		req.Address.City,
		req.Address.State,
		req.Address.Country,
		fmt.Sprintf("%d", req.Address.ZipCode),
		orderResult.ShippingTrackingId,
		orderResult.ShippingCost.CurrencyCode,
		orderResult.ShippingCost.Units,
		orderResult.ShippingCost.Nanos,
		len(orderResult.Items),
	)
	
	if err != nil {
		return fmt.Errorf("failed to insert order: %w", err)
	}
	
	// Insert order items
	itemQuery := `
		INSERT INTO order_items (
			order_id, product_id, quantity, cost_currency, cost_units, cost_nanos
		) VALUES ($1, $2, $3, $4, $5, $6)
	`
	
	for _, item := range orderResult.Items {
		_, err = tx.ExecContext(ctx, itemQuery,
			orderResult.OrderId,
			item.Item.ProductId,
			item.Item.Quantity,
			item.Cost.CurrencyCode,
			item.Cost.Units,
			item.Cost.Nanos,
		)
		if err != nil {
			return fmt.Errorf("failed to insert order item: %w", err)
		}
	}
	
	// Commit transaction
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	
	log.Infof("Successfully saved order %s to database", orderResult.OrderId)
	return nil
}
