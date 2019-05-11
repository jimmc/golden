package base_test

import (
  "io"
  "testing"

  "github.com/jimmc/golden/base"
)

func TestRun(t *testing.T) {
  r := base.NewTester("run-example")
  r.Test = func(r *base.Tester) error {
    s := example("run")
    _, err := io.WriteString(r.OutW, s)
    return err
  }
  if err := base.Run(r); err != nil {
    t.Fatalf("Error in Run: %v", err)
  }
}
