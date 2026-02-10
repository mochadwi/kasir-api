package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"kasir-api/internal/transaction/model"
	"kasir-api/internal/transaction/service"
)

// Handler handles HTTP requests for transactions
type Handler struct {
	service service.Service
}

// New creates a new transaction handler
func New(service service.Service) *Handler {
	return &Handler{service: service}
}

// HandleCheckout handles POST /api/checkout endpoint
func (h *Handler) HandleCheckout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondWithError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var req model.CheckoutRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	transaction, err := h.service.CreateTransaction(r.Context(), req.Items)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrInvalidTransaction):
			respondWithError(w, http.StatusBadRequest, err.Error())
		case errors.Is(err, model.ErrProductNotFound):
			respondWithError(w, http.StatusNotFound, err.Error())
		case errors.Is(err, model.ErrInsufficientStock):
			respondWithError(w, http.StatusConflict, err.Error())
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	respondWithJSON(w, http.StatusCreated, transaction)
}

func respondWithJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}

func respondWithError(w http.ResponseWriter, status int, message string) {
	respondWithJSON(w, status, map[string]string{"error": message})
}
