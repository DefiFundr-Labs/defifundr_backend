package sqlc

import (
	"context"
	"errors"
)

// SeedSize defines the size of the seed operation
type SeedSize string

const (
	SeedSizeSmall  SeedSize = "small"
	SeedSizeMedium SeedSize = "medium"
	SeedSizeLarge  SeedSize = "large"
)

// SeedOptions configures the database seeding operation
type SeedOptions struct {
	Size            SeedSize
	Tables          []string
	PreserveData    []string
	RandomSeed      int64
	CleanBeforeRun  bool
	Verbose         bool
	UserCount       int
	TransactionDays int
}

// DefaultSeedOptions returns default seeding options
func DefaultSeedOptions() SeedOptions {
	return SeedOptions{
		Size:            SeedSizeSmall,
		CleanBeforeRun:  false,
		Verbose:         false,
		UserCount:       10,
		TransactionDays: 30,
	}
}

// Seeder seeds the database with test data
type Seeder struct {
	queries *Queries
	options SeedOptions
}

// NewSeeder creates a new Seeder
func NewSeeder(queries *Queries, options SeedOptions) *Seeder {
	return &Seeder{
		queries: queries,
		options: options,
	}
}

// SeedDB seeds the database with test data
// TODO: implement actual seeding logic
func (s *Seeder) SeedDB(ctx context.Context) error {
	return errors.New("seeder not yet implemented")
}
