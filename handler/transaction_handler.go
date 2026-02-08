package handler

import (
	"encoding/json"
	"kasir-api/model"
	"kasir-api/service"
	"net/http"
)

type TransactionHandler struct {
	service *service.TransactionService
}

func NewTransactionHandler(service *service.TransactionService) *TransactionHandler {
	return &TransactionHandler{service: service}
}

func (h *TransactionHandler) CheckoutItem(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req model.CheckoutRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		RespondError(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	transaction, err := h.service.CheckoutItem(req.Items)
	if err != nil {
		RespondError(w, http.StatusInternalServerError, "Failed to checkout", err)
		return
	}

	RespondSuccess(w, http.StatusCreated, "Checkout successfully", transaction)
}

func (h *TransactionHandler) HandleCheckoutItem(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.CheckoutItem(w, r)
	default:
		RespondError(w, http.StatusMethodNotAllowed, "Method not allowed", nil)
	}
}