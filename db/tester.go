package db

import (
  "database/sql"
  "io"

  "github.com/jimmc/golden/base"
)

// Tester provides the structure for running unit tests with database setup files.
type Tester struct {
  base.Tester

  // Base name for the test setup file; if not set, uses BaseName.
  SetupBaseName string
  // Path to the test setup file; if not set, uses SetupBaseName.
  SetupPath string

  DB *sql.DB
}

// NewTester creates a new instance of a Tester that will call the specified
// callback as the test function.
func NewTester(basename string, callback func(*sql.DB, io.Writer) error) *Tester {
  r := &Tester{}
  r.BaseName = basename
  r.Test = func(baseR *base.Tester) error {
    return callback(r.DB, r.OutW)
  }
  return r
}

// SetupFilePath returns the complete path to the setup file.
func (r *Tester) SetupFilePath() string {
  return r.GetFilePath(r.SetupPath, r.SetupBaseName, "setup")
}

// Init initializes our database.
func (r *Tester) Init() error {
  db, err := EmptyDb()
  if err != nil {
    return err
  }
  r.DB = db
  return nil
}

// Arrange prepares the output file and loads the setup file.
func (r *Tester) Arrange() error {
  if err := r.Tester.Arrange(); err != nil {
    return err
  }
  if err := LoadSetupFile(r.DB, r.SetupFilePath()); err != nil {
    return err
  }
  return nil
}

// Close closes the database
func (r *Tester) Close() error {
  if r.DB != nil {
    r.DB.Close()
  }
  return nil
}
