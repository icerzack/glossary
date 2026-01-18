package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/icerzack/glossary/database"
	"github.com/icerzack/glossary/models"
)

func setupTestRelationshipHandler(t *testing.T) (*RelationshipHandler, *database.DB) {
	t.Helper()
	db, err := database.NewDB(":memory:")
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}
	return NewRelationshipHandler(db), db
}

func TestGetAllRelationships(t *testing.T) {
	handler, db := setupTestRelationshipHandler(t)
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("Failed to close database: %v", err)
		}
	}()

	term1, _ := db.CreateTerm(&models.CreateTermRequest{
		Name:       "Term 1",
		Definition: "Definition 1",
		Category:   "Category 1",
	})
	term2, _ := db.CreateTerm(&models.CreateTermRequest{
		Name:       "Term 2",
		Definition: "Definition 2",
		Category:   "Category 2",
	})

	_, err := db.CreateRelationship(models.CreateRelationshipRequest{
		SourceTermID: term1.ID,
		TargetTermID: term2.ID,
		Type:         "related_to",
		Description:  "Test relationship",
	})
	if err != nil {
		t.Fatalf("Failed to create relationship: %v", err)
	}

	req := httptest.NewRequest("GET", "/api/relationships", nil)
	w := httptest.NewRecorder()

	handler.GetAllRelationships(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var relationships []models.Relationship
	if err := json.NewDecoder(w.Body).Decode(&relationships); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if len(relationships) != 1 {
		t.Errorf("Expected 1 relationship, got %d", len(relationships))
	}
}

func TestCreateRelationship(t *testing.T) {
	handler, db := setupTestRelationshipHandler(t)
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("Failed to close database: %v", err)
		}
	}()

	term1, _ := db.CreateTerm(&models.CreateTermRequest{
		Name:       "Term 1",
		Definition: "Definition 1",
		Category:   "Category 1",
	})
	term2, _ := db.CreateTerm(&models.CreateTermRequest{
		Name:       "Term 2",
		Definition: "Definition 2",
		Category:   "Category 2",
	})

	reqBody := models.CreateRelationshipRequest{
		SourceTermID: term1.ID,
		TargetTermID: term2.ID,
		Type:         "depends_on",
		Description:  "Test dependency",
	}

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/api/relationships", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.CreateRelationship(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", w.Code)
	}

	var result models.Relationship
	if err := json.NewDecoder(w.Body).Decode(&result); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if result.Type != reqBody.Type {
		t.Errorf("Expected type %s, got %s", reqBody.Type, result.Type)
	}
}

func TestCreateRelationshipInvalidRequest(t *testing.T) {
	handler, db := setupTestRelationshipHandler(t)
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("Failed to close database: %v", err)
		}
	}()

	reqBody := models.CreateRelationshipRequest{
		SourceTermID: 1,
	}

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/api/relationships", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.CreateRelationship(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestDeleteRelationship(t *testing.T) {
	handler, db := setupTestRelationshipHandler(t)
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("Failed to close database: %v", err)
		}
	}()

	term1, _ := db.CreateTerm(&models.CreateTermRequest{
		Name:       "Term 1",
		Definition: "Definition 1",
		Category:   "Category 1",
	})
	term2, _ := db.CreateTerm(&models.CreateTermRequest{
		Name:       "Term 2",
		Definition: "Definition 2",
		Category:   "Category 2",
	})

	rel, err := db.CreateRelationship(models.CreateRelationshipRequest{
		SourceTermID: term1.ID,
		TargetTermID: term2.ID,
		Type:         "related_to",
	})
	if err != nil {
		t.Fatalf("Failed to create relationship: %v", err)
	}

	req := httptest.NewRequest("DELETE", "/api/relationships/1", nil)
	w := httptest.NewRecorder()

	router := mux.NewRouter()
	router.HandleFunc("/api/relationships/{id}", handler.DeleteRelationship)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("Expected status 204, got %d", w.Code)
	}

	relationships, err := db.GetAllRelationships()
	if err != nil {
		t.Fatalf("Failed to get relationships: %v", err)
	}

	for _, r := range relationships {
		if r.ID == rel.ID {
			t.Error("Relationship should have been deleted")
		}
	}
}

func TestGetGraph(t *testing.T) {
	handler, db := setupTestRelationshipHandler(t)
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("Failed to close database: %v", err)
		}
	}()

	term1, _ := db.CreateTerm(&models.CreateTermRequest{
		Name:       "Term 1",
		Definition: "Definition 1",
		Category:   "Category 1",
	})
	term2, _ := db.CreateTerm(&models.CreateTermRequest{
		Name:       "Term 2",
		Definition: "Definition 2",
		Category:   "Category 2",
	})

	_, err := db.CreateRelationship(models.CreateRelationshipRequest{
		SourceTermID: term1.ID,
		TargetTermID: term2.ID,
		Type:         "related_to",
	})
	if err != nil {
		t.Fatalf("Failed to create relationship: %v", err)
	}

	req := httptest.NewRequest("GET", "/api/graph", nil)
	w := httptest.NewRecorder()

	handler.GetGraph(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var graph models.Graph
	if err := json.NewDecoder(w.Body).Decode(&graph); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if len(graph.Nodes) != 2 {
		t.Errorf("Expected 2 nodes, got %d", len(graph.Nodes))
	}

	if len(graph.Edges) != 1 {
		t.Errorf("Expected 1 edge, got %d", len(graph.Edges))
	}
}
