package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/kuznetsovmaksim/glossary/database"
	"github.com/kuznetsovmaksim/glossary/models"
)

// RelationshipHandler handles relationship-related HTTP requests
type RelationshipHandler struct {
	db *database.DB
}

// NewRelationshipHandler creates a new RelationshipHandler
func NewRelationshipHandler(db *database.DB) *RelationshipHandler {
	return &RelationshipHandler{db: db}
}

// GetAllRelationships handles GET /api/relationships
func (h *RelationshipHandler) GetAllRelationships(w http.ResponseWriter, r *http.Request) {
	relationships, err := h.db.GetAllRelationships()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if encodeErr := json.NewEncoder(w).Encode(relationships); encodeErr != nil {
		http.Error(w, encodeErr.Error(), http.StatusInternalServerError)
	}
}

// CreateRelationship handles POST /api/relationships
func (h *RelationshipHandler) CreateRelationship(w http.ResponseWriter, r *http.Request) {
	var req models.CreateRelationshipRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.SourceTermID == 0 || req.TargetTermID == 0 || req.Type == "" {
		http.Error(w, "Source term ID, target term ID, and type are required", http.StatusBadRequest)
		return
	}

	relationship, err := h.db.CreateRelationship(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if encodeErr := json.NewEncoder(w).Encode(relationship); encodeErr != nil {
		http.Error(w, encodeErr.Error(), http.StatusInternalServerError)
	}
}

// DeleteRelationship handles DELETE /api/relationships/{id}
func (h *RelationshipHandler) DeleteRelationship(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid relationship ID", http.StatusBadRequest)
		return
	}

	if err := h.db.DeleteRelationship(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// GetGraph handles GET /api/graph
func (h *RelationshipHandler) GetGraph(w http.ResponseWriter, r *http.Request) {
	graph, err := h.db.GetGraph()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if encodeErr := json.NewEncoder(w).Encode(graph); encodeErr != nil {
		http.Error(w, encodeErr.Error(), http.StatusInternalServerError)
	}
}
