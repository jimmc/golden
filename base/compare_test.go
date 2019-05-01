package base_test

import (
  "testing"

  "github.com/jimmc/golden/base"
)

func TestCompareGood(t *testing.T) {
  err := base.CompareOutToGolden("testdata/a.txt", "testdata/a.txt")
  if err != nil {
    t.Errorf("CompareOutToGolden: should be same but got error: %v", err)
  }
}

func TestCompareBad(t *testing.T) {
  err := base.CompareOutToGolden("testdata/a.txt", "testdata/b.txt")
  if err == nil {
    t.Fatal("CompareOutToGolden: expected error about different contents")
  }
}

func TestCompareNoOut(t *testing.T) {
  err := base.CompareOutToGolden("no-such-file", "testdata/a.txt")
  if err == nil {
    t.Fatal("CompareOutToGolden: expected error about no output file")
  }
}

func TestCompareNoGolden(t *testing.T) {
  err := base.CompareOutToGolden("testdata/a.txt", "no-such-file")
  if err == nil {
    t.Fatal("CompareOutToGolden: expected error about no golden file")
  }
}
