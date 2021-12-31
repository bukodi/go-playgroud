package main

import (
	"github.com/codenotary/immudb/embedded/sql"
	"github.com/codenotary/immudb/embedded/store"
	"log"
	"os"
)

func main() {
	os.RemoveAll("immudbtest/sqldata")
	dataStore, err := store.Open("immudbtest/sqldata", store.DefaultOptions())
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll("immudbtest/sqldata")

	// And now you can create the SQL engine, passing both stores and a key prefix:
	engine, err := sql.NewEngine(dataStore, sql.DefaultOptions())
	if err != nil {
		log.Fatal(err)
	}

	// Create and use a database:
	_, _, err = engine.Exec("CREATE DATABASE db1", nil, nil)
	if err != nil {
		log.Fatal(err)
	}

	_, _, err = engine.Exec("USE DATABASE db1", nil, nil)
	if err != nil {
		log.Fatal(err)
	}

	// The engine has an API to execute statements and queries. To execute an statement:
	_, _, err = engine.Exec("USE DATABASE db1; CREATE TABLE journal (id INTEGER, date VARCHAR, creditaccount INTEGER, debitaccount INTEGER, amount INTEGER, description VARCHAR, PRIMARY KEY id)", nil, nil)
	if err != nil {
		log.Fatal(err)
	}

	err = engine.SetDefaultDatabase("db1")

	// Queries can be executed using QueryStmt and you can pass a map of parameters to substitute, and whether the engine should wait for indexing:
	r, err := engine.Query("SELECT id, date, creditaccount, debitaccount, amount, description FROM db1.journal WHERE amount > @value", map[string]interface{}{"value": 100}, nil)

	// To iterate over a result set r, just fetch rows until there are no more entries. Every row has a Values member you can index to access the column:
	for {
		row, err := r.Read()
		if err == sql.ErrNoMoreRows {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("row: %v", row.Values)
	}

}
