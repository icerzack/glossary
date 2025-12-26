package database

import (
	_ "embed"
	"encoding/json"
	"fmt"
)

//go:embed seed/terms.json
var termsJSON []byte

//go:embed seed/relationships.json
var relationshipsJSON []byte

// SeedTerm represents a term in the seed data
type SeedTerm struct {
	Name       string `json:"name"`
	Definition string `json:"definition"`
	Category   string `json:"category"`
	Source     string `json:"source"`
	SourceURL  string `json:"source_url"`
}

// SeedRelationship represents a relationship in the seed data
type SeedRelationship struct {
	Source      string `json:"source"`
	Target      string `json:"target"`
	Type        string `json:"type"`
	Description string `json:"description"`
}

// SeedData populates the database with sample terms and relationships
func (db *DB) SeedData() error {
	// Check if data already exists
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM terms").Scan(&count)
	if err != nil {
		return err
	}
	if count > 0 {
		return nil // Data already seeded
	}

	// Parse terms JSON
	var seedTerms []SeedTerm
	if err := json.Unmarshal(termsJSON, &seedTerms); err != nil {
		return fmt.Errorf("failed to parse terms JSON: %w", err)
	}

	// Parse relationships JSON
	var seedRelationships []SeedRelationship
	if err := json.Unmarshal(relationshipsJSON, &seedRelationships); err != nil {
		return fmt.Errorf("failed to parse relationships JSON: %w", err)
	}

	// Insert terms and store their IDs
	termIDs := make(map[string]int64)
	for i := range seedTerms {
		term := &seedTerms[i]
		result, err := db.Exec(
			`INSERT INTO terms (name, definition, category, source, source_url) 
			 VALUES (?, ?, ?, ?, ?)`,
			term.Name, term.Definition, term.Category, term.Source, term.SourceURL,
		)
		if err != nil {
			return fmt.Errorf("failed to insert term %s: %w", term.Name, err)
		}
		id, _ := result.LastInsertId()
		termIDs[term.Name] = id
	}

	// Insert relationships
	for i := range seedRelationships {
		rel := &seedRelationships[i]
		sourceID, ok := termIDs[rel.Source]
		if !ok {
			continue
		}
		targetID, ok := termIDs[rel.Target]
		if !ok {
			continue
		}

		_, err := db.Exec(
			`INSERT INTO relationships (source_term_id, target_term_id, type, description) 
			 VALUES (?, ?, ?, ?)`,
			sourceID, targetID, rel.Type, rel.Description,
		)
		if err != nil {
			return fmt.Errorf("failed to insert relationship: %w", err)
		}
	}

	return nil
}
