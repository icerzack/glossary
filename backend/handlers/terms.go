package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/kuznetsovmaksim/glossary/database"
	"github.com/kuznetsovmaksim/glossary/models"
)

// TermHandler handles term-related HTTP requests
type TermHandler struct {
	db *database.DB
}

// NewTermHandler creates a new TermHandler
func NewTermHandler(db *database.DB) *TermHandler {
	return &TermHandler{db: db}
}

// GetAllTerms handles GET /api/terms
func (h *TermHandler) GetAllTerms(w http.ResponseWriter, r *http.Request) {
	search := r.URL.Query().Get("search")
	category := r.URL.Query().Get("category")

	var terms []models.Term
	var err error

	if search != "" || category != "" {
		terms, err = h.db.SearchTerms(search, category)
	} else {
		terms, err = h.db.GetAllTerms()
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if encodeErr := json.NewEncoder(w).Encode(terms); encodeErr != nil {
		http.Error(w, encodeErr.Error(), http.StatusInternalServerError)
	}
}

// GetTermByID handles GET /api/terms/{id}
func (h *TermHandler) GetTermByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid term ID", http.StatusBadRequest)
		return
	}

	term, err := h.db.GetTermByID(id)
	if err != nil {
		http.Error(w, "Term not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if encodeErr := json.NewEncoder(w).Encode(term); encodeErr != nil {
		http.Error(w, encodeErr.Error(), http.StatusInternalServerError)
	}
}

// CreateTerm handles POST /api/terms
func (h *TermHandler) CreateTerm(w http.ResponseWriter, r *http.Request) {
	var req models.CreateTermRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Name == "" || req.Definition == "" || req.Category == "" {
		http.Error(w, "Name, definition, and category are required", http.StatusBadRequest)
		return
	}

	term, err := h.db.CreateTerm(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if encodeErr := json.NewEncoder(w).Encode(term); encodeErr != nil {
		http.Error(w, encodeErr.Error(), http.StatusInternalServerError)
	}
}

// UpdateTerm handles PUT /api/terms/{id}
func (h *TermHandler) UpdateTerm(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid term ID", http.StatusBadRequest)
		return
	}

	var req models.UpdateTermRequest
	if decodeErr := json.NewDecoder(r.Body).Decode(&req); decodeErr != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Name == "" || req.Definition == "" || req.Category == "" {
		http.Error(w, "Name, definition, and category are required", http.StatusBadRequest)
		return
	}

	term, err := h.db.UpdateTerm(id, &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if encodeErr := json.NewEncoder(w).Encode(term); encodeErr != nil {
		http.Error(w, encodeErr.Error(), http.StatusInternalServerError)
	}
}

// DeleteTerm handles DELETE /api/terms/{id}
func (h *TermHandler) DeleteTerm(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid term ID", http.StatusBadRequest)
		return
	}

	if err := h.db.DeleteTerm(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
