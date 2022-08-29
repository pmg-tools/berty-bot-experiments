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

func TestAddUser(t *testing.T) {
	var user User
	err := testQueries.AddUser("pubKeyTest")
	require.NoError(t, err)
	testQueries.db.First(&user)
	require.Equal(t, "pubKeyTest", user.BertyPubKey)
}

func TestAddChannel(t *testing.T) {
	var channel Channel
	err := testQueries.AddChannel("workspace_test", "channel_test", "link_test")
	require.NoError(t, err)
	testQueries.db.Preload("Workspace").First(&channel)
	require.Equal(t, "workspace_test", channel.Workspace.Name)
	require.Equal(t, "channel_test", channel.Name)
	require.Equal(t, "link_test", channel.BertyLink)
}
