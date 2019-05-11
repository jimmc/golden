package db_test

import (
  "database/sql"
  "fmt"
  "io"
  "testing"

  "github.com/jimmc/golden/db"
)

// Example is used as the function to be tested by our testing code.
func example(db *sql.DB, w io.Writer) error {
  sql := "SELECT s, n FROM test ORDER BY s;"
  rows, err := db.Query(sql)
  if err != nil {
    return err
  }
  defer rows.Close()
  for rows.Next() {
    var s string
    var n int
    err := rows.Scan(&s, &n)
    if err != nil {
      return err
    }
    fmt.Fprintf(w, "s=%q, n=%d\n", s, n)
  }
  return rows.Err()
}

func TestSetupFilePath(t *testing.T) {
  gr := &db.Tester{}
  if got, want := gr.SetupFilePath(), "testdata/test.setup"; got != want {
    t.Errorf("SetupFilePath on empty config: got %q, want %q", got, want)
  }

  gr = &db.Tester{
    SetupBaseName: "abc",
  }
  if got, want := gr.SetupFilePath(), "testdata/abc.setup"; got != want {
    t.Errorf("SetupFilePath with name: got %q, want %q", got, want)
  }

  gr = &db.Tester{
    SetupPath: "foo/abc.sql",
  }
  if got, want := gr.SetupFilePath(), "foo/abc.sql"; got != want {
    t.Errorf("SetupFilePath with path: got %q, want %q", got, want)
  }
}

// TestDbTester tests the happy path, where our function under test
// is working as expected.
func TestDbTester(t *testing.T) {
  r := db.NewTester("example", example)
  if err := r.Init(); err != nil {
    t.Fatalf("Error in Init: %v", err)
  }
  if err := r.Arrange(); err != nil {
    t.Fatalf("Error in Arrange: %v", err)
  }
  if err := r.Act(); err != nil {
    t.Fatalf("Error in Act: %v", err)
  }
  if err := r.Assert(); err != nil {
    t.Fatalf("Error in Assert: %v", err)
  }
  if err := r.Close(); err != nil {
    t.Fatalf("Error in Close: %v", err)
  }
}

// TestSetupError tests the case where the setup file is invalid.
func TestSetupError(t *testing.T) {
  r := db.NewTester("bad-setup", example)
  if err := r.Init(); err != nil {
    t.Fatalf("Error in Init: %v", err)
  }
  if err := r.Arrange(); err == nil {
    t.Fatalf("Expected no-such-table error")
  }
  if err := r.Close(); err != nil {
    t.Fatalf("Error in Close: %v", err)
  }
}

// TestGoldenMismatch tests the case where the output does not match
// what we expect to see.
func TestGoldenMismatch(t *testing.T) {
  r := db.NewTester("example-no-match", example)
  r.SetupBaseName = "example"
  if err := r.Init(); err != nil {
    t.Fatalf("Error in Init: %v", err)
  }
  if err := r.Arrange(); err != nil {
    t.Fatalf("Error in Arrange: %v", err)
  }
  if err := r.Act(); err != nil {
    t.Fatalf("Error in Act: %v", err)
  }
  if err := r.Assert(); err == nil {
    t.Fatalf("Expected error due to golden file mismatch")
  }
  if err := r.Close(); err != nil {
    t.Fatalf("Error in Close: %v", err)
  }
}
