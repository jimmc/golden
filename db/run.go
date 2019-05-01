package db

import (
  "bufio"
  "database/sql"
  "io"
  "os"

  "github.com/jimmc/golden/base"
)

// FromSetupToGolden loads a setup file into a fresh test database, runs the specified
// callback to produce a test output file, and compares it to the golden file.
func FromSetupToGolden(basename string, callback func(*sql.DB, io.Writer) error) error {
  setupfilename := "testdata/" + basename + ".setup"
  r := &base.GoldenRunner{
    BaseName: basename,
    CreateContext: func() (base.GoldenContext, error) {
      db, err := DbWithSetupFile(string(setupfilename))
      if err != nil {
        return nil, err
      }
      return db, err
    },
    Test: func(ctx base.GoldenContext, f *os.File) error {
      w := bufio.NewWriter(f)
      defer w.Flush()
      db := ctx.(*sql.DB)
      return callback(db, w)
    },
  }
  return r.Run()
}
