package main

import (
	"fmt"
	"kasir/database"
	"kasir/handlers"
	"kasir/repositories"
	"kasir/services"
	"net/http"
	"os"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Port   string `mapstructure:"PORT"`
	DBConn string `mapstructure:"DB_CONN"`
}

func main() {

	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if _, err := os.Stat(".env"); err == nil {
		viper.SetConfigFile(".env")
		_ = viper.ReadInConfig()
	}

	config := Config{
		Port:   viper.GetString("PORT"),
		DBConn: viper.GetString("DB_CONN"),
	}

	fmt.Printf("Attempting to connect to database with connection string: %s\n", config.DBConn)

	db, err := database.InitDB(config.DBConn)
	if err != nil {
		fmt.Printf("Failed to connect to database: %v\n", err)
		return
	}
	defer db.Close()

	// Product setup
	productRepository := repositories.NewProductRepository(db)
	productService := services.NewProductService(productRepository)
	productHandler := handlers.NewProductHandler(productService)

	// Category setup
	categoryRepository := repositories.NewCategoryRepository(db)
	categoryService := services.NewCategoryService(categoryRepository)
	categoryHandler := handlers.NewCategoryHandler(categoryService)

	// Transaction setup
	transactionRepository := repositories.NewTransactionRepository(db)
	transactionService := services.NewTransactionService(transactionRepository)
	transactionHandler := handlers.NewTransactionHandler(transactionService)

	// Register routes
	http.HandleFunc("/health", handlers.GetHealthStatus)

	// Transaction routes
	http.HandleFunc("/api/transactions/checkout", transactionHandler.HandleCheckout)

	// Transaction report
	http.HandleFunc("/api/report", transactionHandler.Summary)

	// Category routes
	http.HandleFunc("/api/categories", categoryHandler.HandleCategories)
	http.HandleFunc("/api/categories/", categoryHandler.HandleCategoryByID)

	// Product routes
	http.HandleFunc("/api/produk", productHandler.HandleProducts)
	http.HandleFunc("/api/produk/", productHandler.HandleProductByID)

	addr := ":" + config.Port

	fmt.Printf("Server running on port %s\n", config.Port)
	http.ListenAndServe(addr, nil)
}
