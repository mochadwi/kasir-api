package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"kasir-api/internal/product/model"
	"kasir-api/internal/product/service"
)

// Handler handles HTTP requests for products
type Handler struct {
	service service.Service
}

// New creates a new product handler
func New(service service.Service) *Handler {
	return &Handler{service: service}
}

// Handle handles /api/produk endpoint (list and create)
func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.list(w, r)
	case http.MethodPost:
		h.create(w, r)
	default:
		respondWithError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

// HandleWithID handles /api/produk/{id} endpoint (get, update, delete)
func (h *Handler) HandleWithID(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDFromPath("/api/produk/", r.URL.Path)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid ID")
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.getByID(w, r, id)
	case http.MethodPut:
		h.update(w, r, id)
	case http.MethodDelete:
		h.delete(w, r, id)
	default:
		respondWithError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

func (h *Handler) list(w http.ResponseWriter, r *http.Request) {
	products, err := h.service.ListProducts(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, products)
}

func (h *Handler) create(w http.ResponseWriter, r *http.Request) {
	var p model.Product
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request")
		return
	}

	created, err := h.service.CreateProduct(r.Context(), p)
	if err != nil {
		if err == model.ErrInvalidProduct {
			respondWithError(w, http.StatusBadRequest, err.Error())
			return
		}
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, created)
}

func (h *Handler) getByID(w http.ResponseWriter, r *http.Request, id int64) {
	p, err := h.service.GetProduct(r.Context(), id)
	if err != nil {
		if err == model.ErrProductNotFound {
			respondWithError(w, http.StatusNotFound, "product not found")
			return
		}
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, p)
}

func (h *Handler) update(w http.ResponseWriter, r *http.Request, id int64) {
	var p model.Product
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request")
		return
	}

	updated, err := h.service.UpdateProduct(r.Context(), id, p)
	if err != nil {
		if err == model.ErrProductNotFound {
			respondWithError(w, http.StatusNotFound, "product not found")
			return
		}
		if err == model.ErrInvalidProduct {
			respondWithError(w, http.StatusBadRequest, err.Error())
			return
		}
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, updated)
}

func (h *Handler) delete(w http.ResponseWriter, r *http.Request, id int64) {
	if err := h.service.DeleteProduct(r.Context(), id); err != nil {
		if err == model.ErrProductNotFound {
			respondWithError(w, http.StatusNotFound, "product not found")
			return
		}
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"message": "product deleted"})
}

func parseIDFromPath(prefix, path string) (int64, error) {
	idStr := strings.TrimPrefix(path, prefix)
	return strconv.ParseInt(idStr, 10, 64)
}

func respondWithJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}

func respondWithError(w http.ResponseWriter, status int, message string) {
	respondWithJSON(w, status, map[string]string{"error": message})
}
