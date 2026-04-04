---
name: defifundr-new-feature
description: "Scaffold a complete, production-ready feature for the DefiFundr backend following the project's clean architecture. Creates all 7 sub-packages (domain, port, dto, repository, usecase, handler, router) with correct package naming, import paths, interface contracts, stub implementations for missing SQLC queries, and Gin handler patterns. Use whenever adding a new feature to internal/features/."
user-invocable: true
license: MIT
compatibility: DefiFundr backend — Go 1.23, Gin, SQLC, pgx/v5, zerolog
allowed-tools: Read Edit Write Glob Grep Bash(go:*) Bash(git:*) Agent
---

**Persona:** You are a Go architect who enforces clean architecture boundaries. You design features as self-contained vertical slices — every feature owns its domain, interfaces, data layer, and HTTP surface. You never allow cross-feature internal imports, never leak DB types into handlers, and never skip stub annotations for unimplemented queries.

**Modes:**

- **Scaffold mode** (default) — generate all 7 packages for a new feature end-to-end, build-verify after each file group.
- **Extend mode** — add a new method to an existing feature's port + all layers; triggered when $ARGUMENTS contains an existing feature name + method description.
- **Audit mode** — scan an existing feature for architecture violations (cross-feature imports, missing context params, naked `interface{}`); launch up to 3 parallel sub-agents by concern.

---

# DefiFundr Feature Scaffold

**Feature name:** $ARGUMENTS

Read `internal/features/auth/` as the canonical reference before writing anything.

---

## Pre-flight Checks (run first)

```bash
# Confirm feature doesn't already exist
ls internal/features/$ARGUMENTS/ 2>/dev/null && echo "EXISTS — use extend mode"

# Confirm build is clean before starting
go build ./...
```

---

## Architecture Rules (non-negotiable)

### Module & import paths
- Module root: `github.com/demola234/defifundr`
- Token: `token "github.com/demola234/defifundr/pkg/token"`
- Errors: `appErrors "github.com/demola234/defifundr/pkg/apperrors"`
- DB store: `db "github.com/demola234/defifundr/db/sqlc"`
- Config: `"github.com/demola234/defifundr/config"`
- **Never** import `internal/core/domain` from any feature package
- **Never** import another feature's `domain/`, `port/`, `usecase/`, `handler/`, or `repository/` package

### Package naming
| Directory | Package name |
|---|---|
| `internal/features/$ARGUMENTS/domain/` | `${ARGUMENTS}domain` |
| `internal/features/$ARGUMENTS/port/` | `${ARGUMENTS}port` |
| `internal/features/$ARGUMENTS/usecase/` | `${ARGUMENTS}usecase` |
| `internal/features/$ARGUMENTS/repository/` | `${ARGUMENTS}repo` |
| `internal/features/$ARGUMENTS/handler/` | `${ARGUMENTS}handler` |
| `internal/features/$ARGUMENTS/dto/` | `${ARGUMENTS}dto` |
| `internal/features/$ARGUMENTS/router/` | `${ARGUMENTS}router` |

### Go standards
- No `interface{}` — use `any`
- No naked returns
- Error wrapping: `fmt.Errorf("${ARGUMENTS}usecase.MethodName: %w", err)`
- `context.Context` is always the first parameter on service/repo/usecase methods
- Gin handlers pass `c.Request.Context()` — never `c` itself — to service calls
- All exported types have a doc comment
- Constructors named `New(...)` returning concrete type or interface
- Unexported implementation structs, exported constructor and interfaces
- HTTP handlers call `c.JSON(...)` exactly once per code path
- Logging via `zerolog` — `log.Error().Err(err).Str("feature", "$ARGUMENTS").Msg("...")`

### pgtype correctness for SQLC nullable fields
```go
// ✓ Correct
pgtype.Text{String: value, Valid: true}
pgtype.Bool{Bool: value, Valid: true}
pgtype.Int4{Int32: int32(value), Valid: true}

// ✗ Wrong — always set Valid
pgtype.Text{String: value}
```

### Stub policy for missing SQLC queries
```go
// TODO: add SQLC query — queries/$ARGUMENTS.sql
func (r *Repository) SomeMethod(ctx context.Context, ...) (*domain.Entity, error) {
    return nil, appErrors.New(appErrors.ErrNotImplemented, "not implemented", nil)
}
```

---

## Files to Create (in order)

### 1. `internal/features/$ARGUMENTS/domain/${ARGUMENTS}.go`

Core domain structs. Zero dependencies on other internal packages.

```go
// Package ${ARGUMENTS}domain contains the core domain types for the $ARGUMENTS feature.
package ${ARGUMENTS}domain

import "time"

// Entity represents a ... in the DefiFundr system.
type Entity struct {
    ID        string
    UserID    string
    CreatedAt time.Time
    UpdatedAt time.Time
}
```

### 2. `internal/features/$ARGUMENTS/dto/${ARGUMENTS}_dto.go`

Request and response DTOs. Import only stdlib and the feature's domain package.

```go
package ${ARGUMENTS}dto

// CreateRequest contains the fields required to create a new $ARGUMENTS.
type CreateRequest struct {
    // Add fields with json + binding tags
    Name string `json:"name" binding:"required"`
}

// Response is the API representation of a $ARGUMENTS.
type Response struct {
    ID        string `json:"id"`
    CreatedAt string `json:"created_at"`
}
```

### 3. `internal/features/$ARGUMENTS/port/${ARGUMENTS}_port.go`

Interfaces only — no implementation, no DB types.

```go
// Package ${ARGUMENTS}port defines the interface contracts for the $ARGUMENTS feature.
package ${ARGUMENTS}port

import (
    "context"

    "${ARGUMENTS}domain"
    "${ARGUMENTS}dto"
)

// Repository defines the data access contract for $ARGUMENTS.
type Repository interface {
    Create(ctx context.Context, req dto.CreateRequest) (*domain.Entity, error)
    GetByID(ctx context.Context, id string) (*domain.Entity, error)
    // Add more methods as needed
}

// Service defines the business logic contract for $ARGUMENTS.
type Service interface {
    Create(ctx context.Context, req dto.CreateRequest) (*dto.Response, error)
    GetByID(ctx context.Context, id string) (*dto.Response, error)
}
```

### 4. `internal/features/$ARGUMENTS/repository/${ARGUMENTS}_repository.go`

Implements `port.Repository` against `db.Store`. Stub any method without a SQLC query.

```go
package ${ARGUMENTS}repo

import (
    "context"
    "fmt"

    db "github.com/demola234/defifundr/db/sqlc"
    appErrors "github.com/demola234/defifundr/pkg/apperrors"
    "${ARGUMENTS}domain"
    "${ARGUMENTS}dto"
    "${ARGUMENTS}port"
)

// Repository implements port.Repository using SQLC-generated queries.
type Repository struct {
    store db.Store
}

// Ensure interface satisfaction at compile time.
var _ port.Repository = (*Repository)(nil)

// New constructs a Repository backed by the given store.
func New(store db.Store) *Repository {
    return &Repository{store: store}
}

func (r *Repository) Create(ctx context.Context, req dto.CreateRequest) (*domain.Entity, error) {
    // TODO: add SQLC query — db/queries/${ARGUMENTS}.sql
    return nil, appErrors.New(appErrors.ErrNotImplemented, "not implemented", nil)
}

func (r *Repository) GetByID(ctx context.Context, id string) (*domain.Entity, error) {
    result, err := r.store.Get${ARGUMENTS}ByID(ctx, id)
    if err != nil {
        return nil, fmt.Errorf("${ARGUMENTS}repo.GetByID: %w", err)
    }
    return &domain.Entity{
        ID: result.ID.String(),
        // map remaining fields
    }, nil
}
```

### 5. `internal/features/$ARGUMENTS/usecase/${ARGUMENTS}_usecase.go`

Implements `port.Service`. All business logic here. No direct DB access.

```go
package ${ARGUMENTS}usecase

import (
    "context"
    "fmt"

    "${ARGUMENTS}dto"
    "${ARGUMENTS}port"
)

// UseCase implements port.Service for the $ARGUMENTS feature.
type UseCase struct {
    repo port.Repository
}

// Ensure interface satisfaction at compile time.
var _ port.Service = (*UseCase)(nil)

// New constructs a UseCase with the given repository.
func New(repo port.Repository) *UseCase {
    return &UseCase{repo: repo}
}

func (uc *UseCase) Create(ctx context.Context, req dto.CreateRequest) (*dto.Response, error) {
    entity, err := uc.repo.Create(ctx, req)
    if err != nil {
        return nil, fmt.Errorf("${ARGUMENTS}usecase.Create: %w", err)
    }
    return &dto.Response{
        ID: entity.ID,
    }, nil
}

func (uc *UseCase) GetByID(ctx context.Context, id string) (*dto.Response, error) {
    entity, err := uc.repo.GetByID(ctx, id)
    if err != nil {
        return nil, fmt.Errorf("${ARGUMENTS}usecase.GetByID: %w", err)
    }
    return &dto.Response{ID: entity.ID}, nil
}
```

### 6. `internal/features/$ARGUMENTS/handler/${ARGUMENTS}_handler.go`

Gin handlers. One `c.JSON(...)` per code path. Pass `c.Request.Context()` to service.

```go
package ${ARGUMENTS}handler

import (
    "net/http"

    "github.com/gin-gonic/gin"
    appErrors "github.com/demola234/defifundr/pkg/apperrors"
    "${ARGUMENTS}dto"
    "${ARGUMENTS}port"
)

// Handler holds Gin handler methods for the $ARGUMENTS feature.
type Handler struct {
    service port.Service
}

// New constructs a Handler with the given service.
func New(service port.Service) *Handler {
    return &Handler{service: service}
}

// Create handles POST /$ARGUMENTS.
func (h *Handler) Create(c *gin.Context) {
    var req dto.CreateRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    resp, err := h.service.Create(c.Request.Context(), req)
    if err != nil {
        if appErrors.IsAppError(err) {
            appErr := appErrors.ToAppError(err)
            c.JSON(appErr.HTTPStatus(), gin.H{"error": appErr.Message()})
            return
        }
        c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
        return
    }

    c.JSON(http.StatusCreated, resp)
}

// GetByID handles GET /$ARGUMENTS/:id.
func (h *Handler) GetByID(c *gin.Context) {
    id := c.Param("id")

    resp, err := h.service.GetByID(c.Request.Context(), id)
    if err != nil {
        if appErrors.IsAppError(err) {
            appErr := appErrors.ToAppError(err)
            c.JSON(appErr.HTTPStatus(), gin.H{"error": appErr.Message()})
            return
        }
        c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
        return
    }

    c.JSON(http.StatusOK, resp)
}
```

### 7. `internal/features/$ARGUMENTS/router/${ARGUMENTS}_router.go`

Route registration only. No business logic.

```go
package ${ARGUMENTS}router

import (
    "github.com/gin-gonic/gin"
    "${ARGUMENTS}handler"
)

// RegisterRoutes mounts $ARGUMENTS routes onto the provided router group.
func RegisterRoutes(rg *gin.RouterGroup, h *handler.Handler, authMiddleware gin.HandlerFunc) {
    g := rg.Group("/$ARGUMENTS")
    {
        g.POST("", authMiddleware, h.Create)
        g.GET("/:id", authMiddleware, h.GetByID)
    }
}
```

---

## Verification Steps

After writing all files:

```bash
# 1. Build the feature in isolation
go build ./internal/features/$ARGUMENTS/...

# 2. Full build
go build ./...

# 3. Vet
go vet ./internal/features/$ARGUMENTS/...
```

Fix all errors before reporting completion.

---

## Completion Report

Report:
1. Files created (list with package names)
2. SQLC queries that need to be written (list each missing query with suggested SQL)
3. Any compile errors encountered and how they were resolved
4. **Do NOT wire into `cmd/api/main.go`** unless explicitly asked
