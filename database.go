package main

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"
)

type Record struct {
	App  string
	Tab  string
	When time.Time
}

type Records struct {
  records []Record
}

func (r *Records) String() string {
  var str strings.Builder
  for _, record := range r.records {
    str.WriteString(fmt.Sprintf("App: %s, Tab: %s, When: %v\n", record.App, record.Tab, record.When))
  }
  return str.String()
}

func SaveRecord(db *sql.DB, r *Record) error {
  stmt := "INSERT INTO switches (App, Tab, `When`) VALUES (?, ?, ?)"
  log.Println(stmt, r.App, r.Tab, r.When)

  _, err := db.Exec(stmt, r.App, r.Tab, r.When)
  if err != nil {
    fmt.Println("There was an error saving to the database:", err)
  }
  return err
}

func BulkSaveRecords(db *sql.DB, r *Records) error {
  valueStrings := make([]string, 0, len(r.records))
  valueArgs := make([]interface{}, 0, len(r.records) * 3)
  for _, record := range r.records {
    valueStrings = append(valueStrings, "(?, ?, ?)")
    valueArgs = append(valueArgs, record.App)
    valueArgs = append(valueArgs, record.Tab)
    valueArgs = append(valueArgs, record.When)
  }
  fmt.Println(r)

  stmt := fmt.Sprintf("INSERT INTO switches (App, Tab, `When`) VALUES %s", strings.Join(valueStrings, ","))
  log.Println(stmt)

  _, err := db.Exec(stmt, valueArgs...)
  if err != nil {
    fmt.Println("There was an error saving to the database brother:", err)
  }

  return err
}

