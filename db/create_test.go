package db_test

import (
  "testing"

  "github.com/google/go-cmp/cmp"

  goldendb "github.com/jimmc/golden/db"
)

func TestDbWithFile(t *testing.T) {
  db, err := goldendb.DbWithSetupFile("testdata/db1.txt")
  if err != nil {
    t.Fatalf("DbWithSetupFile unexpected error: %v", err)
  }

  query := "SELECT n, s from test order by n;"
  expectedResult := []*eTestRow{
    &eTestRow{1, "a"},
    &eTestRow{2, "b"},
    &eTestRow{3, "c"},
  }

  rows, err := collectETestRows(db, query)
  if err != nil {
    t.Fatalf("Error collecting rows: %v", err)
  }

  if got, want := len(rows), 3; got != want {
    t.Fatalf("Wrong number of rows, got %d, want %d", got, want)
  }
  got, want := rows, expectedResult
  if diff := cmp.Diff(want, got, cmp.AllowUnexported(eTestRow{})); diff != "" {
    t.Errorf("Results mismatch (-want +got):\n%s", diff)
  }
}
