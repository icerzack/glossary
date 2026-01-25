package database

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/icerzack/glossary/models"
	_ "github.com/mattn/go-sqlite3"
)

type DB struct {
	*sql.DB
}

func NewDB(dataSourceName string) (*DB, error) {
	db, err := sql.Open("sqlite3", dataSourceName)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	database := &DB{db}
	if err := database.createSchema(); err != nil {
		return nil, fmt.Errorf("failed to create schema: %w", err)
	}

	return database, nil
}

func (db *DB) createSchema() error {
	schema := `
	CREATE TABLE IF NOT EXISTS terms (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL UNIQUE,
		definition TEXT NOT NULL,
		category TEXT NOT NULL,
		source TEXT,
		source_url TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS relationships (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		source_term_id INTEGER NOT NULL,
		target_term_id INTEGER NOT NULL,
		type TEXT NOT NULL,
		description TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (source_term_id) REFERENCES terms(id) ON DELETE CASCADE,
		FOREIGN KEY (target_term_id) REFERENCES terms(id) ON DELETE CASCADE,
		UNIQUE(source_term_id, target_term_id, type)
	);

	CREATE INDEX IF NOT EXISTS idx_terms_name ON terms(name);
	CREATE INDEX IF NOT EXISTS idx_terms_category ON terms(category);
	CREATE INDEX IF NOT EXISTS idx_relationships_source ON relationships(source_term_id);
	CREATE INDEX IF NOT EXISTS idx_relationships_target ON relationships(target_term_id);
	`

	_, err := db.Exec(schema)
	return err
}

func (db *DB) GetAllTerms() ([]models.Term, error) {
	rows, err := db.Query(`
		SELECT id, name, definition, category, source, source_url, created_at, updated_at 
		FROM terms 
		ORDER BY name
	`)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = rows.Close()
	}()

	var terms []models.Term
	for rows.Next() {
		var term models.Term
		err := rows.Scan(
			&term.ID, &term.Name, &term.Definition, &term.Category,
			&term.Source, &term.SourceURL, &term.CreatedAt, &term.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		terms = append(terms, term)
	}

	return terms, rows.Err()
}

func (db *DB) GetTermByID(id int64) (*models.Term, error) {
	var term models.Term
	err := db.QueryRow(`
		SELECT id, name, definition, category, source, source_url, created_at, updated_at 
		FROM terms 
		WHERE id = ?
	`, id).Scan(
		&term.ID, &term.Name, &term.Definition, &term.Category,
		&term.Source, &term.SourceURL, &term.CreatedAt, &term.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &term, nil
}

func (db *DB) SearchTerms(query, category string) ([]models.Term, error) {
	var rows *sql.Rows
	var err error

	if category != "" {
		rows, err = db.Query(`
			SELECT id, name, definition, category, source, source_url, created_at, updated_at 
			FROM terms 
			WHERE (name LIKE ? OR definition LIKE ?) AND category = ?
			ORDER BY name
		`, "%"+query+"%", "%"+query+"%", category)
	} else {
		rows, err = db.Query(`
			SELECT id, name, definition, category, source, source_url, created_at, updated_at 
			FROM terms 
			WHERE name LIKE ? OR definition LIKE ?
			ORDER BY name
		`, "%"+query+"%", "%"+query+"%")
	}

	if err != nil {
		return nil, err
	}
	defer func() {
		_ = rows.Close()
	}()

	var terms []models.Term
	for rows.Next() {
		var term models.Term
		err := rows.Scan(
			&term.ID, &term.Name, &term.Definition, &term.Category,
			&term.Source, &term.SourceURL, &term.CreatedAt, &term.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		terms = append(terms, term)
	}

	return terms, rows.Err()
}

func (db *DB) CreateTerm(req *models.CreateTermRequest) (*models.Term, error) {
	result, err := db.Exec(`
		INSERT INTO terms (name, definition, category, source, source_url) 
		VALUES (?, ?, ?, ?, ?)
	`, req.Name, req.Definition, req.Category, req.Source, req.SourceURL)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return db.GetTermByID(id)
}

func (db *DB) UpdateTerm(id int64, req *models.UpdateTermRequest) (*models.Term, error) {
	_, err := db.Exec(`
		UPDATE terms 
		SET name = ?, definition = ?, category = ?, source = ?, source_url = ?, updated_at = ?
		WHERE id = ?
	`, req.Name, req.Definition, req.Category, req.Source, req.SourceURL, time.Now(), id)
	if err != nil {
		return nil, err
	}

	return db.GetTermByID(id)
}

func (db *DB) DeleteTerm(id int64) error {
	_, err := db.Exec("DELETE FROM terms WHERE id = ?", id)
	return err
}

func (db *DB) GetAllRelationships() ([]models.Relationship, error) {
	rows, err := db.Query(`
		SELECT id, source_term_id, target_term_id, type, description, created_at 
		FROM relationships
	`)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = rows.Close()
	}()

	var relationships []models.Relationship
	for rows.Next() {
		var rel models.Relationship
		err := rows.Scan(
			&rel.ID, &rel.SourceTermID, &rel.TargetTermID,
			&rel.Type, &rel.Description, &rel.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		relationships = append(relationships, rel)
	}

	return relationships, rows.Err()
}

func (db *DB) CreateRelationship(req models.CreateRelationshipRequest) (*models.Relationship, error) {
	result, err := db.Exec(`
		INSERT INTO relationships (source_term_id, target_term_id, type, description) 
		VALUES (?, ?, ?, ?)
	`, req.SourceTermID, req.TargetTermID, req.Type, req.Description)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	var rel models.Relationship
	err = db.QueryRow(`
		SELECT id, source_term_id, target_term_id, type, description, created_at 
		FROM relationships 
		WHERE id = ?
	`, id).Scan(
		&rel.ID, &rel.SourceTermID, &rel.TargetTermID,
		&rel.Type, &rel.Description, &rel.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &rel, nil
}

func (db *DB) DeleteRelationship(id int64) error {
	_, err := db.Exec("DELETE FROM relationships WHERE id = ?", id)
	return err
}

func (db *DB) GetGraph() (*models.Graph, error) {
	terms, err := db.GetAllTerms()
	if err != nil {
		return nil, err
	}

	nodes := make([]models.GraphNode, 0, len(terms))
	for i := range terms {
		nodes = append(nodes, models.GraphNode{
			ID:         terms[i].ID,
			Name:       terms[i].Name,
			Definition: terms[i].Definition,
			Category:   terms[i].Category,
		})
	}

	relationships, err := db.GetAllRelationships()
	if err != nil {
		return nil, err
	}

	edges := make([]models.GraphEdge, len(relationships))
	for i, rel := range relationships {
		edges[i] = models.GraphEdge{
			ID:          rel.ID,
			Source:      rel.SourceTermID,
			Target:      rel.TargetTermID,
			Type:        rel.Type,
			Description: rel.Description,
		}
	}

	return &models.Graph{
		Nodes: nodes,
		Edges: edges,
	}, nil
}
