package base_test

import (
  "testing"

  "github.com/jimmc/golden/base"
)

func TestBasePath(t *testing.T) {
  gr := &base.Tester{
    BaseName: "xyz",
  }
  if got, want := gr.OutFilePath(), "testdata/xyz.out"; got != want {
    t.Errorf("OutFilePath with base name: got %q, want %q", got, want)
  }

  gr = &base.Tester{
    BaseDir: "bar",
  }
  if got, want := gr.OutFilePath(), "bar/test.out"; got != want {
    t.Errorf("OutFilePath with base dir: got %q, want %q", got, want)
  }
}

func TestOutFilePath(t *testing.T) {
  gr := &base.Tester{}
  if got, want := gr.OutFilePath(), "testdata/test.out"; got != want {
    t.Errorf("OutFilePath on empty config: got %q, want %q", got, want)
  }

  gr = &base.Tester{
    OutBaseName: "abc",
  }
  if got, want := gr.OutFilePath(), "testdata/abc.out"; got != want {
    t.Errorf("OutFilePath with name: got %q, want %q", got, want)
  }

  gr = &base.Tester{
    OutPath: "foo/abc.oot",
  }
  if got, want := gr.OutFilePath(), "foo/abc.oot"; got != want {
    t.Errorf("OutFilePath with path: got %q, want %q", got, want)
  }
}
