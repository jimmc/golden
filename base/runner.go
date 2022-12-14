package base

import (
  "fmt"
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
// include common init and close steps.
// Typical use for a MultiRunner is to call m.Init(), then do some test-specific
// setup work and call RunTest, repeat that for all tests, and finish by calling m.Close().
type MultiRunner interface {
  Runner

  // Init does one-time initialization of this Runner before running tests.
  Init() error

  // Close cleans everything up. Nothing else can be called after Close.
  Close() error
}

// RunTest runs the test on the Runner by executing the Arrange, Act, and Assert functions.
func RunTest(r Runner) error {
  // Set things up for our one test.
  if err := r.Arrange(); err != nil {
    return fmt.Errorf("error in test Arrange: %v", err)
  }

  // Perform the test action.
  if err := r.Act(); err != nil {
    return fmt.Errorf("error in test Act: %v", err)
  }

  // Check the output against the golden file.
  if err := r.Assert(); err != nil {
    return fmt.Errorf("error in test Assert: %v", err)
  }

  return nil
}

// RunOne runs one test on the MultiRunner by executing
// Init, then running RunTest, then Close.
func RunOne(r MultiRunner) error {
  // Do the one-time initialization.
  if err := r.Init(); err != nil {
    return fmt.Errorf("error in test Init: %v", err)
  }

  // Run one test.
  if err := RunTest(r); err != nil {
    return err
  }

  // Clean up.
  if err := r.Close(); err != nil {
    return fmt.Errorf("error in test Close: %v", err)
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
