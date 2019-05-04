package base

import (
  "bufio"
  "fmt"
  "os"
  "path"
  "testing"
)

// Runner allows for configuring and running the different steps of the test.
type Runner struct {
  // Base name for all files in the test, if not overridden; defaults to "test".
  BaseName string
  // Base directory; if not set, uses "testdata".
  BaseDir string

  // Base name for the test output file; if not set, uses BaseName.
  OutBaseName string
  // Path to the test output file; if not set, uses OutBaseName.
  OutPath string

  // Base name for the golden file; if not set, uses BaseName.
  GoldenBaseName string
  // Path to the golden file; if not set, uses GoldenBaseName.
  GoldenPath string

  // Function to run the test.
  Test func(*Runner) error

  // The output file.
  OutF *os.File;
  // A Writer that can be used to write to the output file.
  OutW *bufio.Writer;
}

func (r *Runner) OutFilePath() string {
  return r.GetFilePath(r.OutPath, r.OutBaseName, "out")
}

func (r *Runner) GoldenFilePath() string {
  return r.GetFilePath(r.GoldenPath, r.GoldenBaseName, "golden")
}

func (r *Runner) GetFilePath(fpath, basename, extension string) string {
  if fpath != "" {
    return fpath
  }
  if basename == "" {
    basename = r.BaseName
    if basename == "" {
      basename = "test"
    }
  }
  basedir := r.BaseDir
  if basedir == "" {
    basedir = "testdata"
  }
  return path.Join(basedir, basename + "." + extension)
}

func (r *Runner) Setup() error {
  outfilepath := r.OutFilePath()
  os.Remove(outfilepath)
  f, err := os.Create(outfilepath)
  if err != nil {
    return fmt.Errorf("error creating output file %q: %v", outfilepath, err)
  }
  w := bufio.NewWriter(f)

  r.OutF = f
  r.OutW = w
  return nil
}

func (r *Runner) Act() error {
  return r.Test(r)
}

func (r *Runner) Finish() error {
  r.OutW.Flush()
  r.OutF.Close()
  return CompareOutToGolden(r.OutFilePath(), r.GoldenFilePath())
}

// SetupT is like Setup except that it calls t.Fatal on error.
func (r *Runner) SetupT(t *testing.T) {
  t.Helper()
  if err := r.Setup(); err != nil {
    t.Fatalf("Error running Setup: %v", err)
  }
}

// FinishT is like Finish except that it calls t.Fatal on error.
func (r *Runner) FinishT(t *testing.T) {
  t.Helper()
  if err := r.Finish(); err != nil {
    t.Fatalf("Error running Finish: %v", err)
  }
}

type GenericRunner interface {
  Setup() error
  Act() error
  Finish() error
}

// Run runs a test using the configuration of the Runner.
func Run(r GenericRunner) error {
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
func RunT(t *testing.T, r GenericRunner) {
  t.Helper()
  if err := Run(r); err != nil {
    t.Fatalf("Error running Run: %v", err)
  }
}
