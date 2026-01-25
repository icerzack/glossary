package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/icerzack/glossary/database"
	"github.com/icerzack/glossary/handlers"
	"github.com/rs/cors"
	httpSwagger "github.com/swaggo/http-swagger"
)

const (
	defaultReadTimeout       = 15 * time.Second
	defaultWriteTimeout      = 15 * time.Second
	defaultIdleTimeout       = 60 * time.Second
	defaultReadHeaderTimeout = 5 * time.Second
)

// @title Glossary API
// @version 1.0
func main() {
	seedFlag := flag.Bool("seed", false, "Seed the database with sample data")
	flag.Parse()

	db := initializeDatabase()
	defer func() {
		if closeErr := db.Close(); closeErr != nil {
			log.Printf("Failed to close database: %v", closeErr)
		}
	}()

	if *seedFlag {
		seedDatabase(db)
	}

	server := setupServer(db)
	startServer(server)
}

func initializeDatabase() *database.DB {
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "./glossary.db"
	}

	db, err := database.NewDB(dbPath)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	return db
}

func seedDatabase(db *database.DB) {
	log.Println("Seeding database with sample data...")
	if err := db.SeedData(); err != nil {
		log.Printf("Failed to seed data: %v", err)
		return
	}
	log.Println("Database seeded successfully")
}

func setupServer(db *database.DB) *http.Server {
	termHandler := handlers.NewTermHandler(db)
	relationshipHandler := handlers.NewRelationshipHandler(db)

	r := mux.NewRouter()
	api := r.PathPrefix("/api").Subrouter()

	setupAPIRoutes(api, termHandler, relationshipHandler)

	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)
	r.HandleFunc("/health", healthCheckHandler).Methods("GET")

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	return &http.Server{
		Addr:              ":" + port,
		Handler:           c.Handler(r),
		ReadTimeout:       defaultReadTimeout,
		WriteTimeout:      defaultWriteTimeout,
		IdleTimeout:       defaultIdleTimeout,
		ReadHeaderTimeout: defaultReadHeaderTimeout,
	}
}

func setupAPIRoutes(
	api *mux.Router,
	termHandler *handlers.TermHandler,
	relationshipHandler *handlers.RelationshipHandler,
) {
	api.HandleFunc("/terms", termHandler.GetAllTerms).Methods("GET")
	api.HandleFunc("/terms", termHandler.CreateTerm).Methods("POST")
	api.HandleFunc("/terms/{id}", termHandler.GetTermByID).Methods("GET")
	api.HandleFunc("/terms/{id}", termHandler.UpdateTerm).Methods("PUT")
	api.HandleFunc("/terms/{id}", termHandler.DeleteTerm).Methods("DELETE")

	api.HandleFunc("/relationships", relationshipHandler.GetAllRelationships).Methods("GET")
	api.HandleFunc("/relationships", relationshipHandler.CreateRelationship).Methods("POST")
	api.HandleFunc("/relationships/{id}", relationshipHandler.DeleteRelationship).Methods("DELETE")

	api.HandleFunc("/graph", relationshipHandler.GetGraph).Methods("GET")
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte("OK")); err != nil {
		log.Printf("Failed to write health check response: %v", err)
	}
}

func startServer(server *http.Server) {
	port := server.Addr
	if port != "" && port[0] == ':' {
		port = port[1:]
	}
	log.Printf("Server starting on port %s", port)
	if err := server.ListenAndServe(); err != nil {
		log.Printf("Server failed to start: %v", err)
	}
}
