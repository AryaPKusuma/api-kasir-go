package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"kasir/models"
	"kasir/services"
)

type TransactionHandler struct {
	service *services.TransactionService
}

func NewTransactionHandler(service *services.TransactionService) *TransactionHandler {
	return &TransactionHandler{service: service}
}

// multiple item apa aja, quantity nya
func (h *TransactionHandler) HandleCheckout(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.Checkout(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *TransactionHandler) Checkout(w http.ResponseWriter, r *http.Request) {
	var req models.CheckoutRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate request items
	for _, item := range req.Items {
		if item.ProductID <= 0 {
			http.Error(w, fmt.Sprintf("invalid product id: %d", item.ProductID), http.StatusBadRequest)
			return
		}
		if item.Quantity <= 0 {
			http.Error(w, fmt.Sprintf("invalid quantity for product %d: %d", item.ProductID, item.Quantity), http.StatusBadRequest)
			return
		}
	}

	transaction, err := h.service.Checkout(req.Items, true) // Enable row-level locking for concurrent transactions
	if err != nil {
		// Check if it's a business logic error (like insufficient stock) vs internal server error
		if strings.Contains(err.Error(), "insufficient stock") || strings.Contains(err.Error(), "not found") {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(transaction)
}

func (h *TransactionHandler) Summary(w http.ResponseWriter, r *http.Request) {
	startDate := r.URL.Query().Get("start_date")
	endDate := r.URL.Query().Get("end_date")

	summary, err := h.service.GetTransactionSummary(startDate, endDate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(summary)
}
