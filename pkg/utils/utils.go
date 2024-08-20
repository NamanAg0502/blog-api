package utils

import (
	"encoding/json"
	"net/http"
	"os"
)

type Response struct {
	Status  string      `json:"status"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

func WriteJSONResponse(w http.ResponseWriter, data interface{}, statusCode int, message string, err error) {
	resp := Response{
		Status:  http.StatusText(statusCode),
		Message: message,
		Data:    data,
	}

	if err != nil {
		resp.Error = err.Error()
		resp.Data = nil
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if encodeErr := json.NewEncoder(w).Encode(resp); encodeErr != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func WriteErrorResponse(w http.ResponseWriter, err error, statusCode int) {
	resp := Response{
		Status: http.StatusText(statusCode),
		Error:  err.Error(),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if encodeErr := json.NewEncoder(w).Encode(resp); encodeErr != nil {
		http.Error(w, "Failed to encode error response", http.StatusInternalServerError)
	}
}

func GetPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port if not specified
	}
	return port
}
