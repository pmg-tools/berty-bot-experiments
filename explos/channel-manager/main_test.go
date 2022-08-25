package main

import (
	"log"
	"os"
	"testing"
)

var testQueries *sqlLite

func TestMain(m *testing.M) {
	var err error
	testQueries, err = NewSqliteDB(":memory:")
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}
	os.Exit(m.Run())
}
