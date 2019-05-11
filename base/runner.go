package base

import (
  "testing"
)

// Runner defines the methods used when running one of our unit tests.
type Runner interface {
  // Init does one-time initialization of this Runner.
  Init() error

  // Arrange sets up for one test, to be performed by Act().
  Arrange() error

  // Act runs the code under test.
  Act() error

  // Assert checks the output against the golden file.
  Assert() error

  // Close cleans everything up. Nothing else can be called after Close.
  Close() error
}

// Run runs one test on the Runner.
func Run(r Runner) error {
  // Do the one-time initialization.
  if err := r.Init(); err != nil {
    return err
  }

  // Set things up for our one test.
  if err := r.Arrange(); err != nil {
    return err
  }

  // Perform the test action.
  if err := r.Act(); err != nil {
    return err
  }

  // Check the output against the golden file.
  if err := r.Assert(); err != nil {
    return err
  }

  // Clean up.
  return r.Close()
}

// FatalIfError calls testing.T.Fatal if there is an error.
// This is typically used to wrap calls to the various Runner steps, for example:
//   base.FatalIfError(t, r.Arrange(), "Arrange")
func FatalIfError(t *testing.T, err error, label string) {
  t.Helper()
  if err != nil {
    t.Fatalf("Error in %s: %v", label, err)
  }
}
