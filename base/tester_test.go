package base_test

import (
  "errors"
  "fmt"
  "io"
  "testing"

  "github.com/jimmc/golden/base"
)

// Example is used as the function to be tested by our testing code.
func example(s string) string {
  return fmt.Sprintf("This is the output of example(%q).\n", s)
}

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

func TestNoTestSet(t *testing.T) {
  r := base.NewTester("example")
  if err := r.Init(); err != nil {
    t.Fatalf("Error in Init: %v", err)
  }
  if err := r.Arrange(); err != nil {
    t.Fatalf("Error in Arrange: %v", err)
  }
  if err := r.Act(); err == nil {
    t.Fatalf("Expected error due to no test being set")
  }
}

func TestActError(t *testing.T) {
  r := base.NewTester("example")
  r.Test = func(r *base.Tester) error {
    return errors.New("intentional error from function under test")
  }
  if err := r.Init(); err != nil {
    t.Fatalf("Error in Init: %v", err)
  }
  if err := r.Arrange(); err != nil {
    t.Fatalf("Error in Arrange: %v", err)
  }
  if err := r.Act(); err == nil {
    t.Fatalf("Expected error from function under test")
  }
}

func TestNoGoldenFile(t *testing.T) {
  r := base.NewTester("example-no-golden")
  r.Test = func(r *base.Tester) error {
    s := example("no-golden")
    _, err := io.WriteString(r.OutW, s)
    return err
  }
  if err := r.Init(); err != nil {
    t.Fatalf("Error in Init: %v", err)
  }
  if err := r.Arrange(); err != nil {
    t.Fatalf("Error in Arrange: %v", err)
  }
  if err := r.Act(); err != nil {
    t.Fatalf("Error in Act: %v", err)
  }
  if err := r.Assert(); err == nil {
    t.Fatalf("Expected error in Assert due to no golden file")
  }
}

func TestLifecycle(t *testing.T) {
  r := base.NewTester("example")
  r.Test = func(r *base.Tester) error {
    s := example("happy-path")
    _, err := io.WriteString(r.OutW, s)
    return err
  }
  if err := r.Init(); err != nil {
    t.Fatalf("Error in Init: %v", err)
  }
  if err := r.Arrange(); err != nil {
    t.Fatalf("Error in Arrange: %v", err)
  }
  if err := r.Act(); err != nil {
    t.Fatalf("Error in Act: %v", err)
  }
  if err := r.Assert(); err != nil {
    t.Fatalf("Error in Assert: %v", err)
  }
  if err := r.Close(); err != nil {
    t.Fatalf("Error in Close: %v", err)
  }
}

func TestTwoTests(t *testing.T) {
  r := base.NewTester("example1")
  if err := r.Init(); err != nil {
    t.Fatalf("Error in Init: %v", err)
  }

  r.Test = func(r *base.Tester) error {
    s := example("test1")
    _, err := io.WriteString(r.OutW, s)
    return err
  }
  if err := r.Arrange(); err != nil {
    t.Fatalf("Error in Arrange: %v", err)
  }
  if err := r.Act(); err != nil {
    t.Fatalf("Error in Act: %v", err)
  }
  if err := r.Assert(); err != nil {
    t.Fatalf("Error in Assert: %v", err)
  }

  r.BaseName = "example2"
  r.Test = func(r *base.Tester) error {
    s := example("test2")
    _, err := io.WriteString(r.OutW, s)
    return err
  }
  if err := r.Arrange(); err != nil {
    t.Fatalf("Error in Arrange: %v", err)
  }
  if err := r.Act(); err != nil {
    t.Fatalf("Error in Act: %v", err)
  }
  if err := r.Assert(); err != nil {
    t.Fatalf("Error in Assert: %v", err)
  }

  if err := r.Close(); err != nil {
    t.Fatalf("Error in Close: %v", err)
  }
}
