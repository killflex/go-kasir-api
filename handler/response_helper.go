package handler

import (
	"encoding/json"
	"kasir-api/model"
	"log"
	"net/http"
)

func RespondJSON(w http.ResponseWriter, status int, payload model.APIResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		log.Printf("Failed to encode JSON response: %v\n", err)
	}
}

func RespondSuccess(w http.ResponseWriter, status int, msg string, data interface{}) {
	RespondJSON(w, status, model.APIResponse{
		Success: true,
		Message: msg,
		Data: data,
	})
}

func RespondError(w http.ResponseWriter, status int, msg string, internalErr error) {
	// log the real error for internal debugging
	if internalErr != nil {
		log.Printf("[ERROR] %d: %s | internal: %v", status, msg, internalErr)
	}

	RespondJSON(w, status, model.APIResponse{
		Success: false,
		Message: msg,
	})
}