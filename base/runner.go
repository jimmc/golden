package base

import (
  "fmt"
  "io"
  "os"
  "path"
)

// GoldenRunner allows for configuring and running the different steps of the test.
type GoldenRunner struct {
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

func (r *GoldenRunner) SetupFilePath() string {
  return r.getFilePath(r.SetupPath, r.SetupBaseName, "setup")
}

func (r *GoldenRunner) OutFilePath() string {
  return r.getFilePath(r.OutPath, r.OutBaseName, "out")
}

func (r *GoldenRunner) GoldenFilePath() string {
  return r.getFilePath(r.GoldenPath, r.GoldenBaseName, "golden")
}

func (r *GoldenRunner) getFilePath(fpath, basename, extension string) string {
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

// Run runs a test using the configuration of the Runner.
func (r *GoldenRunner) Run() error {
  var ctx GoldenContext
  if r.CreateContext != nil {
    ctx, err := r.CreateContext()
    if err != nil {
      return err
    }
    defer ctx.Close()
  }

  outfilepath := r.OutFilePath()
  os.Remove(outfilepath)
  outfile, err := os.Create(outfilepath)
  if err != nil {
    return fmt.Errorf("error creating output file: %v", outfile)
  }

  // Run the specific test step.
  err = r.Test(ctx, outfile)

  if err != nil {
    return err
  }
  outfile.Close()
  goldenfilepath := r.GoldenFilePath()
  return CompareOutToGolden(outfilepath, goldenfilepath)
}
