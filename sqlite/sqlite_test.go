package main

import (
	"database/sql"
	"fmt"
	"strconv"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func TestRAWSqlite(t *testing.T) {
	fmt.Println("Start")
	db, _ := sql.Open("sqlite3", ":memory:")
	//db, _ := sql.Open("sqlite3", "/tmp/nraboy.db")
	stmt, _ := db.Prepare("CREATE TABLE IF NOT EXISTS people (id INTEGER PRIMARY KEY, firstname TEXT, lastname TEXT)")
	stmt.Exec()
	stmt, _ = db.Prepare("INSERT INTO people (firstname, lastname) VALUES (?, ?)")
	stmt.Exec("Nic", "Raboy")
	rows, _ := db.Query("SELECT id, firstname, lastname FROM people")
	var id int
	var firstname string
	var lastname string
	for rows.Next() {
		rows.Scan(&id, &firstname, &lastname)
		fmt.Println(strconv.Itoa(id) + ": " + firstname + " " + lastname)
	}
}
