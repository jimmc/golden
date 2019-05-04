package base

import (
  "bufio"
  "fmt"
  "os"
  "path"
  "testing"
)

// Tester allows for configuring and running the different steps of the test.
type Tester struct {
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
  Test func(*Tester) error

  // The output file.
  OutF *os.File;
  // A Writer that can be used to write to the output file.
  OutW *bufio.Writer;
}

func NewTester(basename string) *Tester {
  r := &Tester{
    BaseName: basename,
  }
  return r
}

func (r *Tester) OutFilePath() string {
  return r.GetFilePath(r.OutPath, r.OutBaseName, "out")
}

func (r *Tester) GoldenFilePath() string {
  return r.GetFilePath(r.GoldenPath, r.GoldenBaseName, "golden")
}

func (r *Tester) GetFilePath(fpath, basename, extension string) string {
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

func (r *Tester) Setup() error {
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

func (r *Tester) Act() error {
  return r.Test(r)
}

func (r *Tester) Finish() error {
  r.OutW.Flush()
  r.OutF.Close()
  return CompareOutToGolden(r.OutFilePath(), r.GoldenFilePath())
}

// SetupT is like Setup except that it calls t.Fatal on error.
func (r *Tester) SetupT(t *testing.T) {
  t.Helper()
  if err := r.Setup(); err != nil {
    t.Fatalf("Error running Setup: %v", err)
  }
}

// FinishT is like Finish except that it calls t.Fatal on error.
func (r *Tester) FinishT(t *testing.T) {
  t.Helper()
  if err := r.Finish(); err != nil {
    t.Fatalf("Error running Finish: %v", err)
  }
}
