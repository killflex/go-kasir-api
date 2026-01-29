package main

import (
	"encoding/json"
	"fmt"
	"kasir-api/database"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Port   string `mapstructure:"PORT"`
	DBconn string `mapstructure:"DB_CONN"`
}

type Product struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Price string `json:"price"`
	Stock int    `json:"stock"`
}

type Category struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

var produk = []Product{
	{ID: 1, Name: "Indomie Goreng", Price: "Rp 3.500", Stock: 100},
	{ID: 2, Name: "Teh Botol Sosro", Price: "Rp 5.000", Stock: 50},
	{ID: 3, Name: "Aqua 600ml", Price: "Rp 4.000", Stock: 200},
}

var category = []Category{
	{ID: 1, Name: "Makanan", Description: "Product makanan ringan dan berat"},
	{ID: 2, Name: "Minuman", Description: "Berbagai jenis minuman kemasan"},
	{ID: 3, Name: "Sembako", Description: "Kebutuhan pokok sehari-hari"},
	{ID: 4, Name: "Snack", Description: "Camilan dan makanan ringan"},
	{ID: 5, Name: "Perlengkapan Rumah", Description: "Alat dan perlengkapan rumah tangga"},
}

func httpError(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(map[string]string{"error": message}); err != nil {
		fmt.Println("Failed to encode error message:", err)
	}
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]string{"message": "API is running"}); err != nil {
		httpError(w, http.StatusInternalServerError, "Failed to encode health check response")
	}
}

func produkHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case "GET":
		if err := json.NewEncoder(w).Encode(produk); err != nil {
			httpError(w, http.StatusInternalServerError, "Failed to encode produk data")
		}
	case "POST":
		var newProduk Product
		if err := json.NewDecoder(r.Body).Decode(&newProduk); err != nil {
			httpError(w, http.StatusBadRequest, "Invalid request payload")
			return
		}

		maxID := 0
		for _, p := range produk {
			if p.ID > maxID {
				maxID = p.ID
			}
		}
		newProduk.ID = maxID + 1

		produk = append(produk, newProduk)
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(map[string]interface{}{
			"data":    newProduk,
			"message": "Product created successfully",
		}); err != nil {
			httpError(w, http.StatusInternalServerError, "Failed to encode new produk data")
		}
	default:
		httpError(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

func produkByIDHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := r.URL.Path[len("/api/produk/"):]
	if id == "" {
		httpError(w, http.StatusBadRequest, "Product ID is required")
		return
	}

	switch r.Method {
	case "GET":
		for _, p := range produk {
			if fmt.Sprintf("%d", p.ID) == id {
				if err := json.NewEncoder(w).Encode(p); err != nil {
					httpError(w, http.StatusInternalServerError, "Failed to encode produk data")
				}
				return
			}
		}
		httpError(w, http.StatusNotFound, "Product not found")
	case "PUT":
		var updatedProduk Product
		if err := json.NewDecoder(r.Body).Decode(&updatedProduk); err != nil {
			httpError(w, http.StatusBadRequest, "Invalid request payload")
			return
		}

		for i, p := range produk {
			if fmt.Sprintf("%d", p.ID) == id {
				updatedProduk.ID = p.ID
				produk[i] = updatedProduk
				if err := json.NewEncoder(w).Encode(map[string]interface{}{
					"data":    updatedProduk,
					"message": "Product updated successfully",
				}); err != nil {
					httpError(w, http.StatusInternalServerError, "Failed to encode updated produk data")
				}
				return
			}
		}
		httpError(w, http.StatusNotFound, "Product not found")
	case "DELETE":
		for i, p := range produk {
			if fmt.Sprintf("%d", p.ID) == id {
				produk = append(produk[:i], produk[i+1:]...)
				if err := json.NewEncoder(w).Encode(map[string]interface{}{
					"id":      p.ID,
					"message": "Product deleted successfully",
				}); err != nil {
					httpError(w, http.StatusInternalServerError, "Failed to encode delete confirmation message")
				}
				return
			}
		}
		httpError(w, http.StatusNotFound, "Product not found")
	default:
		httpError(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

func CategoryHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case "GET":
		if err := json.NewEncoder(w).Encode(category); err != nil {
			httpError(w, http.StatusInternalServerError, "Failed to encode category data")
		}
	case "POST":
		var newCategory Category
		if err := json.NewDecoder(r.Body).Decode(&newCategory); err != nil {
			httpError(w, http.StatusBadRequest, "Invalid request payload")
		}

		maxID := 0
		for _, p := range category {
			if p.ID > maxID {
				maxID = p.ID
			}
		}
		newCategory.ID = maxID + 1

		category = append(category, newCategory)
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(map[string]interface{}{
			"data":    newCategory,
			"message": "Category created successfully",
		}); err != nil {
			httpError(w, http.StatusInternalServerError, "Failed to encode new category data")
		}
	default:
		httpError(w, http.StatusInternalServerError, "Method not allowed")
	}
}

func CategoryByIDHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := r.URL.Path[len("/category/"):]
	if id == "" {
		httpError(w, http.StatusBadRequest, "Category ID is required")
		return
	}

	switch r.Method {
	case "GET":
		for _, p := range category {
			if fmt.Sprintf("%d", p.ID) == id {
				if err := json.NewEncoder(w).Encode(p); err != nil {
					httpError(w, http.StatusInternalServerError, "Failed to encode category data by ID")
				}
				return
			}
		}
		httpError(w, http.StatusNotFound, "Category ID is not found")
	case "PUT":
		var updatedCategory Category
		if err := json.NewDecoder(r.Body).Decode(&updatedCategory); err != nil {
			httpError(w, http.StatusBadRequest, "Invalid request payload")
			return
		}

		for i, p := range category {
			if fmt.Sprintf("%d", p.ID) == id {
				updatedCategory.ID = p.ID
				category[i] = updatedCategory
				if err := json.NewEncoder(w).Encode(map[string]interface{}{
					"data":    updatedCategory,
					"message": "Category updated successfully",
				}); err != nil {
					httpError(w, http.StatusInternalServerError, "Failed to encode updated category data")
				}
				return
			}
		}
		httpError(w, http.StatusNotFound, "Category by ID not found")
	case "DELETE":
		for i, p := range category {
			if fmt.Sprintf("%d", p.ID) == id {
				category = append(category[:i], category[i+1:]...)
				if err := json.NewEncoder(w).Encode(map[string]interface{}{
					"id":      p.ID,
					"message": "Category deleted successfully",
				}); err != nil {
					httpError(w, http.StatusInternalServerError, "Failed to encode delete confirmation message")
				}
				return
			}
		}
		httpError(w, http.StatusNotFound, "Category by ID not found")
	default:
		httpError(w, http.StatusMethodNotAllowed, "Method is not allowed")
	}
}

func main() {
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	config := Config{
		Port:   viper.GetString("PORT"),
		DBconn: viper.GetString("DB_CONN"),
	}

	if _, err := os.Stat(".env"); err == nil {
		viper.SetConfigFile(".env")
		_ = viper.ReadInConfig()
	}

	db, err := database.Connect(config.DBconn)
	if err != nil {
		log.Fatal("Error connecting to the database: %v\n", err)
		return
	}
	defer db.Close()

	// GET localhost:8080/health
	http.HandleFunc("/health", healthHandler)

	// GET localhost:8080/api/produk (get all)
	// POST localhost:8080/api/produk (create new)
	http.HandleFunc("/api/produk", produkHandler)

	// GET localhost:8080/api/produk/{id} (get by ID)
	// PUT localhost:8080/api/produk/{id} (update by ID)
	// DELETE localhost:8080/api/produk/{id} (delete by ID)
	http.HandleFunc("/api/produk/", produkByIDHandler)

	// GET localhost:8080/category (get all category)
	// POST localhost:8080/category (create new category)
	http.HandleFunc("/category", CategoryHandler)

	// GET localhost:8080/category/{id} (get by ID)
	// PUT localhost:8080/category/{id} (update by ID)
	// DELETE localhost:8080/category/{id} (delete by ID)
	http.HandleFunc("/category/", CategoryByIDHandler)

	fmt.Printf("Server running on http://localhost:%s\n", config.Port)
	if err := http.ListenAndServe(":"+config.Port, nil); err != nil {
		fmt.Printf("Error starting server %v\n", err)
	}
}
