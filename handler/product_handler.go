package handler

import (
	"encoding/json"
	"kasir-api/model"
	"kasir-api/service"
	"net/http"
	"strconv"
	"strings"
)

type ProductHandler struct {
	service *service.ProductService
}

func NewProductHandler(service *service.ProductService) *ProductHandler {
	return &ProductHandler{service: service}
}

func (h *ProductHandler) GetAllProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	products, err := h.service.GetAllProduct()
	if err != nil {
		RespondError(w, http.StatusInternalServerError, "Failed to retrieve products", err.Error())
		return
	}

	RespondSuccess(w, http.StatusOK, "Product found",products)
}

func (h *ProductHandler) GetProductById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	idStr := strings.TrimPrefix(r.URL.Path, "/api/product/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		RespondError(w, http.StatusBadRequest, err.Error())
		http.Error(w, "Invalid Product ID", http.StatusBadRequest)
		return
	}

	product, err := h.service.GetProductById(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := json.NewEncoder(w).Encode(product); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var product model.Product
	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if product.Name == "" {
		http.Error(w, "Product name is required", http.StatusBadRequest)
		return
	}
	if product.Price <= 0 {
		http.Error(w, "Product price must be greater than 0", http.StatusBadRequest)
		return
	}
	if product.Stock < 0 {
		http.Error(w, "Product stock cannot be negative", http.StatusBadRequest)
		return
	}

	err = h.service.CreateProduct(&product)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Product created successfully",
		"data":    product,
	}); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (h *ProductHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	idStr := strings.TrimPrefix(r.URL.Path, "/api/product/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid Product ID", http.StatusBadRequest)
		return
	}

	var product model.Product
	err = json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	product.ID = id
	err = h.service.UpdateProduct(&product)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := json.NewEncoder(w).Encode(product); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (h *ProductHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	idStr := strings.TrimPrefix(r.URL.Path, "/api/product/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid Product ID", http.StatusBadRequest)
		return
	}

	err = h.service.DeleteProduct(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(map[string]string{
		"message": "Product deleted successfully",
	}); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (h *ProductHandler) HandleProducts(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.GetAllProduct(w, r)
	case http.MethodPost:
		h.CreateProduct(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *ProductHandler) HandleProductByID(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.GetProductById(w, r)
	case http.MethodPut:
		h.UpdateProduct(w, r)
	case http.MethodDelete:
		h.DeleteProduct(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
