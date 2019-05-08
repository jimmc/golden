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

// RunT is like Run except that it calls t.Fatal on error.
func RunT(t *testing.T, r Runner) {
  t.Helper()
  if err := Run(r); err != nil {
    t.Fatalf("Error running Run: %v", err)
  }
}
