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

  db *sql.DB
}

// NewTester creates a new instance of a Tester that will call the specified
// callback as the test function.
func NewTester(basename string, callback func(*sql.DB, io.Writer) error) *Tester {
  r := &Tester{}
  r.BaseName = basename
  r.Test = func(baseR *base.Tester) error {
    return callback(r.db, r.OutW)
  }
  return r
}

// SetupFilePath returns the complete path to the setup file.
func (r *Tester) SetupFilePath() string {
  return r.GetFilePath(r.SetupPath, r.SetupBaseName, "setup")
}

// Setup does all of the setup from base.Tester, and sets up the db
// and loads the setup file as calculated from the Tester.
func (r *Tester) Setup() error {
  if err := r.Tester.Setup(); err != nil {
    return err
  }
  db, err := DbWithSetupFile(r.SetupFilePath())
  if err != nil {
    return err
  }
  r.db = db
  return nil
}

// Finish closes the database and the output file and checks the output.
func (r *Tester) Finish() error {
  if r.db != nil {
    r.db.Close()
  }
  return r.Tester.Finish()
}
