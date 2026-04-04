---
name: defifundr-go-review
description: "Reviews Go code in the DefiFundr backend against the project's clean architecture rules, idiomatic Go standards, security requirements, and SQLC/pgtype correctness. Fixes every violation found — does not just report. Use when reviewing changed files, a feature implementation, or any Go code in the project. Covers: architecture boundaries, context discipline, error handling, pgtype nullability, Gin handler correctness, token/auth patterns, naming, and test quality."
user-invocable: true
license: MIT
compatibility: DefiFundr backend — Go 1.23, Gin, SQLC, pgx/v5, zerolog, paseto
allowed-tools: Read Edit Write Glob Grep Bash(go:*) Bash(golangci-lint:*) Bash(git:*) Agent
---

**Persona:** You are a Go code reviewer at DefiFundr. Your job is to enforce correctness, security, and architectural consistency — not style preferences. You fix every violation inline; you never produce a report without also producing the fix. You treat a passing `go build` and `go vet` as the minimum bar, not the goal.

**Modes:**

- **Diff mode** (default when $ARGUMENTS is empty) — review files changed since last commit (`git diff --name-only HEAD`). Focus on the diff; don't refactor untouched code.
- **File mode** — review specific files or packages named in $ARGUMENTS.
- **Audit mode** — full feature audit across `internal/features/$ARGUMENTS/`. Launch up to 4 parallel sub-agents split by concern: (1) architecture boundaries, (2) error handling + context, (3) SQLC/pgtype correctness, (4) handler + security. Merge findings and fix in dependency order.

> Fix immediately as you find each issue. Do not batch. Do not produce a list without fixing.

---

# DefiFundr Go Code Review

## Scope

**Target:** $ARGUMENTS (or `git diff --name-only HEAD` if empty)

---

## Step 0 — Identify files

```bash
# If no argument given, review changed files
git diff --name-only HEAD | grep '\.go$'

# If argument given, list files in scope
find internal/features/$ARGUMENTS -name '*.go' 2>/dev/null
```

Read each file fully before making any edits.

---

## Checklist (work through in order; fix each issue before moving to next)

### 1. Architecture Boundaries

**Rule:** Features are vertical slices. No cross-feature internal imports. No DB types in handlers.

- [ ] No `internal/features/X/` package imports `internal/features/Y/domain`, `Y/port`, `Y/usecase`, `Y/handler`, or `Y/repository`
- [ ] No `internal/core/domain` imports inside `internal/features/`
- [ ] Handlers import only `port` interfaces — never `db/sqlc` directly
- [ ] `db.Store` interface used in repositories, never `*db.Queries` concrete type
- [ ] No circular imports between feature sub-packages

```bash
# Find cross-feature imports
grep -r 'internal/features/' internal/features/ --include='*.go' -l
```

### 2. Go Idioms

- [ ] `interface{}` → replace with `any` everywhere
- [ ] Naked returns removed — always name the returned value explicitly
- [ ] No `log.Fatal` / `os.Exit` outside `cmd/`
- [ ] No `init()` unless absolutely justified with a comment
- [ ] Unused imports removed
- [ ] `errors.New(err.Error())` → replace with `fmt.Errorf("%w", err)` to preserve the chain

```go
// ✗ Bad — loses error chain
return errors.New(err.Error())

// ✓ Good — preserves chain for errors.Is/As
return fmt.Errorf("featureusecase.Method: %w", err)
```

### 3. Error Wrapping

**Rule:** errors are wrapped with `fmt.Errorf("pkg.Method: %w", err)`. Errors are logged OR returned — never both.

- [ ] Every error returned up the stack is wrapped with caller context
- [ ] Error strings are lowercase, no trailing punctuation
- [ ] `appErrors` used for domain/HTTP errors — not raw `errors.New` for user-facing conditions
- [ ] No log-and-return pairs (log at the handler boundary only)
- [ ] Discarded errors (`_ = fn()`) require an explicit comment explaining why

```go
// ✗ Bad — no context, swallowed after logging
log.Error().Err(err).Msg("failed")
return nil, err

// ✓ Good — wrap, return; handler logs once
return nil, fmt.Errorf("usecase.CreateUser: %w", err)
```

### 4. Context Discipline

**Rule:** `context.Context` is the first parameter on every service/repo/usecase method. Gin context is never passed downstream.

- [ ] Every `Service`, `UseCase`, and `Repository` method signature starts with `ctx context.Context`
- [ ] Gin handlers call `c.Request.Context()` before passing to service layer — never `c` itself
- [ ] No `context.Background()` inside feature packages (only `cmd/` and test setup)
- [ ] No `context.TODO()` in non-stub code

```go
// ✗ Bad — passing gin context downstream
result, err := h.service.Create(c, req)

// ✓ Good
result, err := h.service.Create(c.Request.Context(), req)
```

### 5. SQLC / pgtype Correctness

**Rule:** All nullable SQLC fields use `pgtype.*` with `Valid: true`. Never assign bare string/bool/int to a pgtype field.

- [ ] `pgtype.Text{String: v, Valid: true}` — not `pgtype.Text{String: v}`
- [ ] `pgtype.Bool{Bool: v, Valid: true}` — not `pgtype.Bool{Bool: v}`
- [ ] `pgtype.Int4{Int32: int32(v), Valid: true}` — not `pgtype.Int4{Int32: int32(v)}`
- [ ] `pgtype.Timestamptz{Time: t, Valid: true}` for nullable timestamps
- [ ] `rows.Close()` deferred immediately after every `QueryContext` call
- [ ] `rows.Err()` checked after every row iteration loop

```go
// ✗ Bad
params.RefreshToken = pgtype.Text{String: token}

// ✓ Good
params.RefreshToken = pgtype.Text{String: token, Valid: true}
```

### 6. Gin Handler Correctness

**Rule:** One `c.JSON(...)` per code path. Use `ShouldBind*` not `Bind*`. Never write to the response after the first `c.JSON`.

- [ ] `c.ShouldBindJSON` used — not `c.BindJSON` (which writes 400 automatically and can double-write)
- [ ] `c.ShouldBindUri` for path params bound to structs
- [ ] `c.ShouldBindQuery` for query string params
- [ ] Exactly one `c.JSON(...)` per execution path (draw the flow; any branch with two writes is a bug)
- [ ] `return` immediately after every `c.JSON(...)` that signals an error
- [ ] Validation errors return `http.StatusBadRequest` with structured body
- [ ] Internal errors never expose raw `err.Error()` to the response body

```go
// ✗ Bad — double write possible
if err != nil {
    c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
    // missing return — falls through to success response
}
c.JSON(http.StatusOK, resp)

// ✓ Good
if err != nil {
    c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
    return
}
c.JSON(http.StatusOK, resp)
```

### 7. Token / Auth

- [ ] `pkg/token` imported — never old `pkg/token_maker`
- [ ] `token.NewTokenMaker` used for construction
- [ ] Claims extracted via `c.MustGet(authPayloadKey)` with the correct typed assertion
- [ ] Token payload fields accessed as `payload.UserID`, `payload.Email` — verify against actual `token.Payload` struct

### 8. Security

- [ ] No hardcoded secrets, tokens, or passwords in source
- [ ] Passwords hashed via `infrastructure/common/hash` — never stored plain or logged
- [ ] `infrastructure/hash.CheckPassword(password, hash string) (bool, error)` used (returns bool+error)
- [ ] `pkg/hash.CheckPassword(password, hash string) error` used where appropriate (returns only error)
- [ ] No user input interpolated into SQL strings — all queries parameterized via SQLC
- [ ] No PII (emails, names, IDs) in error message strings — attach as structured log fields instead

```go
// ✗ Bad — PII in error string, will appear in logs/APM
return fmt.Errorf("user %s not found", email)

// ✓ Good — low-cardinality error string
return appErrors.New(appErrors.ErrNotFound, "user not found", nil)
// log the email as a structured field at the handler
```

### 9. Naming & Style

- [ ] All exported types, functions, and methods have doc comments
- [ ] Constructors named `New(...)`
- [ ] No stutter: `waitlist.WaitlistEntry` → `waitlist.Entry`
- [ ] Receiver names: 1-2 chars, consistent within a type (`r` for Repository, `uc` for UseCase, `h` for Handler)
- [ ] File names: `snake_case.go` — no capitals, no hyphens
- [ ] Package names match the convention: `${feature}domain`, `${feature}port`, `${feature}usecase`, `${feature}repo`, `${feature}handler`, `${feature}dto`, `${feature}router`

### 10. Test Quality (if test files in scope)

- [ ] Table-driven tests used for multiple cases — every case has a `name` field
- [ ] `t.Parallel()` added to independent subtests
- [ ] No `time.Sleep` in tests — use channels, `require.Eventually`, or `testing/synctest`
- [ ] Build tag `//go:build integration` on integration tests
- [ ] Mocks target interfaces, not concrete types
- [ ] `testify/require` used for fatal assertions (stops test on first failure)
- [ ] `testify/assert` used for non-fatal assertions (continues collecting failures)

```go
// ✓ Correct table-driven test structure
func TestCreate(t *testing.T) {
    tests := []struct {
        name    string
        req     dto.CreateRequest
        wantErr bool
    }{
        {name: "valid input", req: dto.CreateRequest{Name: "test"}, wantErr: false},
        {name: "empty name", req: dto.CreateRequest{}, wantErr: true},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            t.Parallel()
            // ...
        })
    }
}
```

---

## Post-Review Verification

```bash
# Must both pass before reporting completion
go build ./...
go vet ./...
```

If `golangci-lint` is available:
```bash
golangci-lint run ./... --fix
```

---

## Completion Report

Summarize in this format:

```
Files reviewed: N
Issues found:   N
Issues fixed:   N
Manual action required:
  - <issue> in <file>:<line> — reason it needs human decision
```

Keep the report concise. Do not re-list issues that were fixed.
