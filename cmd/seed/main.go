package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	db "github.com/demola234/defifundr/db/sqlc"
	"github.com/jackc/pgx/v5/pgxpool"
	// db "github.com/demola234/defifundr/internal/adapters/secondary/db/postgres/sqlc"
)

func main() {
	// Parse command line flags
	var (
		size            string
		tables          string
		preserveData    string
		randomSeed      int64
		cleanBeforeRun  bool
		verbose         bool
		userCount       int
		transactionDays int
	)

	// Define command line flags
	flag.StringVar(&size, "size", "medium", "Size of the dataset to generate (small, medium, large)")
	flag.StringVar(&tables, "tables", "", "Comma-separated list of specific tables to seed (empty means all)")
	flag.StringVar(&preserveData, "preserve", "", "Comma-separated list of tables to preserve existing data")
	flag.Int64Var(&randomSeed, "seed", 0, "Seed for random generator (0 means use time)")
	flag.BoolVar(&cleanBeforeRun, "clean", true, "Clean tables before seeding")
	flag.BoolVar(&verbose, "verbose", true, "Enable verbose output")
	flag.IntVar(&userCount, "users", 0, "Number of users to generate (0 means use size-based default)")
	flag.IntVar(&transactionDays, "tx-days", 0, "How many days of transaction history to generate (0 means use size-based default)")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of database seeder:\n")
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "  %s [options]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "Examples:\n")
		fmt.Fprintf(os.Stderr, "  Seed all tables with medium data volume:\n")
		fmt.Fprintf(os.Stderr, "    %s -size=medium\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "  Seed only users and transactions tables:\n")
		fmt.Fprintf(os.Stderr, "    %s -tables=users,transactions\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "  Seed with large data volume but preserve existing users:\n")
		fmt.Fprintf(os.Stderr, "    %s -size=large -preserve=users\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "  Use a specific random seed for reproducible data:\n")
		fmt.Fprintf(os.Stderr, "    %s -seed=12345\n", os.Args[0])
	}

	flag.Parse()

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

	// Create seed options
	options := db.DefaultSeedOptions()

	// Apply command line options
	if size != "" {
		switch strings.ToLower(size) {
		case "small":
			options.Size = db.SeedSizeSmall
		case "medium":
			options.Size = db.SeedSizeMedium
		case "large":
			options.Size = db.SeedSizeLarge
		default:
			log.Fatalf("Invalid size value: %s. Must be 'small', 'medium', or 'large'", size)
		}
	}

	if tables != "" {
		options.Tables = strings.Split(tables, ",")
	}

	if preserveData != "" {
		options.PreserveData = strings.Split(preserveData, ",")
	}

	options.RandomSeed = randomSeed
	options.CleanBeforeRun = cleanBeforeRun
	options.Verbose = verbose

	if userCount > 0 {
		options.UserCount = userCount
	}

	if transactionDays > 0 {
		options.TransactionDays = transactionDays
	}

	// Create seeder and run
	startTime := time.Now()

	log.Println("Initializing database seeder...")
	log.Printf("Options: size=%s, tables=%v, clean=%v", options.Size, options.Tables, options.CleanBeforeRun)

	seeder := db.NewSeeder(queries, options)

	if err := seeder.SeedDB(ctx); err != nil {
		log.Fatalf("Error seeding database: %v", err)
	}

	duration := time.Since(startTime)
	log.Printf("Database successfully seeded in %s!", duration)
}
