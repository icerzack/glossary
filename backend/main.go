package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/kuznetsovmaksim/glossary/database"
	"github.com/kuznetsovmaksim/glossary/handlers"
	"github.com/rs/cors"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title Glossary API
// @version 1.0
func main() {
	seedFlag := flag.Bool("seed", false, "Seed the database with sample data")
	flag.Parse()

	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "./glossary.db"
	}

	db, err := database.NewDB(dbPath)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer func() {
		if closeErr := db.Close(); closeErr != nil {
			log.Printf("Failed to close database: %v", closeErr)
		}
	}()

	if *seedFlag {
		log.Println("Seeding database with sample data...")
		if err := db.SeedData(); err != nil {
			log.Printf("Failed to seed data: %v", err)
			return
		}
		log.Println("Database seeded successfully")
	}

	termHandler := handlers.NewTermHandler(db)
	relationshipHandler := handlers.NewRelationshipHandler(db)

	r := mux.NewRouter()

	api := r.PathPrefix("/api").Subrouter()

	api.HandleFunc("/terms", termHandler.GetAllTerms).Methods("GET")
	api.HandleFunc("/terms", termHandler.CreateTerm).Methods("POST")
	api.HandleFunc("/terms/{id}", termHandler.GetTermByID).Methods("GET")
	api.HandleFunc("/terms/{id}", termHandler.UpdateTerm).Methods("PUT")
	api.HandleFunc("/terms/{id}", termHandler.DeleteTerm).Methods("DELETE")

	api.HandleFunc("/relationships", relationshipHandler.GetAllRelationships).Methods("GET")
	api.HandleFunc("/relationships", relationshipHandler.CreateRelationship).Methods("POST")
	api.HandleFunc("/relationships/{id}", relationshipHandler.DeleteRelationship).Methods("DELETE")

	api.HandleFunc("/graph", relationshipHandler.GetGraph).Methods("GET")

	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte("OK")); err != nil {
			log.Printf("Failed to write health check response: %v", err)
		}
	}).Methods("GET")

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})

	handler := c.Handler(r)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	server := &http.Server{
		Addr:              ":" + port,
		Handler:           handler,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Printf("Server starting on port %s", port)
	if err := server.ListenAndServe(); err != nil {
		log.Printf("Server failed to start: %v", err)
	}
}
