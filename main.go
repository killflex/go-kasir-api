package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Produk struct {
	ID int `json:"id"`
	Name string `json:"name"`
	Harga string `json:"harga"`
	Stok int `json:"stok"`
}

var produk = []Produk{
	{ID: 1, Name: "Indomie Goreng", Harga: "Rp 3.500", Stok: 100},
	{ID: 2, Name: "Teh Botol Sosro", Harga: "Rp 5.000", Stok: 50},
	{ID: 3, Name: "Aqua 600ml", Harga: "Rp 4.000", Stok: 200},
}

func httpError(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	err := json.NewEncoder(w).Encode(map[string]string{"error": message})
	if err != nil {
		fmt.Println("Failed to encode error message:", err)
	}
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(map[string]string{"status": "OK", "message": "API is running"})
	if err != nil {
		httpError(w, http.StatusInternalServerError, "Failed to encode health check response")
	}
}

func produkHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
		case "GET":
			w.WriteHeader(http.StatusOK)
			err := json.NewEncoder(w).Encode(produk)
			if err != nil {
				httpError(w, http.StatusInternalServerError, "Failed to encode produk data")
			}
		case "POST":
			var newProduk Produk
			err := json.NewDecoder(r.Body).Decode(&newProduk)
			if err != nil {
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
			err = json.NewEncoder(w).Encode(map[string]interface{}{
				"message": "Produk created successfully",
				"data": newProduk,
			})
			if err != nil {
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
					w.WriteHeader(http.StatusOK)
					err := json.NewEncoder(w).Encode(p)
					if err != nil {
						httpError(w, http.StatusInternalServerError, "Failed to encode produk data")
					}
					return
				}
			}
			httpError(w, http.StatusNotFound, "Produk not found")
		case "PUT":
			var updatedProduk Produk
			err := json.NewDecoder(r.Body).Decode(&updatedProduk)
			if err != nil {
				httpError(w, http.StatusBadRequest, "Invalid request payload")
				return
			}

			for i, p := range produk {
				if fmt.Sprintf("%d", p.ID) == id {
					updatedProduk.ID = p.ID
					produk[i] = updatedProduk
					w.WriteHeader(http.StatusOK)
					err = json.NewEncoder(w).Encode(map[string]interface{}{
						"message": "Produk updated successfully",
						"data": updatedProduk,
					})
					if err != nil {
						httpError(w, http.StatusInternalServerError, "Failed to encode updated produk data")
					}
					return
				}
			}
			httpError(w, http.StatusNotFound, "Produk not found")
		case "DELETE":
			for i, p := range produk {
				if fmt.Sprintf("%d", p.ID) == id {
					produk = append(produk[:i], produk[i+1:]...)
					w.WriteHeader(http.StatusOK)
					err := json.NewEncoder(w).Encode(map[string]interface{}{
						"message": "Produk deleted successfully",
						"id": p.ID,
					})
					if err != nil {
						httpError(w, http.StatusInternalServerError, "Failed to encode delete confirmation message")
					}
					return
				}
			}
			httpError(w, http.StatusNotFound, "Produk not found")
		default:
			httpError(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

func main() {
	// localhost:8080/health
	http.HandleFunc("/health", healthHandler)

	// GET localhost:8080/api/produk (get all)
	// POST localhost:8080/api/produk (create new)
	http.HandleFunc("/api/produk", produkHandler)

	// GET localhost:8080/api/produk/{id} (get by ID)
	// PUT localhost:8080/api/produk/{id} (update by ID)
	// DELETE localhost:8080/api/produk/{id} (delete by ID)
	http.HandleFunc("/api/produk/", produkByIDHandler)

	fmt.Println("Starting server on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}