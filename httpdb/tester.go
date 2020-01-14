package httpdb

import (
  "errors"
  "fmt"
  "io/ioutil"
  "net/http"
  "net/http/httptest"
  "os"

  goldenbase "github.com/jimmc/golden/base"
  goldendb "github.com/jimmc/golden/db"
)

// Tester provides the structure for running API unit tests.
// For a single test, the typical calling sequence is:
//   r := NewTester(handlerCreateFunc)
//   r.Run(t, basename, callback)
// For multiple tests, maintaining the Tester state across tests as it changes:
//   r := NewTester(handlerCreateFunc)
//   r.Init()
//   r.RunTestWith(t, basename, callback)
//   r.RunTestWith(t, basename2, callback2)
//   r.Close()
type Tester struct {
  goldendb.Tester

  CreateHandler func(r *Tester) http.Handler
  Callback func() (*http.Request, error)
}

type TesterApi interface {
  goldenbase.MultiRunner
  SetBaseNameAndCallback(basename string, callback func() (*http.Request, error))
}

// NewTester creates a new instance of a Tester that will use the specified
// function to create an http.Handler.
func NewTester(createHandler func(r *Tester) http.Handler) *Tester {
  r := &Tester{}
  r.CreateHandler = createHandler
  return r
}

// SetBaseNameAndCallback resets the basename and callback of the Tester in preparation for running a test.
func (r *Tester) SetBaseNameAndCallback(basename string, callback func() (*http.Request, error)) {
  r.BaseName = basename
  r.Callback = callback
}

// Act sets up the handler, calls the request, and records the result to the output file.
func (r *Tester) Act() error {
  handler := r.CreateHandler(r)

  req, err := r.Callback()
  if err != nil {
    return fmt.Errorf("error calling callback in Tester.Act: %v", err)
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

  outfilepath := r.OutFilePath()
  os.Remove(outfilepath)
  if err := ioutil.WriteFile(outfilepath, body, 0644); err != nil {
    return err
  }
  return nil
}

// RunTestWith runs a test using the specified basename and callback.
// This can be used multiple times within a Tester. The database state is maintained across tests,
// allowing a sequence of calls that builds up and modifies a database.
func RunTestWith(r TesterApi, basename string, callback func() (*http.Request, error)) error {
  r.SetBaseNameAndCallback(basename, callback)
  return goldenbase.RunTest(r)
}

// Run initializes the tester, runs a test, and closes it, calling Fatalf on any error.
func RunOneWith(r TesterApi, basename string, callback func() (*http.Request, error)) error {
  if err := r.Init(); err != nil {
    return err
  }
  if err := RunTestWith(r, basename, callback); err != nil {
    return err
  }
  return r.Close()
}
