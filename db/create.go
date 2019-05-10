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

// EmptyDb creates an empty database from the values in our variables DbType and DbName.
func EmptyDb() (*sql.DB, error) {
  return sql.Open(DbType, DbName)
}

// LoadSetupFile reads and executes SQL commands from the specified file.
func LoadSetupFile(db *sql.DB, filename string) error {
  setupSql, err := ioutil.ReadFile(filename)
  if err != nil {
    return err
  }
  return LoadSetupString(db, string(setupSql))
}

// LoadSetupString reads and executes SQL commands from the given string.
func LoadSetupString(db *sql.DB, setupSql string) error {
  return ExecMulti(db, setupSql)
}

// DbWithSetupFile creates a new database and executes SQL commands from the given file.
func DbWithSetupFile(filename string) (*sql.DB, error) {
  setupSql, err := ioutil.ReadFile(filename)
  if err != nil {
    return nil, err
  }
  return DbWithSetupString(string(setupSql))
}

// DbWithSetupFile creates a new database and executes SQL commands from the given string.
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
