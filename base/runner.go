package base

import (
  "testing"
)

// Runner defines the methods used when running one of our unit tests.
type Runner interface {
  Setup() error
  Act() error
  Finish() error
}

// Run runs a test on the Runner.
func Run(r Runner) error {
  if err := r.Setup(); err != nil {
    return err
  }

  // Perform the test action.
  if err := r.Act(); err != nil {
    return err
  }

  // Check the output to see if we got the right data.
  return r.Finish()
}

// RunT is like Run except that it calls t.Fatal on error.
func RunT(t *testing.T, r Runner) {
  t.Helper()
  if err := Run(r); err != nil {
    t.Fatalf("Error running Run: %v", err)
  }
}
