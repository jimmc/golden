package db

import (
  "database/sql"
  "io"

  "github.com/jimmc/golden/base"
)

type Tester struct {
  base.Tester

  // Base name for the test output file; if not set, uses BaseName.
  SetupBaseName string
  // Path to the test output file; if not set, uses SetupBaseName.
  SetupPath string

  db *sql.DB
}

func NewTester(basename string, callback func(*sql.DB, io.Writer) error) *Tester {
  r := &Tester{}
  r.BaseName = basename
  r.Test = func(baseR *base.Tester) error {
    return callback(r.db, r.OutW)
  }
  return r
}

func (r *Tester) SetupFilePath() string {
  return r.GetFilePath(r.SetupPath, r.SetupBaseName, "setup")
}

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

func (r *Tester) Finish() error {
  if r.db != nil {
    r.db.Close()
  }
  return r.Tester.Finish()
}
