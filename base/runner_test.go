package base_test

import (
  "testing"

  "github.com/jimmc/golden/base"
)

func TestBasePath(t *testing.T) {
  gr := &base.Runner{
    BaseName: "xyz",
  }
  if got, want := gr.OutFilePath(), "testdata/xyz.out"; got != want {
    t.Errorf("SetupFilePath with base name: got %q, want %q", got, want)
  }

  gr = &base.Runner{
    BaseDir: "bar",
  }
  if got, want := gr.OutFilePath(), "bar/test.out"; got != want {
    t.Errorf("SetupFilePath with base dir: got %q, want %q", got, want)
  }
}

func TestSetupFilePath(t *testing.T) {
  gr := &base.Runner{}
  if got, want := gr.SetupFilePath(), "testdata/test.setup"; got != want {
    t.Errorf("SetupFilePath on empty config: got %q, want %q", got, want)
  }

  gr = &base.Runner{
    SetupBaseName: "abc",
  }
  if got, want := gr.SetupFilePath(), "testdata/abc.setup"; got != want {
    t.Errorf("SetupFilePath with name: got %q, want %q", got, want)
  }

  gr = &base.Runner{
    SetupPath: "foo/abc.set-up",
  }
  if got, want := gr.SetupFilePath(), "foo/abc.set-up"; got != want {
    t.Errorf("SetupFilePath with path: got %q, want %q", got, want)
  }
}
