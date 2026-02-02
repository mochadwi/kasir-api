package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"kasir-api/internal/category/model"
	"kasir-api/internal/category/service"
)

// Handler handles HTTP requests for categories
type Handler struct {
	service service.Service
}

// New creates a new category handler
func New(service service.Service) *Handler {
	return &Handler{service: service}
}

// Handle handles /categories endpoint (list and create)
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

// HandleWithID handles /categories/{id} endpoint (get, update, delete)
func (h *Handler) HandleWithID(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDFromPath("/categories/", r.URL.Path)
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
	categories, err := h.service.ListCategories(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, categories)
}

func (h *Handler) create(w http.ResponseWriter, r *http.Request) {
	var c model.Category
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request")
		return
	}

	created, err := h.service.CreateCategory(r.Context(), c)
	if err != nil {
		if err == model.ErrInvalidCategory {
			respondWithError(w, http.StatusBadRequest, err.Error())
			return
		}
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, created)
}

func (h *Handler) getByID(w http.ResponseWriter, r *http.Request, id int64) {
	c, err := h.service.GetCategory(r.Context(), id)
	if err != nil {
		if err == model.ErrCategoryNotFound {
			respondWithError(w, http.StatusNotFound, "category not found")
			return
		}
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, c)
}

func (h *Handler) update(w http.ResponseWriter, r *http.Request, id int64) {
	var c model.Category
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request")
		return
	}

	updated, err := h.service.UpdateCategory(r.Context(), id, c)
	if err != nil {
		if err == model.ErrCategoryNotFound {
			respondWithError(w, http.StatusNotFound, "category not found")
			return
		}
		if err == model.ErrInvalidCategory {
			respondWithError(w, http.StatusBadRequest, err.Error())
			return
		}
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, updated)
}

func (h *Handler) delete(w http.ResponseWriter, r *http.Request, id int64) {
	if err := h.service.DeleteCategory(r.Context(), id); err != nil {
		if err == model.ErrCategoryNotFound {
			respondWithError(w, http.StatusNotFound, "category not found")
			return
		}
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"message": "category deleted"})
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
