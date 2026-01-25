package models

import (
	"testing"
)

func TestCreateTermRequest(t *testing.T) {
	req := CreateTermRequest{
		Name:       "Test Term",
		Definition: "Test Definition",
		Category:   "Test Category",
		Source:     "Test Source",
		SourceURL:  "https://test.com",
	}

	if req.Name == "" {
		t.Error("Name should not be empty")
	}

	if req.Definition == "" {
		t.Error("Definition should not be empty")
	}

	if req.Category == "" {
		t.Error("Category should not be empty")
	}
}

func TestGraphStructure(t *testing.T) {
	graph := Graph{
		Nodes: []GraphNode{
			{ID: 1, Name: "Node 1", Definition: "Def 1", Category: "Cat 1"},
			{ID: 2, Name: "Node 2", Definition: "Def 2", Category: "Cat 2"},
		},
		Edges: []GraphEdge{
			{ID: 1, Source: 1, Target: 2, Type: "related_to", Description: "Test"},
		},
	}

	if len(graph.Nodes) != 2 {
		t.Errorf("Expected 2 nodes, got %d", len(graph.Nodes))
	}

	if len(graph.Edges) != 1 {
		t.Errorf("Expected 1 edge, got %d", len(graph.Edges))
	}

	if graph.Edges[0].Source != graph.Nodes[0].ID {
		t.Error("Edge source should match first node ID")
	}

	if graph.Edges[0].Target != graph.Nodes[1].ID {
		t.Error("Edge target should match second node ID")
	}
}
