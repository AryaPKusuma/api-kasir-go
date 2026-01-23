package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

type Produk struct {
	ID     int    `json:"id"`
	Nama   string `json:"nama"`
	Harga  int    `json:"harga"`
	Stok   int    `json:"stok"`
	Detail string `json:"detail,omitempty"`
}

var produk = []Produk{
	{ID: 1, Nama: "Indomie Godog", Harga: 3500, Stok: 10, Detail: "Mie instan dengan rasa pedas"},
	{ID: 2, Nama: "Vit 1000ml", Harga: 3000, Stok: 40, Detail: "Minuman vitamin dengan rasa buah"},
	{ID: 3, Nama: "kecap", Harga: 12000, Stok: 20, Detail: "Kecap manis"},
}

func Products(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		GetAllProducts(w, r)
		return
	} else if r.Method == "POST" {
		AddProducts(w, r)
		return
	}
	http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
}

func ProductById(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		GetProductByID(w, r)
		return
	} else if r.Method == "PUT" {
		UpdateProducts(w, r)
		return
	} else if r.Method == "DELETE" {
		DeleteProducts(w, r)
		return
	}
	http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
}

func GetAllProducts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(produk)
}

func AddProducts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var newProduct Produk
	id := len(produk) + 1
	newProduct.ID = id

	err := json.NewDecoder(r.Body).Decode(&newProduct)
	if err != nil {
		http.Error(w, "invalid input", http.StatusBadRequest)
	}
	produk = append(produk, newProduct)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newProduct)
}

func GetProductByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	idStr := strings.TrimPrefix(r.URL.Path, "/api/products/")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		http.Error(w, "invalid product id", http.StatusBadRequest)
		return
	}

	for _, p := range produk {
		if p.ID == id {
			json.NewEncoder(w).Encode(p)
			return
		}
	}

	http.Error(w, "product not found", http.StatusNotFound)
}

func UpdateProducts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	idStr := strings.TrimPrefix(r.URL.Path, "/api/products/")
	id, err := strconv.Atoi(idStr)
	var updateProducts Produk

	if err != nil {
		http.Error(w, "invalid product id", http.StatusBadRequest)
		return
	}

	err = json.NewDecoder(r.Body).Decode(&updateProducts)
	if err != nil {
		http.Error(w, "invalid input", http.StatusBadRequest)
	}

	for i := range produk {
		if produk[i].ID == id {
			updateProducts.ID = id
			produk[i] = updateProducts
			json.NewEncoder(w).Encode(updateProducts)
			return
		}
	}
}

func DeleteProducts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	idStr := strings.TrimPrefix(r.URL.Path, "/api/products/")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		http.Error(w, "invalid product id", http.StatusBadRequest)
		return
	}

	for i := range produk {
		if produk[i].ID == id {
			produk = append(produk[:i], produk[i+1:]...)
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}

	http.Error(w, "product not found", http.StatusNotFound)
}
