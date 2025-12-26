package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/kuznetsovmaksim/glossary/database"
	"github.com/kuznetsovmaksim/glossary/models"
)

func setupTestHandler(t *testing.T) (*TermHandler, *database.DB) {
	t.Helper()
	db, err := database.NewDB(":memory:")
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}
	return NewTermHandler(db), db
}

func TestGetAllTerms(t *testing.T) {
	handler, db := setupTestHandler(t)
	defer db.Close()

	// Create test term
	_, err := db.CreateTerm(&models.CreateTermRequest{
		Name:       "Test Term",
		Definition: "Test Definition",
		Category:   "Test Category",
	})
	if err != nil {
		t.Fatalf("Failed to create test term: %v", err)
	}

	req := httptest.NewRequest("GET", "/api/terms", nil)
	w := httptest.NewRecorder()

	handler.GetAllTerms(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var terms []models.Term
	if err := json.NewDecoder(w.Body).Decode(&terms); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if len(terms) != 1 {
		t.Errorf("Expected 1 term, got %d", len(terms))
	}
}

func TestGetTermByID(t *testing.T) {
	handler, db := setupTestHandler(t)
	defer db.Close()

	// Create test term
	term, err := db.CreateTerm(&models.CreateTermRequest{
		Name:       "Test Term",
		Definition: "Test Definition",
		Category:   "Test Category",
	})
	if err != nil {
		t.Fatalf("Failed to create test term: %v", err)
	}

	req := httptest.NewRequest("GET", "/api/terms/1", nil)
	w := httptest.NewRecorder()

	// Set up router to extract ID from URL
	router := mux.NewRouter()
	router.HandleFunc("/api/terms/{id}", handler.GetTermByID)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var result models.Term
	if err := json.NewDecoder(w.Body).Decode(&result); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if result.ID != term.ID {
		t.Errorf("Expected term ID %d, got %d", term.ID, result.ID)
	}
}

func TestCreateTerm(t *testing.T) {
	handler, db := setupTestHandler(t)
	defer db.Close()

	reqBody := models.CreateTermRequest{
		Name:       "New Term",
		Definition: "New Definition",
		Category:   "New Category",
		Source:     "Test Source",
		SourceURL:  "https://test.com",
	}

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/api/terms", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.CreateTerm(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", w.Code)
	}

	var result models.Term
	if err := json.NewDecoder(w.Body).Decode(&result); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if result.Name != reqBody.Name {
		t.Errorf("Expected name %s, got %s", reqBody.Name, result.Name)
	}
}

func TestCreateTermInvalidRequest(t *testing.T) {
	handler, db := setupTestHandler(t)
	defer db.Close()

	// Missing required fields
	reqBody := models.CreateTermRequest{
		Name: "Incomplete Term",
	}

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/api/terms", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.CreateTerm(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestUpdateTerm(t *testing.T) {
	handler, db := setupTestHandler(t)
	defer db.Close()

	// Create test term
	term, err := db.CreateTerm(&models.CreateTermRequest{
		Name:       "Original Name",
		Definition: "Original Definition",
		Category:   "Original Category",
	})
	if err != nil {
		t.Fatalf("Failed to create test term: %v", err)
	}

	updateReq := models.UpdateTermRequest{
		Name:       "Updated Name",
		Definition: "Updated Definition",
		Category:   "Updated Category",
	}

	body, _ := json.Marshal(updateReq)
	req := httptest.NewRequest("PUT", "/api/terms/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router := mux.NewRouter()
	router.HandleFunc("/api/terms/{id}", handler.UpdateTerm)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var result models.Term
	if err := json.NewDecoder(w.Body).Decode(&result); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if result.Name != updateReq.Name {
		t.Errorf("Expected name %s, got %s", updateReq.Name, result.Name)
	}

	if result.ID != term.ID {
		t.Errorf("Expected ID to remain %d, got %d", term.ID, result.ID)
	}
}

func TestDeleteTerm(t *testing.T) {
	handler, db := setupTestHandler(t)
	defer db.Close()

	// Create test term
	term, err := db.CreateTerm(&models.CreateTermRequest{
		Name:       "Term to Delete",
		Definition: "This will be deleted",
		Category:   "Test",
	})
	if err != nil {
		t.Fatalf("Failed to create test term: %v", err)
	}

	req := httptest.NewRequest("DELETE", "/api/terms/1", nil)
	w := httptest.NewRecorder()

	router := mux.NewRouter()
	router.HandleFunc("/api/terms/{id}", handler.DeleteTerm)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("Expected status 204, got %d", w.Code)
	}

	// Verify term is deleted
	_, err = db.GetTermByID(term.ID)
	if err == nil {
		t.Error("Expected error when getting deleted term, got nil")
	}
}

func TestSearchTerms(t *testing.T) {
	handler, db := setupTestHandler(t)
	defer db.Close()

	// Create test terms
	terms := []models.CreateTermRequest{
		{Name: "Docker", Definition: "Container platform", Category: "DevOps"},
		{Name: "Kubernetes", Definition: "Container orchestration", Category: "DevOps"},
		{Name: "API", Definition: "Application Programming Interface", Category: "Software"},
	}

	for i := range terms {
		_, err := db.CreateTerm(&terms[i])
		if err != nil {
			t.Fatalf("Failed to create test term: %v", err)
		}
	}

	// Search by keyword
	req := httptest.NewRequest("GET", "/api/terms?search=Container", nil)
	w := httptest.NewRecorder()

	handler.GetAllTerms(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var results []models.Term
	if err := json.NewDecoder(w.Body).Decode(&results); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if len(results) != 2 {
		t.Errorf("Expected 2 results, got %d", len(results))
	}
}

func TestFilterByCategory(t *testing.T) {
	handler, db := setupTestHandler(t)
	defer db.Close()

	// Create test terms
	terms := []models.CreateTermRequest{
		{Name: "Docker", Definition: "Container platform", Category: "DevOps"},
		{Name: "Kubernetes", Definition: "Container orchestration", Category: "DevOps"},
		{Name: "API", Definition: "Application Programming Interface", Category: "Software"},
	}

	for i := range terms {
		_, err := db.CreateTerm(&terms[i])
		if err != nil {
			t.Fatalf("Failed to create test term: %v", err)
		}
	}

	// Filter by category
	req := httptest.NewRequest("GET", "/api/terms?category=DevOps", nil)
	w := httptest.NewRecorder()

	handler.GetAllTerms(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var results []models.Term
	if err := json.NewDecoder(w.Body).Decode(&results); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if len(results) != 2 {
		t.Errorf("Expected 2 results, got %d", len(results))
	}

	for _, term := range results {
		if term.Category != "DevOps" {
			t.Errorf("Expected category DevOps, got %s", term.Category)
		}
	}
}
