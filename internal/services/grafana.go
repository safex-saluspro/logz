package services

import (
	"encoding/json"
	"net/http"
)

type GrafanaData struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

// GrafanaHandler é o endpoint para receber consultas e comandos do Grafana.
func GrafanaHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST requests are allowed", http.StatusMethodNotAllowed)
		return
	}

	// Simulação de resposta para Grafana
	response := GrafanaData{
		Status:  "success",
		Message: "Grafana integration is active",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
