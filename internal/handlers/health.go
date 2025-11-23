package handlers

import (
	"encoding/json"
	"net/http"
)

// HealthResponse represents a simple JSON health response.
type HealthResponse struct {
	Status string `json:"status"`
}

// HealthHandler returns 200 OK with a basic JSON payload {"status":"ok"}.
func HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(HealthResponse{Status: "ok"})
}
