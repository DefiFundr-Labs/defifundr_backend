## Performance Testing (continued)

### Load Testing

Use tools like k6, JMeter, or custom Go benchmarks to test system performance under load.

Example k6 script:

```js
import http from 'k6/http';
import { check, sleep } from 'k6';

export let options = {
  vus: 100,
  duration: '30s',
};

export default function() {
  let res = http.get('http://localhost:8080/api/v1/users');
  check(res, {
    'status is 200': (r) => r.status === 200,
    'response time < 200ms': (r) => r.timings.duration < 200,
  });
  sleep(1);
}
```

### Benchmarks

Use Go's built-in benchmarking capability for performance-critical code:

```go
func BenchmarkUserRepository_GetUser(b *testing.B) {
    // Setup
    db := setupTestDB()
    repo := NewUserRepository(db)
    
    // Run the benchmark
    for i := 0; i < b.N; i++ {
        user, err := repo.GetUser(context.Background(), 1)
        if err != nil {
            b.Fatal(err)
        }
    }
}
```

Run benchmarks with:

```bash
go test -bench=. ./path/to/package
```

## Continuous Integration

Tests are automatically run in our CI pipeline to ensure code quality.

### CI Workflow

1. **Lint Check**: Verify code style and formatting
   ```bash
   make lint
   ```

2. **Unit Tests**: Run all unit tests
   ```bash
   make test
   ```

3. **Integration Tests**: Run integration tests
   ```bash
   go test -tags=integration ./test/integration
   ```

4. **Coverage Report**: Generate and enforce coverage thresholds
   ```bash
   go test -coverprofile=coverage.out ./...
   go tool cover -func=coverage.out
   ```

### Pre-commit Hooks

We recommend setting up pre-commit hooks to run tests before committing:

```bash
# Add to .git/hooks/pre-commit
#!/bin/sh
go test ./...
if [ $? -ne 0 ]; then
    echo "Tests failed, commit aborted"
    exit 1
fi
```

## Best Practices

### Writing Effective Tests

1. **Arrange-Act-Assert**: Structure tests with clear setup, action, and verification
2. **One Assertion Per Test**: Keep tests focused on a single behavior
3. **Test Edge Cases**: Include tests for boundary conditions and error scenarios
4. **Use Test Helpers**: Create helper functions for common test operations
5. **Descriptive Names**: Use clear test names that describe what is being tested

### Testing Anti-patterns to Avoid

1. **Flaky Tests**: Tests that sometimes pass and sometimes fail
2. **Test Interdependence**: Tests that depend on the state from other tests
3. **Testing Implementation Details**: Focus on testing behavior, not implementation
4. **Slow Tests**: Keep tests fast to encourage frequent running

## Debugging Tests

### Verbose Output

Use the `-v` flag for detailed test output:

```bash
go test -v ./...
```

### Logging in Tests

Use `t.Logf()` to add debug information to tests:

```go
func TestSomething(t *testing.T) {
    result := functionUnderTest()
    t.Logf("Got result: %+v", result)
    assert.Equal(t, expected, result)
}
```

### Skipping Tests

Skip tests that aren't ready or are environment-specific:

```go
func TestFeatureInProgress(t *testing.T) {
    if !isFeatureEnabled() {
        t.Skip("Feature not enabled, skipping test")
    }
    // Test code...
}
```

## Test Data Management

### Test Fixtures

Store test data in `testdata` directories or embedded in test files as constants.

### Factories

Create test data factories for complex objects:

```go
func createTestUser(t *testing.T) *domain.User {
    return &domain.User{
        ID: 1,
        Name: "Test User",
        Email: "test@example.com",
        // ...
    }
}
```

### Database Seeding

For integration tests, seed the database with test data:

```go
func seedTestData(t *testing.T, db *sql.DB) {
    // Execute SQL to insert test data
    _, err := db.Exec(`INSERT INTO users (name, email) VALUES ('Test User', 'test@example.com')`)
    require.NoError(t, err)
}
```

## Mocking With Counterfeiter

### Installation

To use our mock generation capabilities, you need to install counterfeiter:

```bash
go install github.com/maxbrunsfeld/counterfeiter/v6@latest
```

Ensure that your `$GOPATH/bin` is in your `$PATH` so you can execute the counterfeiter command.

### Generating Mocks

We provide a script to generate mocks for all interface dependencies:

```bash
# Make the script executable if needed
chmod +x scripts/generate_mocks.sh

# Generate all mocks
./scripts/generate_mocks.sh
```

This will create mock implementations in the `internal/core/ports/mocks` directory.

### Using Mocks in Tests

Here's an example of how to use the generated mocks in tests:

```go
import (
    "testing"
    "github.com/demola234/defifundr/internal/core/ports/mocks"
    "github.com/stretchr/testify/assert"
)

func TestServiceWithMocks(t *testing.T) {
    // Arrange
    mockRepo := new(mocks.FakeUserRepository)
    
    // Set up expectations
    mockRepo.GetUserByIDReturns(&domain.User{
        ID: uuid.New(),
        Email: "test@example.com",
    }, nil)
    
    // Create service with mock
    service := NewUserService(mockRepo)
    
    // Act
    result, err := service.GetUserByID(context.Background(), userID)
    
    // Assert
    assert.NoError(t, err)
    assert.Equal(t, "test@example.com", result.Email)
    
    // Verify mock was called
    assert.Equal(t, 1, mockRepo.GetUserByIDCallCount())
}
```

### Mocking Patterns

1. **Arrange-Act-Assert with Mocks**:
   - Set up your mock objects and their expectations
   - Execute the code under test
   - Verify the results and mock interactions

2. **Verifying Calls**:
   - `mockObj.MethodNameCallCount()` - Returns number of calls
   - `mockObj.MethodNameArgsForCall(index)` - Returns arguments for a specific call

3. **Stubbing Return Values**:
   - `mockObj.MethodNameReturns(returnValues...)` - Sets return values for all calls
   - `mockObj.MethodNameReturnsOnCall(index, returnValues...)` - Sets return values for a specific call

4. **Testing Error Paths**:
   - `mockObj.MethodNameReturns(nil, errors.New("expected error"))` - Test error handling

Remember to regenerate mocks when interfaces change to ensure they match the latest definitions.

## Conclusion

Thorough testing is essential to maintain high quality in the DefiFundr codebase. By following the guidelines in this document, we can ensure that our tests are effective, maintainable, and provide confidence in our code.

Remember:
- Write tests for new features and bug fixes
- Run the test suite regularly
- Keep tests fast and reliable
- Use mocks appropriately to isolate components
- Maintain high test coverage for critical components