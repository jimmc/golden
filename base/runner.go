package base

import (
  "testing"
)

// Runner defines the methods used when running one of our unit tests.
type Runner interface {
  // Arrange sets up for one test, to be performed by Act().
  Arrange() error

  // Act runs the code under test.
  Act() error

  // Assert checks the output against the golden file.
  Assert() error
}

// MultiRunner defines the methods used when running multiple tests that
// include common setup and close steps.
// Typical use for a MultiRunner is to call m.Init(), then do some test-specific
// setup work and call RunTest, repeat that for all tests, and finish by calling m.Close().
type MultiRunner interface {
  Runner

  // Init does one-time initialization of this Runner before running tests.
  Init() error

  // Close cleans everything up. Nothing else can be called after Close.
  Close() error
}

// RunTest runs one test on the Runner.
func RunTest(r Runner) error {
  // Set things up for our one test.
  if err := r.Arrange(); err != nil {
    return err
  }

  // Perform the test action.
  if err := r.Act(); err != nil {
    return err
  }

  // Check the output against the golden file.
  return r.Assert()
}

// Run runs one test on the Runner.
// If the Runner is also a MultiRunner, this also runs the Init and Close functions.
// If not, this is the same as calling RunTest.
func Run(r Runner) error {
  if m, ok := r.(MultiRunner); ok {
    // Do the one-time initialization.
    if err := m.Init(); err != nil {
      return err
    }
  }

  // Run one test.
  if err := RunTest(r); err != nil {
    return err
  }

  if m, ok := r.(MultiRunner); ok {
    // Clean up.
    if err := m.Close(); err != nil {
      return err
    }
  }

  return nil
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
