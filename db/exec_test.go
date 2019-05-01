package db_test

import (
  "database/sql"
  "fmt"
  "reflect"
  "testing"

  goldendb "github.com/jimmc/golden/db"
)

type eTestRow struct {
  n int;
  s string;
}

func collectETestRows(db *sql.DB, query string) ([]*eTestRow, error) {
  rows := make([]*eTestRow, 0)
  row := &eTestRow{}
  targets := []interface{}{
    &row.n,
    &row.s,
  }
  collector := func() {
    rowCopy := eTestRow(*row)
    rows = append(rows, &rowCopy)
  }
  err := queryAndCollect(db, query, targets, collector)
  return rows, err
}

func setupAndCollectETestRows(setup, query string) ([]*eTestRow, error) {
  db, err := goldendb.EmptyDb()
  if err != nil {
    return nil, fmt.Errorf("error opening test database: %v", err)
  }
  defer db.Close()

  if err := goldendb.ExecMulti(db,setup); err != nil {
    return nil, fmt.Errorf("error calling ExecMulti: %v", err)
  }

  return collectETestRows(db, query)
}

// queryAndCollect issues a Query for the given sql, then interates through
// the returned rows. For each row, it retrieves the data into targets, then
// calls the collect function. The assumption is that the targets store the
// results into data that is accessible to the collect function.
func queryAndCollect(db *sql.DB, sql string, targets []interface{}, collect func()) error {
  rows, err := db.Query(sql)
  if err != nil {
    return err
  }
  defer rows.Close()
  for rows.Next() {
    err := rows.Scan(targets...)
    if err != nil {
      return err
    }
    collect()
  }
  return rows.Err()
}

func TestExecMulti(t *testing.T) {
  setup := `
CREATE table test(n int, s string);
INSERT into test(n, s) values(1, 'a'), (2, 'b'), (3, 'c');
`
  query := "SELECT n, s from test order by n;"
  expectedResult := []*eTestRow{
    &eTestRow{1, "a"},
    &eTestRow{2, "b"},
    &eTestRow{3, "c"},
  }

  rows, err := setupAndCollectETestRows(setup, query)
  if err != nil {
    t.Fatalf("Error collecting rows: %v", err)
  }

  if got, want := len(rows), 3; got != want {
    t.Fatalf("Wrong number of rows, got %d, want %d", got, want)
  }
  if got, want := rows, expectedResult; !reflect.DeepEqual(got, want) {
    t.Errorf("Results array, got %v, want %v", got, want)
  }
}

func TestComments(t *testing.T) {
  setup := `
CREATE table test(n int, s string);
# This is a comment
INSERT into test(n, s)
# another comment
  values(1, 'a');
`
  query := "SELECT n, s from test order by n;"
  expectedResult := []*eTestRow{
    &eTestRow{1, "a"},
  }

  rows, err := setupAndCollectETestRows(setup, query)
  if err != nil {
    t.Fatalf("Error collecting rows: %v", err)
  }

  if got, want := len(rows), 1; got != want {
    t.Fatalf("Wrong number of rows, got %d, want %d", got, want)
  }
  if got, want := rows, expectedResult; !reflect.DeepEqual(got, want) {
    t.Errorf("Results array, got %v, want %v", got, want)
  }
}

func TestExecErrors(t *testing.T) {
  setup := "invalid sql"
  db, err := goldendb.EmptyDb()
  if err != nil {
    t.Fatalf("error opening test database: %v", err)
  }
  defer db.Close()

  if err := goldendb.ExecMulti(db,setup); err == nil {
    t.Errorf("Expected error for invalid sql")
  }
}
