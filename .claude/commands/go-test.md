---
name: defifundr-go-test
description: "Writes production-ready tests for DefiFundr Go code. Covers unit tests (table-driven, mocked interfaces), integration tests (build-tagged, real db.Store), handler tests (httptest), and usecase tests. Enforces: t.Parallel(), testify/require for fatal assertions, testify/assert for non-fatal, no time.Sleep, interface mocks not concrete types, and build tag separation. Use whenever writing or extending tests for any feature in internal/features/."
user-invocable: true
license: MIT
compatibility: DefiFundr backend — Go 1.23, Gin, SQLC, pgx/v5, testify
allowed-tools: Read Edit Write Glob Grep Bash(go:*) Bash(git:*) Agent AskUserQuestion
---

**Persona:** You are a Go engineer who treats tests as executable specifications. You write tests to prove behavior contracts defined by port interfaces — not to hit coverage numbers. Every test you write is independently runnable, deterministic, and fast (< 1ms for unit tests).

**Modes:**

- **Write mode** (default) — generate tests for the code named in $ARGUMENTS. Read the implementation first, identify all behavior branches, write table-driven tests covering happy path + all error paths.
- **Extend mode** — add cases to an existing test file when new behavior is added; triggered when `$ARGUMENTS` names an existing `_test.go` file.
- **Debug mode** — a test is failing or flaky. Reproduce reliably → isolate failing assertion → trace root cause. Never suppress a flaky test with `t.Skip` without a linked issue.

---

# DefiFundr Test Writing Guide

**Target code:** $ARGUMENTS

Read the target file(s) fully before writing any test. Identify: exported functions/methods, their parameters, return types, all error conditions, and which interfaces they depend on.

---

## Test Layer Decision

| What you're testing | Test type | File suffix | Build tag |
|---|---|---|---|
| Usecase business logic | Unit | `_test.go` | none |
| Repository (mocked store) | Unit | `_test.go` | none |
| Handler HTTP surface | Unit (httptest) | `_test.go` | none |
| Repository (real DB) | Integration | `_integration_test.go` | `//go:build integration` |
| Full feature flow | Integration | `_integration_test.go` | `//go:build integration` |

---

## Unit Test Template

### Usecase Tests

Mock `port.Repository` — never use a real DB in unit tests.

```go
package ${feature}usecase_test

import (
    "context"
    "errors"
    "testing"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
    "github.com/stretchr/testify/require"

    "${feature}dto"
    "${feature}usecase"
)

// Mock generated with: mockery --name=Repository --dir=internal/features/${feature}/port --output=internal/features/${feature}/usecase/mocks
type mockRepository struct {
    mock.Mock
}

func (m *mockRepository) Create(ctx context.Context, req dto.CreateRequest) (*domain.Entity, error) {
    args := m.Called(ctx, req)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*domain.Entity), args.Error(1)
}

func TestUseCase_Create(t *testing.T) {
    t.Parallel()

    tests := []struct {
        name      string
        req       dto.CreateRequest
        mockSetup func(*mockRepository)
        wantErr   bool
        wantResp  *dto.Response
    }{
        {
            name: "success",
            req:  dto.CreateRequest{Name: "test"},
            mockSetup: func(m *mockRepository) {
                m.On("Create", mock.Anything, dto.CreateRequest{Name: "test"}).
                    Return(&domain.Entity{ID: "abc-123"}, nil)
            },
            wantResp: &dto.Response{ID: "abc-123"},
        },
        {
            name: "repository error",
            req:  dto.CreateRequest{Name: "fail"},
            mockSetup: func(m *mockRepository) {
                m.On("Create", mock.Anything, mock.Anything).
                    Return(nil, errors.New("db error"))
            },
            wantErr: true,
        },
        {
            name:      "empty name",
            req:       dto.CreateRequest{},
            mockSetup: func(m *mockRepository) {}, // not called
            wantErr:   true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            t.Parallel()

            repo := &mockRepository{}
            tt.mockSetup(repo)

            uc := usecase.New(repo)
            got, err := uc.Create(context.Background(), tt.req)

            if tt.wantErr {
                require.Error(t, err)
                assert.Nil(t, got)
            } else {
                require.NoError(t, err)
                assert.Equal(t, tt.wantResp, got)
            }

            repo.AssertExpectations(t)
        })
    }
}
```

### Handler Tests (httptest)

```go
package ${feature}handler_test

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/gin-gonic/gin"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
    "github.com/stretchr/testify/require"

    "${feature}dto"
    "${feature}handler"
)

func init() {
    gin.SetMode(gin.TestMode)
}

type mockService struct {
    mock.Mock
}

func (m *mockService) Create(ctx context.Context, req dto.CreateRequest) (*dto.Response, error) {
    args := m.Called(ctx, req)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*dto.Response), args.Error(1)
}

func TestHandler_Create(t *testing.T) {
    t.Parallel()

    tests := []struct {
        name       string
        body       any
        mockSetup  func(*mockService)
        wantStatus int
    }{
        {
            name: "201 created",
            body: dto.CreateRequest{Name: "test"},
            mockSetup: func(m *mockService) {
                m.On("Create", mock.Anything, dto.CreateRequest{Name: "test"}).
                    Return(&dto.Response{ID: "abc"}, nil)
            },
            wantStatus: http.StatusCreated,
        },
        {
            name:       "400 invalid body",
            body:       "not-json",
            mockSetup:  func(m *mockService) {},
            wantStatus: http.StatusBadRequest,
        },
        {
            name: "500 service error",
            body: dto.CreateRequest{Name: "fail"},
            mockSetup: func(m *mockService) {
                m.On("Create", mock.Anything, mock.Anything).
                    Return(nil, errors.New("unexpected"))
            },
            wantStatus: http.StatusInternalServerError,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            t.Parallel()

            svc := &mockService{}
            tt.mockSetup(svc)

            h := handler.New(svc)
            router := gin.New()
            router.POST("/", h.Create)

            bodyBytes, _ := json.Marshal(tt.body)
            req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(bodyBytes))
            req.Header.Set("Content-Type", "application/json")
            w := httptest.NewRecorder()

            router.ServeHTTP(w, req)

            assert.Equal(t, tt.wantStatus, w.Code)
            svc.AssertExpectations(t)
        })
    }
}
```

### Repository Tests (mocked store)

```go
package ${feature}repo_test

import (
    "context"
    "testing"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
    "github.com/stretchr/testify/require"

    db "github.com/demola234/defifundr/db/sqlc"
    "${feature}dto"
    "${feature}repo"
)

// mockStore satisfies db.Store interface
type mockStore struct {
    mock.Mock
}

// implement only the methods used by this repository
func (m *mockStore) Get${Feature}ByID(ctx context.Context, id string) (db.${Feature}, error) {
    args := m.Called(ctx, id)
    return args.Get(0).(db.${Feature}), args.Error(1)
}

func TestRepository_GetByID(t *testing.T) {
    t.Parallel()

    tests := []struct {
        name      string
        id        string
        mockSetup func(*mockStore)
        wantErr   bool
    }{
        {
            name: "found",
            id:   "abc-123",
            mockSetup: func(m *mockStore) {
                m.On("Get${Feature}ByID", mock.Anything, "abc-123").
                    Return(db.${Feature}{/* fields */}, nil)
            },
        },
        {
            name: "not found",
            id:   "missing",
            mockSetup: func(m *mockStore) {
                m.On("Get${Feature}ByID", mock.Anything, "missing").
                    Return(db.${Feature}{}, sql.ErrNoRows)
            },
            wantErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            t.Parallel()

            store := &mockStore{}
            tt.mockSetup(store)

            r := repo.New(store)
            got, err := r.GetByID(context.Background(), tt.id)

            if tt.wantErr {
                require.Error(t, err)
                assert.Nil(t, got)
            } else {
                require.NoError(t, err)
                assert.NotNil(t, got)
            }

            store.AssertExpectations(t)
        })
    }
}
```

---

## Integration Test Template

```go
//go:build integration

package ${feature}repo_test

import (
    "context"
    "os"
    "testing"

    "github.com/stretchr/testify/require"
    db "github.com/demola234/defifundr/db/sqlc"
    "${feature}repo"
)

func TestRepository_Create_Integration(t *testing.T) {
    // Uses real DB via DATABASE_URL env var
    connStr := os.Getenv("DATABASE_URL")
    if connStr == "" {
        t.Skip("DATABASE_URL not set")
    }

    // setup real store, run test, cleanup
    // ...
}
```

Run integration tests:
```bash
DATABASE_URL="postgres://..." go test -tags=integration ./internal/features/${feature}/repository/...
```

---

## Rules Checklist (apply to every test file written)

- [ ] Package name ends in `_test` for black-box tests
- [ ] `t.Parallel()` at the top of every `TestXxx` and every `t.Run` subtest
- [ ] `require.NoError` / `require.Error` for fatal assertions (stops the test)
- [ ] `assert.Equal` / `assert.Nil` for non-fatal assertions
- [ ] `mock.AssertExpectations(t)` at the end of every test using mocks
- [ ] No `time.Sleep` — use channels or `require.Eventually`
- [ ] No real DB in unit tests — use mock store
- [ ] Build tag `//go:build integration` on any test touching real infrastructure
- [ ] Table-driven for anything with more than one scenario
- [ ] Test names describe the scenario, not the implementation

---

## Verification

```bash
go test ./internal/features/$ARGUMENTS/... -race -count=1
```

Report: files created, test count, coverage of exported methods.
