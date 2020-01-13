package base

import (
  "bufio"
  "errors"
  "fmt"
  "os"
  "path"
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

// NewTester creates a new Tester instance.
func NewTester(basename string) *Tester {
  r := &Tester{
    BaseName: basename,
  }
  return r
}

// OutFilePath returns the complete path to the output file.
func (r *Tester) OutFilePath() string {
  return r.GetFilePath(r.OutPath, r.OutBaseName, "out")
}

// GoldenFilePath returns the complete path to the golden file.
func (r *Tester) GoldenFilePath() string {
  return r.GetFilePath(r.GoldenPath, r.GoldenBaseName, "golden")
}

// GetFilePath calculates and returns the complete path to a file.
// If fpath is set, it returns it, else it uses basename, with default "test",
// and the Tester's BaseDir, with default "testdata", plus the given extension
// to generate the path to return.
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

// Init is a no-op for this Tester.
func (r *Tester) Init() error {
  return nil
}

// Arrange creates the output files for the test to write to and sets
// OutF and OutW in the Tester.
func (r *Tester) Arrange() error {
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

// Act calls the users Test function.
func (r *Tester) Act() error {
  if r.Test == nil {
    return errors.New("no Test method set")
  }
  return r.Test(r)
}

// Assert closes the output and compares it to the golden file.
func (r *Tester) Assert() error {
  r.OutW.Flush()
  r.OutF.Close()
  return CompareOutToGolden(r.OutFilePath(), r.GoldenFilePath())
}

// Close is a no-op in this Tester.
func (r *Tester) Close() error {
  return nil
}
