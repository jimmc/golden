package http

import (
  "database/sql"
  "errors"
  "fmt"
  "io/ioutil"
  "net/http"
  "net/http/httptest"
  "os"

  "github.com/jimmc/golden/base"
  goldendb "github.com/jimmc/golden/db"
)

func SetupToGolden(db *sql.DB, handler http.Handler, basename string,
    callback func() (*http.Request, error)) error {
  return SetupDbToGolden(db, handler, basename, basename, callback)
}

// SetupDbToGolden loads a database setup file into the given database, runs the
// test callback, records the output into the out file, and compares to the golden file.
// All files are located in the testdata folder. The basename arg is used to make
// the filenames for both the output and golden files.
func SetupDbToGolden(db *sql.DB, handler http.Handler, dbsetupbasename, basename string,
    callback func() (*http.Request, error)) error {
  setupfilename := "testdata/" + dbsetupbasename + ".setup"
  outfilename := "testdata/" + basename + ".out"
  goldenfilename := "testdata/" + basename + ".golden"

  if err := goldendb.LoadSetupFile(db, setupfilename); err != nil {
    return fmt.Errorf("error loading setup file %v: %v", setupfilename, err)
  }

  req, err := callback()
  if err != nil {
    return fmt.Errorf("error calling callback in SetupToGolden: %v", err)
  }

  rr := httptest.NewRecorder()
  handler.ServeHTTP(rr, req)

  if got, want := rr.Code, http.StatusOK; got != want {
    return fmt.Errorf("HTTP response status for request %v: got %d, want %d\nBody: %v",
        req.URL, got, want, rr.Body.String())
  }

  body := rr.Body.Bytes()
  if len(body) == 0 {
    return errors.New("response body should not be empty")
  }

  os.Remove(outfilename)
  if err := ioutil.WriteFile(outfilename, body, 0644); err != nil {
    return err
  }

  if err := base.CompareOutToGolden(outfilename, goldenfilename); err != nil {
    return err
  }
  return nil
}
