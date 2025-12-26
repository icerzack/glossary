package models

import "time"

// Term represents a glossary term with its definition and metadata
type Term struct {
	ID         int64     `json:"id"`
	Name       string    `json:"name"`
	Definition string    `json:"definition"`
	Category   string    `json:"category"`
	Source     string    `json:"source"`
	SourceURL  string    `json:"source_url"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// Relationship represents a semantic relationship between two terms
type Relationship struct {
	ID           int64     `json:"id"`
	SourceTermID int64     `json:"source_term_id"`
	TargetTermID int64     `json:"target_term_id"`
	Type         string    `json:"type"` // e.g., "related_to", "part_of", "depends_on", "implements"
	Description  string    `json:"description"`
	CreatedAt    time.Time `json:"created_at"`
}

// GraphNode represents a node in the semantic graph
type GraphNode struct {
	ID         int64  `json:"id"`
	Name       string `json:"name"`
	Definition string `json:"definition"`
	Category   string `json:"category"`
}

// GraphEdge represents an edge in the semantic graph
type GraphEdge struct {
	ID          int64  `json:"id"`
	Source      int64  `json:"source"`
	Target      int64  `json:"target"`
	Type        string `json:"type"`
	Description string `json:"description"`
}

// Graph represents the complete semantic graph
type Graph struct {
	Nodes []GraphNode `json:"nodes"`
	Edges []GraphEdge `json:"edges"`
}

// CreateTermRequest represents the request body for creating a term
type CreateTermRequest struct {
	Name       string `json:"name"`
	Definition string `json:"definition"`
	Category   string `json:"category"`
	Source     string `json:"source"`
	SourceURL  string `json:"source_url"`
}

// UpdateTermRequest represents the request body for updating a term
type UpdateTermRequest struct {
	Name       string `json:"name"`
	Definition string `json:"definition"`
	Category   string `json:"category"`
	Source     string `json:"source"`
	SourceURL  string `json:"source_url"`
}

// CreateRelationshipRequest represents the request body for creating a relationship
type CreateRelationshipRequest struct {
	SourceTermID int64  `json:"source_term_id"`
	TargetTermID int64  `json:"target_term_id"`
	Type         string `json:"type"`
	Description  string `json:"description"`
}
