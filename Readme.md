# DefiFundr - A Decentralized Payroll Platform

[![GitHub Workflow Status (with event)](https://img.shields.io/github/actions/workflow/status/demola234/deFICrowdFunding-Backend/test.yml)](https://github.com/DefiFundr-Labs/defifundr_backend/actions)
![GitHub go.mod Go version (branch & subdirectory of monorepo)](https://img.shields.io/github/go-mod/go-version/demola234/deFICrowdFunding-Backend/main)
[![GitHub issues](https://img.shields.io/github/issues/demola234/deFICrowdFunding-Backend)](https://github.com/DefiFundr-Labs/defifundr_backend/issues?q=is%3Aissue%20state%3Aopen)
[![GitHub Repo stars](https://img.shields.io/github/stars/demola234/deFICrowdFunding-Backend)](https://github.com/DefiFundr-Labs/defifundr_backend/stargazers)
[![GitHub Forks](https://img.shields.io/github/forks/demola234/deFICrowdFunding-Backend)](https://github.com/DefiFundr-Labs/defifundr_backend/network/members)

## 📋 Table of Contents

- [DefiFundr - A Decentralized Payroll Platform](#defifundr---a-decentralized-payroll-platform)
  - [📋 Table of Contents](#-table-of-contents)
  - [🚀 What is DefiFundr?](#-what-is-defifundr)
    - [🌟 Features](#-features)
  - [🏗️ Architecture](#️-architecture)
    - [Key Components:](#key-components)
  - [📁 Project Structure](#-project-structure)
  - [🛠️ Technologies](#️-technologies)
  - [🏁 Getting Started](#-getting-started)
    - [Prerequisites](#prerequisites)
    - [Installation](#installation)
    - [Environment Setup](#environment-setup)
  - [💻 Development](#-development)
    - [Running the Application](#running-the-application)
    - [Available Commands](#available-commands)
    - [🐚 Run Migration with Shell Commands](#-run-migration-with-shell-commands)
  - [📚 API Documentation](#-api-documentation)
  - [🗄️ Database Management](#️-database-management)
    - [Creating a New Migration](#creating-a-new-migration)
    - [Running Migrations](#running-migrations)
    - [Generate SQL Code](#generate-sql-code)
  - [⛓️ Smart Contracts](#️-smart-contracts)
    - [Contracts](#contracts)
    - [Generating Go Bindings](#generating-go-bindings)
  - [🧪 Testing](#-testing)
    - [Running Tests](#running-tests)
    - [Test Structure](#test-structure)
  - [👥 Contributing](#-contributing)
    - [Development Workflow](#development-workflow)
  - [Contributors](#contributors)
  - [📄 License](#-license)

## 🚀 What is DefiFundr?

DefiFundr is a revolutionary decentralized payroll and invoice management system that bridges the gap between traditional financial systems and blockchain technology. The platform provides a seamless, secure, and transparent solution for businesses to manage employee payments, handle freelancer invoices, and automate salary disbursements across both fiat and cryptocurrency channels.

### 🌟 Features

- **Automated Payroll Management**: Schedule and automate regular salary payments
- **Multi-currency Support**: Pay in both fiat and cryptocurrency
- **Invoice Processing**: Create, manage, and process invoices efficiently
- **Secure Authentication**: PASETO token-based authentication with robust password hashing
- **User Management**: Comprehensive user account management with KYC verification
- **Transaction History**: Detailed tracking of all financial transactions
- **Smart Contract Integration**: Direct interaction with Ethereum-based smart contracts
- **API-First Design**: RESTful API architecture for seamless integration

## 🏗️ Architecture

DefiFundr implements a **Hexagonal Architecture** (also known as Ports and Adapters) to achieve:

- Separation of business logic from external concerns
- Improved testability through clear boundaries
- Greater flexibility in replacing components
- Enhanced maintainability with well-defined interfaces

### Key Components:

- **Core Domain** (internal/core): Business rules and entities
- **Ports** (internal/core/ports): Interface definitions
- **Adapters** (internal/adapters): Implementation of interfaces
- **Infrastructure** (infrastructure/): Cross-cutting concerns

## 📁 Project Structure

```
defifundr_backend/
├── cmd/                        # Application entry points
│   ├── api/                    # API server
│   │   ├── docs/               # Swagger documentation
│   │   └── main.go             # API server entry point
│   └── seed/                   # Database seeding
├── config/                     # Configuration management
├── db/                         # Database related code
│   ├── migrations/             # SQL migrations
│   ├── query/                  # SQL queries
│   └── sqlc/                   # Generated Go code
├── docs/                       # Project documentation
├── infrastructure/             # Cross-cutting concerns
│   ├── common/                 # Shared utilities
│   ├── hash/                   # Password hashing
│   └── middleware/             # HTTP middleware
├── internal/                   # Application core code
│   ├── adapters/               # Ports implementation
│   └── core/                   # Business logic and domain
│       ├── domain/             # Domain models
│       ├── ports/              # Interface definitions
│       └── services/           # Business logic
├── pkg/                        # Reusable packages
├── scripts/                    # Utility scripts
└── test/                       # Test suites
```

## 🛠️ Technologies

- **Go**: Main programming language
- **PostgreSQL**: Primary database
- **Docker**: Containerization
- **SQLC**: Type-safe SQL query generation
- **Goose**: Database migration management
- **Swagger**: API documentation
- **PASETO**: Modern security token framework
- **Solidity**: Smart contract development

## 🏁 Getting Started

### Prerequisites

- Go 1.21+
- Docker and Docker Compose
- PostgreSQL 14+
- Make
- Git

### Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/DefiFundr-Labs/defifundr_backend.git
   cd defifundr_backend
   ```

2. Install required tools:
   ```bash
   make install-tools
   ```

3. Set up the development environment:
   ```bash
   make docker-up
   ```

4. Set up the database:
   ```bash
   make postgres
   make createdb
   make migrate-up
   ```

5. Generate SQL code:
   ```bash
   make sqlc
   ```

6. Run the server:
   ```bash
   make server
   ```

### Environment Setup

Create a `.env` file in the project root:

```
DB_SOURCE=postgres://root:secret@localhost:5433/defi?sslmode=disable
SERVER_ADDRESS=0.0.0.0:8080
TOKEN_SYMMETRIC_KEY=your-secret-key-at-least-32-bytes-long
ACCESS_TOKEN_DURATION=15m
REFRESH_TOKEN_DURATION=24h
```

## 💻 Development

### Running the Application

You can run the application in several ways:

```bash
# Standard mode
make server

# Hot reload mode (recommended for development)
make air

# Using Docker
make docker-up
```

### Available Commands

Run `make help` to see a list of all available commands. Key commands include:

```bash
# Database commands
make postgres         # Start PostgreSQL
make createdb         # Create the database
make dropdb           # Drop the database

# Migration commands
make migrate-up       # Apply migrations
make migrate-down     # Revert migrations
make migrate-create   # Create a new migration

# Development commands
make sqlc             # Generate SQL code
make mock             # Generate mock code
make test             # Run tests
make lint             # Run linter
make swagger          # Generate Swagger documentation
```

### 🐚 Run Migration with Shell Commands
```bash
# Create a new migration
cd scripts
sh create_migration.sh
```

```bash
# Apply all pending migrations
cd scripts
sh migrate_up.sh
```

```bash
# Revert the last migration
cd scripts
sh migrate_down.sh
```

```bash
# Reset Migrations
cd scripts
sh migrate_reset.sh
```
```bash
# Migration Status
cd scripts
sh migrate_status.sh
```

```bash 
# Run migrations up to the latest version
cd scripts
sh migrate.sh
```


## 📚 API Documentation

DefiFundr provides comprehensive API documentation using Swagger.

1. Generate the Swagger documentation:
   ```bash
   make swagger
   ```

2. Access the Swagger UI:
   ```
   http://localhost:8080/swagger/index.html
   ```

The API follows RESTful principles with these main endpoints:

- **Authentication**: `/v1/auth/*` (register, login, refresh, logout)
- **Users**: `/v1/users/*` (user management)
- **Transactions**: `/v1/transactions/*` (payment operations)
- **KYC**: `/v1/kyc/*` (verification processes)

For detailed API specifications, see [API_DOCUMENTATION.md](documentation/API_DOCUMENTATION.md).

## 🗄️ Database Management

DefiFundr uses [Goose](https://github.com/pressly/goose) for database migrations and [SQLC](https://sqlc.dev/) for type-safe SQL queries.

### Creating a New Migration

```bash
make migrate-create
# When prompted, enter a descriptive name
```

### Running Migrations

```bash
# Apply all pending migrations
make migrate-up

# Revert the last migration
make migrate-down

# Check migration status
make migrate-status
```

### Generate SQL Code

After adding or modifying queries in the `db/query/` directory:

```bash
make sqlc
```

For more details, see [DATABASE.md](documentation/DATABASE.md).



### Test Structure

- **Unit Tests**: Located alongside the code being tested
- **Integration Tests**: In the `test/integration/` directory
- **End-to-End Tests**: In the `test/e2e/` directory

For detailed testing information, see [TESTING.md](documentation/TESTING.md).


## 👥 Contributing

We welcome contributions to DefiFundr! Please review our [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines on how to contribute.

### Development Workflow

1. Create a feature branch from `main`
2. Implement your changes with appropriate tests
3. Ensure all tests pass with `make test`
4. Create a pull request following our [PR guidelines](ISSUE_PR_GUIDELINES.md)

## Contributors

<a href="https://github.com/DefiFundr-Labs/defifundr_backend/graphs/contributors">
  <img src="https://contrib.rocks/image?repo=DefiFundr-Labs/defifundr_backend" />
</a>

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

<p align="center">
  <b>DefiFundr - Revolutionizing Payroll with Blockchain Technology</b>
</p>
