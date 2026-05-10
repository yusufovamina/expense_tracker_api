package response

import (
	"encoding/json"
	"log"
	"net/http"
)

func Error(w http.ResponseWriter, message string, statusCode int) {
	log.Printf("[ERROR] %d: %s", statusCode, message)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

func JSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}
