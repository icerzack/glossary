package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/kuznetsovmaksim/glossary/database"
	"github.com/kuznetsovmaksim/glossary/models"
)

type RelationshipHandler struct {
	db *database.DB
}

func NewRelationshipHandler(db *database.DB) *RelationshipHandler {
	return &RelationshipHandler{db: db}
}

// GetAllRelationships godoc
// @Summary Get all relationships
// @Description Get all term relationships
// @Tags relationships
// @Accept json
// @Produce json
// @Success 200 {array} models.Relationship
// @Failure 500 {string} string "Internal server error"
// @Router /relationships [get]
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

// CreateRelationship godoc
// @Summary Create a new relationship
// @Description Create a relationship between two terms
// @Tags relationships
// @Accept json
// @Produce json
// @Param relationship body models.CreateRelationshipRequest true "Relationship data"
// @Success 201 {object} models.Relationship
// @Failure 400 {string} string "Invalid request body or missing required fields"
// @Failure 500 {string} string "Internal server error"
// @Router /relationships [post]
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

// DeleteRelationship godoc
// @Summary Delete a relationship
// @Description Delete a relationship between terms
// @Tags relationships
// @Accept json
// @Produce json
// @Param id path int true "Relationship ID"
// @Success 204 "Relationship deleted successfully"
// @Failure 400 {string} string "Invalid relationship ID"
// @Failure 500 {string} string "Internal server error"
// @Router /relationships/{id} [delete]
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

// GetGraph godoc
// @Summary Get semantic graph
// @Description Get the complete semantic graph with all terms and relationships
// @Tags graph
// @Accept json
// @Produce json
// @Success 200 {object} models.Graph
// @Failure 500 {string} string "Internal server error"
// @Router /graph [get]
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
