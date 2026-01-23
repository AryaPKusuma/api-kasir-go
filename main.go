package main

import (
	"fmt"
	"kasir/handlers"
	"net/http"
)

func Routes() {
	// health
	http.HandleFunc("/health", handlers.GetHealthStatus)

	// Product
	// get all dan post
	http.HandleFunc("/api/products", handlers.Products)
	// endpoint untuk delete, update, get by id
	http.HandleFunc("/api/products/{id}", handlers.ProductById)

	// Category routes
	http.HandleFunc("/api/categories", handlers.Categories)
	http.HandleFunc("/api/categories/{id}", handlers.CategoriesById)
}

func main() {
	Routes()
	port := ":9000"
	fmt.Println("Server running on port 8000")
	http.ListenAndServe(port, nil)
}
