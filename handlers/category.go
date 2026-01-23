package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

type Category struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

var categories = []Category{
	{ID: 1, Name: "Makanan", Description: "Kategori untuk semua jenis makanan"},
	{ID: 2, Name: "Minuman", Description: "Kategori untuk semua jenis minuman"},
	{ID: 3, Name: "Bumbu", Description: "Kategori untuk bumbu dapur"},
}

func Categories(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		GetAllCategories(w, r)
		return
	} else if r.Method == "POST" {
		AddCategory(w, r)
		return
	}
	http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
}

func CategoriesById(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		GetCategoryByID(w, r)
		return
	} else if r.Method == "PUT" {
		UpdateCategory(w, r)
		return
	} else if r.Method == "DELETE" {
		DeleteCategory(w, r)
		return
	}
	http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
}

// Get all categories
func GetAllCategories(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(categories)
}

// Add a new category
func AddCategory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var newCategory Category
	id := len(categories) + 1
	newCategory.ID = id

	err := json.NewDecoder(r.Body).Decode(&newCategory)
	if err != nil {
		http.Error(w, "invalid input", http.StatusBadRequest)
		return
	}
	categories = append(categories, newCategory)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newCategory)
}

// Get category by ID
func GetCategoryByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	idStr := strings.TrimPrefix(r.URL.Path, "/api/categories/")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		http.Error(w, "invalid category id", http.StatusBadRequest)
		return
	}

	for _, c := range categories {
		if c.ID == id {
			json.NewEncoder(w).Encode(c)
			return
		}
	}

	http.Error(w, "category not found", http.StatusNotFound)
}

// Update category by ID
func UpdateCategory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	idStr := strings.TrimPrefix(r.URL.Path, "/api/categories/")
	id, err := strconv.Atoi(idStr)
	var updateCategory Category

	if err != nil {
		http.Error(w, "invalid category id", http.StatusBadRequest)
		return
	}

	err = json.NewDecoder(r.Body).Decode(&updateCategory)
	if err != nil {
		http.Error(w, "invalid input", http.StatusBadRequest)
		return
	}

	for i := range categories {
		if categories[i].ID == id {
			updateCategory.ID = id
			categories[i] = updateCategory
			json.NewEncoder(w).Encode(updateCategory)
			return
		}
	}

	http.Error(w, "category not found", http.StatusNotFound)
}

// Delete category by ID
func DeleteCategory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	idStr := strings.TrimPrefix(r.URL.Path, "/api/categories/")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		http.Error(w, "invalid category id", http.StatusBadRequest)
		return
	}

	for i := range categories {
		if categories[i].ID == id {
			categories = append(categories[:i], categories[i+1:]...)
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}

	http.Error(w, "category not found", http.StatusNotFound)
}
