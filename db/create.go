// Package db contains test helper functions dealing with a database.
package db

import (
  "database/sql"
  "io/ioutil"

  _ "github.com/mattn/go-sqlite3"       // driver name: sqlite3
)

var (
  DbType = "sqlite3"
  DbName = ":memory:"
)

func EmptyDb() (*sql.DB, error) {
  return sql.Open(DbType, DbName)
}

func LoadSetupFile(db *sql.DB, filename string) error {
  setupSql, err := ioutil.ReadFile(filename)
  if err != nil {
    return err
  }
  return LoadSetupString(db, string(setupSql))
}

func LoadSetupString(db *sql.DB, setupSql string) error {
  return ExecMulti(db, setupSql)
}

func DbWithSetupFile(filename string) (*sql.DB, error) {
  setupSql, err := ioutil.ReadFile(filename)
  if err != nil {
    return nil, err
  }
  return DbWithSetupString(string(setupSql))
}

func DbWithSetupString(setupSql string) (*sql.DB, error) {
  db, err := EmptyDb()
  if err != nil {
    return nil, err
  }
  err = LoadSetupString(db, setupSql)
  if err != nil {
    db.Close()
    return nil, err
  }
  return db, nil
}
