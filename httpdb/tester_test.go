package httpdb_test

import (
  "database/sql"
  "fmt"
  "net/http"
  "testing"

  goldenhttpdb "github.com/jimmc/golden/httpdb"
)

type dbhandler struct {
  db *sql.DB
}

type exampleRow struct {
  s string
  n int
}

func readDbRows(db *sql.DB) ([]exampleRow, error) {
  sql := "SELECT s, n FROM test ORDER BY s;"
  rows, err := db.Query(sql)
  if err != nil {
    return nil, err
  }
  defer rows.Close()
  xRows := []exampleRow{}
  for rows.Next() {
    var s string
    var n int
    err := rows.Scan(&s, &n)
    if err != nil {
      return nil, err
    }
    xRow := exampleRow{s: s, n: n}
    xRows = append(xRows, xRow)
  }
  if err := rows.Err(); err != nil {
    return nil, err
  }
  return xRows, nil
}

func (h *dbhandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
  xRows, err := readDbRows(h.db)
  if err != nil {
    http.Error(w, fmt.Sprintf("Error reading database: %v", err), http.StatusInternalServerError)
    return
  }
  w.WriteHeader(http.StatusOK)
  for _, xRow := range xRows {
    s := fmt.Sprintf("%v\n", xRow)
    w.Write([]byte(s))
  }
}

func TestHttpDbTester(t *testing.T) {
  request := func() (*http.Request, error) {
    return http.NewRequest("GET", "/api/foo/", nil)
  }
  r := goldenhttpdb.NewTester(func (r *goldenhttpdb.Tester) http.Handler {
    h := &dbhandler{}
    h.db = r.DB
    return h
  })
  if err := goldenhttpdb.RunOneWith(r, "foo-db", request); err != nil {
    t.Fatalf("Error in Run: %s", err)
  }
}
