package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	db "github.com/defifundr/defifundr-backend/internal/adapters/secondary/db/postgres/sqlc"
)

func main() {
	// Setup connection to database
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:postgres@localhost:5432/defifundr?sslmode=disable"
	}

	// Connect using pgx
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	conn, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}
	defer conn.Close()

	queries := db.New(conn)

	if err := db.SeedDB(ctx, queries); err != nil {
		log.Fatalf("Error seeding database: %v", err)
	}

	log.Println("Database successfully seeded!")
}
