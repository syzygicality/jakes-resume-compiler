package health

import (
	"encoding/json"
	"net/http"
)

func SetupHandlers(mux *http.ServeMux) {
	mux.HandleFunc("GET /health", healthHandler)
	mux.HandleFunc("GET /health/protected", protectedHealthHandler)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(map[string]string{"message": "PONG"})
}

func protectedHealthHandler(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(map[string]string{"message": "PONG"})
}
