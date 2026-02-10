package handler

import (
	"encoding/json"
	"net/http"

	"kasir-api/internal/report/service"
)

type Handler struct {
	service service.Service
}

func New(service service.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) HandleGetReport(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondWithError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var startDate, endDate *string

	if sd := r.URL.Query().Get("start_date"); sd != "" {
		startDate = &sd
	}

	if ed := r.URL.Query().Get("end_date"); ed != "" {
		endDate = &ed
	}

	report, err := h.service.GetSalesReport(r.Context(), startDate, endDate)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, report)
}

func respondWithJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}

func respondWithError(w http.ResponseWriter, status int, message string) {
	respondWithJSON(w, status, map[string]string{"error": message})
}
