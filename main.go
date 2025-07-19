package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"

	"internal-transfer/handlers"
)

func main() {
	dsn := os.Getenv("POSTGRES_DSN")
	if dsn == "" {
		log.Fatal("POSTGRES_DSN env var is required")
	}

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}
	defer db.Close()

	// Ping to verify connection
	if err := db.Ping(); err != nil {
		log.Fatalf("db ping failed: %v", err)
	}

	h := handlers.NewHandler(db)

	r := mux.NewRouter()
	r.HandleFunc("/accounts", h.CreateAccount).Methods("POST")
	r.HandleFunc("/accounts/{id}", h.GetAccount).Methods("GET")
	r.HandleFunc("/transactions", h.MakeTransaction).Methods("POST")

	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
