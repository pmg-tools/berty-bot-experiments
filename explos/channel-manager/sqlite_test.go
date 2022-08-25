package main

import (
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestNewSqliteDB(t *testing.T) {
	db, err := NewSqliteDB("tmp.db")
	require.NoError(t, err)
	require.NotNil(t, db)
	_, err = os.Stat("tmp.db")
	require.NoError(t, err)
	err = os.Remove("tmp.db")
	if err != nil {
		t.Error("Can't remove tmp.db")
	}
}
