# Database Seeder for DefiFundr

This tool provides a flexible mechanism for seeding the DefiFundr database with realistic test data. It's designed to support various development and testing scenarios, from small datasets for quick tests to large datasets for performance testing.

## Features

- Seed all database entities with realistic data
- Support for different data volumes (small, medium, large)
- Flexible configuration through command-line options
- Selective seeding of specific tables
- Preservation of existing data in selected tables
- Deterministic seeding for reproducible test data
- Comprehensive console output showing progress

## Entities Seeded

The seeder can create data for the following entities:

- **Users**: User accounts with various roles and attributes
- **Sessions**: User authentication sessions
- **OTP Verifications**: One-time password verification records
- **KYC Records**: Know Your Customer verification data
- **User Devices**: User device information
- **Transactions**: Financial transaction records

## Usage

### Basic Usage

```bash
# Seed database with default options (medium size)
make seed

# Or run directly
go run cmd/seed/main.go
```

### Command-Line Options

The seeder supports various command-line options for more control:

```bash
# Specify data volume
go run cmd/seed/main.go -size=small|medium|large

# Seed only specific tables
go run cmd/seed/main.go -tables=users,transactions

# Preserve existing data in certain tables
go run cmd/seed/main.go -preserve=users

# Skip cleaning before seeding
go run cmd/seed/main.go -clean=false

# Generate a specific number of users
go run cmd/seed/main.go -users=20

# Generate transactions for a specific period
go run cmd/seed/main.go -tx-days=60

# Use a specific random seed for reproducible data
go run cmd/seed/main.go -seed=12345

# Disable verbose output
go run cmd/seed/main.go -verbose=false
```

### Example Scenarios

#### Development Setup

For initial development setup with a minimal dataset:

```bash
go run cmd/seed/main.go -size=small -clean=true
```

#### QA Testing

For QA testing with a comprehensive dataset:

```bash
go run cmd/seed/main.go -size=large -clean=true
```

#### Performance Testing

For performance testing with a large number of users and transactions:

```bash
go run cmd/seed/main.go -size=large -users=100 -tx-days=180
```

#### Add Data Without Removing Existing Records

To add more data without clearing existing records:

```bash
go run cmd/seed/main.go -clean=false
```

#### Refresh Specific Tables

To refresh only specific tables while preserving others:

```bash
go run cmd/seed/main.go -tables=transactions,user_device_tokens -preserve=users,sessions
```

## Data Profiles

### Small Profile

- 5 users
- 1-3 sessions per user
- Few OTP verifications 
- Limited transaction history (7 days)
- Suitable for quick tests and development

### Medium Profile (Default)

- 10 users with various roles
- 1-3 sessions per user
- Moderate number of OTP verifications
- Reasonable transaction history (30 days)
- Suitable for most development and testing needs

### Large Profile

- 50+ users with diverse characteristics
- Multiple sessions per user
- Comprehensive OTP verification history
- Extensive transaction history (90 days)
- Suitable for performance testing and complex scenarios

## Entity Relationships

The seeder maintains proper relationships between entities:

```
Users
  └── Sessions
  └── OTP Verifications
  └── KYC Records
  └── User Devices
  └── Transactions
```

## Implementation Details

- Uses deterministic random generation for reproducibility
- Follows proper database constraints and business rules
- Respects referential integrity between tables
- Creates realistic data patterns and distributions

## Development

To extend the seeder for new entities or fields:

1. Update the entity seeding function in `db/sqlc/seed.go`
2. Add the new table to the seeding order in `SeedDB()`
3. Update the table list in `getTablesToSeed()`
4. Include the table in the clean-up list in `cleanTables()` 