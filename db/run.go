package db

import (
  "database/sql"
  "io"

  "github.com/jimmc/golden/base"
)

type Runner struct {
  base.Runner

  // Base name for the test output file; if not set, uses BaseName.
  SetupBaseName string
  // Path to the test output file; if not set, uses SetupBaseName.
  SetupPath string

  db *sql.DB
}

func NewRunner(basename string, callback func(*sql.DB, io.Writer) error) *Runner {
  r := &Runner{}
  r.BaseName = basename
  r.Test = func(baseR *base.Runner) error {
    return callback(r.db, r.OutW)
  }
  return r
}

func (r *Runner) SetupFilePath() string {
  return r.GetFilePath(r.SetupPath, r.SetupBaseName, "setup")
}

func (r *Runner) Setup() error {
  if err := r.Runner.Setup(); err != nil {
    return err
  }
  db, err := DbWithSetupFile(r.SetupFilePath())
  if err != nil {
    return err
  }
  r.db = db
  return nil
}

func (r *Runner) Finish() error {
  if r.db != nil {
    r.db.Close()
  }
  return r.Runner.Finish()
}

// Run loads a setup file into a fresh test database, runs the specified
// callback to produce a test output file, and compares it to the golden file.
func (r *Runner) Run() error {
  if err := r.Setup(); err != nil {
    return err
  }

  // Run the specific test step.
  if err := r.Test(&r.Runner); err != nil {
    return err
  }

  // Check the output to see if we got the right data.
  return r.Finish()
}
