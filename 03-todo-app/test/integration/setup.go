package integration

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/chameleon-db/chameleondb/chameleon/pkg/engine"
)

// setupTestEngine creates an engine connected to the test database and empties
// the tables the tests use.
//
// It needs TEST_DATABASE_URL: a database separate from the one you develop
// against, because every test starts by deleting its rows. See the Testing
// section of the README.
func setupTestEngine(t *testing.T) *engine.Engine {
	t.Helper()

	dbURL := os.Getenv("TEST_DATABASE_URL")
	if dbURL == "" {
		t.Skip("TEST_DATABASE_URL is not set; see the Testing section of the README")
	}

	// engine.NewEngine reads .chameleon.yml and the vault from the working
	// directory, and `go test` runs in this package's directory. Move to the
	// project root (this file is in test/integration).
	_, thisFile, _, _ := runtime.Caller(0)
	root := filepath.Join(filepath.Dir(thisFile), "..", "..")
	if err := os.Chdir(root); err != nil {
		t.Fatalf("Failed to change to the project root: %v", err)
	}

	eng, err := engine.NewEngine()
	if err != nil {
		t.Fatalf("Failed to create test engine: %v", err)
	}

	ctx := context.Background()
	cfg, err := engine.ParseConnectionString(dbURL)
	if err != nil {
		t.Fatalf("Invalid TEST_DATABASE_URL: %v", err)
	}
	if err := eng.Connect(ctx, cfg); err != nil {
		t.Fatalf("Failed to connect to the test database: %v", err)
	}
	t.Cleanup(eng.Close)

	if _, err := eng.Connector().Pool().Exec(ctx, "TRUNCATE todos, users"); err != nil {
		t.Fatalf("Failed to empty the test tables: %v", err)
	}

	return eng
}

// TestHelper provides common test utilities
type TestHelper struct {
	T *testing.T
}

// NewTestHelper creates a new test helper
func NewTestHelper(t *testing.T) *TestHelper {
	return &TestHelper{T: t}
}

// AssertError asserts that an error matches expected
func (th *TestHelper) AssertError(err, expected error, message string) {
	if err != expected {
		th.T.Errorf("%s: expected %v, got %v", message, expected, err)
	}
}

// AssertNotNil asserts that a value is not nil
func (th *TestHelper) AssertNotNil(value interface{}, message string) {
	if value == nil {
		th.T.Errorf("%s: expected non-nil value", message)
	}
}

// AssertEqual asserts that two values are equal
func (th *TestHelper) AssertEqual(actual, expected interface{}, message string) {
	if actual != expected {
		th.T.Errorf("%s: expected %v, got %v", message, expected, actual)
	}
}

// AssertTrue asserts that a value is true
func (th *TestHelper) AssertTrue(value bool, message string) {
	if !value {
		th.T.Errorf("%s: expected true, got false", message)
	}
}

// AssertFalse asserts that a value is false
func (th *TestHelper) AssertFalse(value bool, message string) {
	if value {
		th.T.Errorf("%s: expected false, got true", message)
	}
}

// Fatalf logs a fatal error with formatted message
func (th *TestHelper) Fatalf(format string, args ...interface{}) {
	th.T.Fatalf(format, args...)
}

// Errorf logs an error with formatted message
func (th *TestHelper) Errorf(format string, args ...interface{}) {
	th.T.Errorf(format, args...)
}

// Logf logs a message
func (th *TestHelper) Logf(format string, args ...interface{}) {
	th.T.Logf(format, args...)
}

// Printf prints a message
func Printf(format string, args ...interface{}) {
	fmt.Printf(format, args...)
}

// Println prints a message with newline
func Println(args ...interface{}) {
	fmt.Println(args...)
}
