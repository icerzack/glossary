package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/kuznetsovmaksim/glossary/database"
	"github.com/kuznetsovmaksim/glossary/models"
)

type TermHandler struct {
	db *database.DB
}

func NewTermHandler(db *database.DB) *TermHandler {
	return &TermHandler{db: db}
}

// GetAllTerms godoc
// @Summary Get all terms
// @Description Get all glossary terms with optional search and filtering
// @Tags terms
// @Accept json
// @Produce json
// @Param search query string false "Search query for term name or definition"
// @Param category query string false "Filter by category"
// @Success 200 {array} models.Term
// @Failure 500 {string} string "Internal server error"
// @Router /terms [get]
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

// GetTermByID godoc
// @Summary Get term by ID
// @Description Get a specific term by its ID
// @Tags terms
// @Accept json
// @Produce json
// @Param id path int true "Term ID"
// @Success 200 {object} models.Term
// @Failure 400 {string} string "Invalid term ID"
// @Failure 404 {string} string "Term not found"
// @Router /terms/{id} [get]
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

// CreateTerm godoc
// @Summary Create a new term
// @Description Add a new term to the glossary
// @Tags terms
// @Accept json
// @Produce json
// @Param term body models.CreateTermRequest true "Term data"
// @Success 201 {object} models.Term
// @Failure 400 {string} string "Invalid request body or missing required fields"
// @Failure 500 {string} string "Internal server error"
// @Router /terms [post]
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

// UpdateTerm godoc
// @Summary Update a term
// @Description Update an existing term
// @Tags terms
// @Accept json
// @Produce json
// @Param id path int true "Term ID"
// @Param term body models.UpdateTermRequest true "Updated term data"
// @Success 200 {object} models.Term
// @Failure 400 {string} string "Invalid request"
// @Failure 500 {string} string "Internal server error"
// @Router /terms/{id} [put]
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

// DeleteTerm godoc
// @Summary Delete a term
// @Description Delete a term from the glossary
// @Tags terms
// @Accept json
// @Produce json
// @Param id path int true "Term ID"
// @Success 204 "Term deleted successfully"
// @Failure 400 {string} string "Invalid term ID"
// @Failure 500 {string} string "Internal server error"
// @Router /terms/{id} [delete]
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
