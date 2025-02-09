package helpers

import (
	"encoding/json"
	"log"
	"net/http"
)

// ErrorResponse standardizes error responses
type ErrorResponse struct {
	Message string `json:"message"`
}

// SuccessResponse standardizes successful responses
type SuccessResponse struct {
	Data interface{} `json:"data"`
}

// WriteJSONError sends a standardized JSON error response
func WriteJSONError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	resp := ErrorResponse{Message: message}
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("Failed to write JSON error response: %v", err)
	}
}

// WriteJSONResponse sends a standardized JSON success response
func WriteJSONResponse(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	encoder := json.NewEncoder(w)
	encoder.SetEscapeHTML(false)
	if err := encoder.Encode(SuccessResponse{Data: data}); err != nil {
		log.Printf("Failed to write JSON response: %v", err)
		WriteJSONError(w, "Failed to encode response", http.StatusInternalServerError)
	}
}
