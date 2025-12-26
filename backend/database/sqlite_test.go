package database

import (
	"os"
	"testing"

	"github.com/kuznetsovmaksim/glossary/models"
)

func setupTestDB(t *testing.T) *DB {
	t.Helper()
	// Create temporary database
	dbPath := ":memory:"
	db, err := NewDB(dbPath)
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}
	return db
}

func TestNewDB(t *testing.T) {
	db := setupTestDB(t)
	if db == nil {
		t.Fatal("Expected database connection, got nil")
	}
	if err := db.Close(); err != nil {
		t.Errorf("Failed to close database: %v", err)
	}
}

func TestCreateTerm(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	req := models.CreateTermRequest{
		Name:       "Test Term",
		Definition: "Test Definition",
		Category:   "Test Category",
		Source:     "Test Source",
		SourceURL:  "https://test.com",
	}

	term, err := db.CreateTerm(&req)
	if err != nil {
		t.Fatalf("Failed to create term: %v", err)
	}

	if term.Name != req.Name {
		t.Errorf("Expected name %s, got %s", req.Name, term.Name)
	}

	if term.Definition != req.Definition {
		t.Errorf("Expected definition %s, got %s", req.Definition, term.Definition)
	}
}

func TestGetTermByID(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	// Create a term first
	req := models.CreateTermRequest{
		Name:       "Test Term",
		Definition: "Test Definition",
		Category:   "Test Category",
	}

	created, err := db.CreateTerm(&req)
	if err != nil {
		t.Fatalf("Failed to create term: %v", err)
	}

	// Get the term by ID
	term, err := db.GetTermByID(created.ID)
	if err != nil {
		t.Fatalf("Failed to get term by ID: %v", err)
	}

	if term.ID != created.ID {
		t.Errorf("Expected ID %d, got %d", created.ID, term.ID)
	}

	if term.Name != created.Name {
		t.Errorf("Expected name %s, got %s", created.Name, term.Name)
	}
}

func TestGetAllTerms(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	// Create multiple terms
	terms := []models.CreateTermRequest{
		{Name: "Term 1", Definition: "Definition 1", Category: "Category 1"},
		{Name: "Term 2", Definition: "Definition 2", Category: "Category 2"},
		{Name: "Term 3", Definition: "Definition 3", Category: "Category 1"},
	}

	for i := range terms {
		_, err := db.CreateTerm(&terms[i])
		if err != nil {
			t.Fatalf("Failed to create term: %v", err)
		}
	}

	// Get all terms
	allTerms, err := db.GetAllTerms()
	if err != nil {
		t.Fatalf("Failed to get all terms: %v", err)
	}

	if len(allTerms) != len(terms) {
		t.Errorf("Expected %d terms, got %d", len(terms), len(allTerms))
	}
}

func TestSearchTerms(t *testing.T) {
	db := setupTestDB(t)
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
			t.Fatalf("Failed to create term: %v", err)
		}
	}

	// Search by keyword
	results, err := db.SearchTerms("Container", "")
	if err != nil {
		t.Fatalf("Failed to search terms: %v", err)
	}

	if len(results) != 2 {
		t.Errorf("Expected 2 results, got %d", len(results))
	}

	// Search by category
	results, err = db.SearchTerms("", "DevOps")
	if err != nil {
		t.Fatalf("Failed to search terms by category: %v", err)
	}

	if len(results) != 2 {
		t.Errorf("Expected 2 results, got %d", len(results))
	}
}

func TestUpdateTerm(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	// Create a term
	req := models.CreateTermRequest{
		Name:       "Original Name",
		Definition: "Original Definition",
		Category:   "Original Category",
	}

	created, err := db.CreateTerm(&req)
	if err != nil {
		t.Fatalf("Failed to create term: %v", err)
	}

	// Update the term
	updateReq := models.UpdateTermRequest{
		Name:       "Updated Name",
		Definition: "Updated Definition",
		Category:   "Updated Category",
	}

	updated, err := db.UpdateTerm(created.ID, &updateReq)
	if err != nil {
		t.Fatalf("Failed to update term: %v", err)
	}

	if updated.Name != updateReq.Name {
		t.Errorf("Expected name %s, got %s", updateReq.Name, updated.Name)
	}

	if updated.Definition != updateReq.Definition {
		t.Errorf("Expected definition %s, got %s", updateReq.Definition, updated.Definition)
	}
}

func TestDeleteTerm(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	// Create a term
	req := models.CreateTermRequest{
		Name:       "Term to Delete",
		Definition: "This will be deleted",
		Category:   "Test",
	}

	created, err := db.CreateTerm(&req)
	if err != nil {
		t.Fatalf("Failed to create term: %v", err)
	}

	// Delete the term
	err = db.DeleteTerm(created.ID)
	if err != nil {
		t.Fatalf("Failed to delete term: %v", err)
	}

	// Try to get the deleted term
	_, err = db.GetTermByID(created.ID)
	if err == nil {
		t.Error("Expected error when getting deleted term, got nil")
	}
}

func TestCreateRelationship(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	// Create two terms
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

	// Create relationship
	relReq := models.CreateRelationshipRequest{
		SourceTermID: term1.ID,
		TargetTermID: term2.ID,
		Type:         "related_to",
		Description:  "Test relationship",
	}

	rel, err := db.CreateRelationship(relReq)
	if err != nil {
		t.Fatalf("Failed to create relationship: %v", err)
	}

	if rel.SourceTermID != term1.ID {
		t.Errorf("Expected source term ID %d, got %d", term1.ID, rel.SourceTermID)
	}

	if rel.TargetTermID != term2.ID {
		t.Errorf("Expected target term ID %d, got %d", term2.ID, rel.TargetTermID)
	}
}

func TestGetGraph(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	// Create terms
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

	// Create relationship
	_, err := db.CreateRelationship(models.CreateRelationshipRequest{
		SourceTermID: term1.ID,
		TargetTermID: term2.ID,
		Type:         "related_to",
	})
	if err != nil {
		t.Fatalf("Failed to create relationship: %v", err)
	}

	// Get graph
	graph, err := db.GetGraph()
	if err != nil {
		t.Fatalf("Failed to get graph: %v", err)
	}

	if len(graph.Nodes) != 2 {
		t.Errorf("Expected 2 nodes, got %d", len(graph.Nodes))
	}

	if len(graph.Edges) != 1 {
		t.Errorf("Expected 1 edge, got %d", len(graph.Edges))
	}
}

func TestSeedData(t *testing.T) {
	// Create temporary file for database
	tmpFile, err := os.CreateTemp("", "test-*.db")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	tmpFile.Close()

	db, err := NewDB(tmpFile.Name())
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}
	defer db.Close()

	// Seed data
	err = db.SeedData()
	if err != nil {
		t.Fatalf("Failed to seed data: %v", err)
	}

	// Verify data was seeded
	terms, err := db.GetAllTerms()
	if err != nil {
		t.Fatalf("Failed to get terms: %v", err)
	}

	if len(terms) == 0 {
		t.Error("Expected seeded terms, got none")
	}

	// Try seeding again - should not duplicate
	err = db.SeedData()
	if err != nil {
		t.Fatalf("Failed to seed data second time: %v", err)
	}

	termsAfter, err := db.GetAllTerms()
	if err != nil {
		t.Fatalf("Failed to get terms after second seed: %v", err)
	}

	if len(termsAfter) != len(terms) {
		t.Errorf("Expected %d terms after second seed, got %d", len(terms), len(termsAfter))
	}
}
