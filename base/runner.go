package base

import (
  "bufio"
  "fmt"
  "io"
  "os"
  "path"
  "testing"
)

// Runner allows for configuring and running the different steps of the test.
type Runner struct {
  BaseName string       // Base name for all files in the test, if not overridden; defaults to "test".
  BaseDir string        // Base directory; if not set, uses "testdata".

  SetupBaseName string  // Base name for the setup file; if not set, uses BaseName.
  SetupPath string      // Path to the setup file; if not set, uses SetupBaseName.

  OutBaseName string    // Base name for the test output file; if not set, uses BaseName.
  OutPath string        // Path to the test output file; if not set, uses OutBaseName.

  GoldenBaseName string // Base name for the golden file; if not set, uses BaseName.
  GoldenPath string     // Path to the golden file; if not set, uses GoldenBaseName.

  CreateContext func() (GoldenContext, error)    // Function to create the context for testing.
  Test func(GoldenContext, *os.File) error     // Function to run the test.
}

type GoldenContext interface {
  io.Closer
}

func (r *Runner) SetupFilePath() string {
  return r.getFilePath(r.SetupPath, r.SetupBaseName, "setup")
}

func (r *Runner) OutFilePath() string {
  return r.getFilePath(r.OutPath, r.OutBaseName, "out")
}

func (r *Runner) GoldenFilePath() string {
  return r.getFilePath(r.GoldenPath, r.GoldenBaseName, "golden")
}

func (r *Runner) getFilePath(fpath, basename, extension string) string {
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

type RunData struct {
  OutW *bufio.Writer;
  OutF *os.File;
  Ctx GoldenContext
}

func (r *Runner) Setup() (*RunData, error) {
  var ctx GoldenContext
  if r.CreateContext != nil {
    var err error
    ctx, err = r.CreateContext()
    if err != nil {
      return nil, err
    }
  }

  outfilepath := r.OutFilePath()
  os.Remove(outfilepath)
  f, err := os.Create(outfilepath)
  if err != nil {
    if ctx != nil {
      ctx.Close()
    }
    return nil, fmt.Errorf("error creating output file %q: %v", outfilepath, err)
  }
  w := bufio.NewWriter(f)

  return &RunData{
    OutF: f,
    OutW: w,
    Ctx: ctx,
  }, nil
}

func (r *Runner) Finish(data *RunData) error {
  data.OutW.Flush()
  data.OutF.Close()
  if data.Ctx != nil {
    data.Ctx.Close()
  }
  return CompareOutToGolden(r.OutFilePath(), r.GoldenFilePath())
}

// Run runs a test using the configuration of the Runner.
func (r *Runner) Run() error {
  runData, err := r.Setup()
  if err != nil {
    return err
  }
  if runData.Ctx != nil {
    defer runData.Ctx.Close()
  }

  // Run the specific test step.
  err = r.Test(runData.Ctx, runData.OutF)
  if err != nil {
    return err
  }

  // Check the output to see if we got the right data.
  return r.Finish(runData)
}

// SetupT is like Setup except that it calls t.Fatal on error.
func (r *Runner) SetupT(t *testing.T) *RunData {
  t.Helper()
  d, err := r.Setup()
  if err != nil {
    t.Fatalf("Error running Setup: %v", err)
  }
  return d
}

// FinishT is like Finish except that it calls t.Fatal on error.
func (r *Runner) FinishT(t *testing.T, data *RunData) {
  t.Helper()
  err := r.Finish(data)
  if err != nil {
    t.Fatalf("Error running Finish: %v", err)
  }
}

// RunT is like Run except that it calls t.Fatal on error.
func (r *Runner) RunT(t *testing.T) {
  t.Helper()
  err := r.Run()
  if err != nil {
    t.Fatalf("Error running Run: %v", err)
  }
}
