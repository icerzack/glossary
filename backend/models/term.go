package models

import "time"

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

type Relationship struct {
	ID           int64     `json:"id"`
	SourceTermID int64     `json:"source_term_id"`
	TargetTermID int64     `json:"target_term_id"`
	Type         string    `json:"type"`
	Description  string    `json:"description"`
	CreatedAt    time.Time `json:"created_at"`
}

type GraphNode struct {
	ID         int64  `json:"id"`
	Name       string `json:"name"`
	Definition string `json:"definition"`
	Category   string `json:"category"`
}

type GraphEdge struct {
	ID          int64  `json:"id"`
	Source      int64  `json:"source"`
	Target      int64  `json:"target"`
	Type        string `json:"type"`
	Description string `json:"description"`
}

type Graph struct {
	Nodes []GraphNode `json:"nodes"`
	Edges []GraphEdge `json:"edges"`
}

type CreateTermRequest struct {
	Name       string `json:"name"`
	Definition string `json:"definition"`
	Category   string `json:"category"`
	Source     string `json:"source"`
	SourceURL  string `json:"source_url"`
}

type UpdateTermRequest struct {
	Name       string `json:"name"`
	Definition string `json:"definition"`
	Category   string `json:"category"`
	Source     string `json:"source"`
	SourceURL  string `json:"source_url"`
}

type CreateRelationshipRequest struct {
	SourceTermID int64  `json:"source_term_id"`
	TargetTermID int64  `json:"target_term_id"`
	Type         string `json:"type"`
	Description  string `json:"description"`
}
