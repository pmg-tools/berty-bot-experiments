package main

import (
	"errors"
	"fmt"
)

/*type channel struct {
	ID   string
	discordLink string
}*/

type Channel map[string]string

type mockDb struct {
	Workspaces map[string]Channel
	// TODO: add a way to link channels to discord
	// Channels map[channel][]message
}

func (m *mockDb) AddWorkspace(workspaceID string) {
	if m.Workspaces == nil {
		m.Workspaces = make(map[string]Channel)
	}
	if m.Workspaces[workspaceID] == nil {
		m.Workspaces[workspaceID] = make(map[string]string)
	}
}

func (m *mockDb) AddChannel(workspaceID string, channelID string) {
	if m.Workspaces == nil || m.Workspaces[workspaceID] == nil {
		m.AddWorkspace(workspaceID)
	}
	if m.Workspaces[workspaceID][channelID] == "" {
		m.Workspaces[workspaceID][channelID] = channelID
	}
}

func (m *mockDb) WorkspaceExist(workspaceID string) bool {
	return m.Workspaces[workspaceID] != nil
}

func (m *mockDb) ChannelExist(workspaceID string, channelID string) bool {
	if m.Workspaces == nil {
		return false
	}
	return m.Workspaces[workspaceID][channelID] != ""
}

func (m *mockDb) ListWorkspaces() ([]string, error) {
	if m.Workspaces == nil {
		return nil, errors.New("no workspaces")
	}
	var workspaces []string
	for workspaceID := range m.Workspaces {
		workspaces = append(workspaces, workspaceID)
	}
	return workspaces, nil
}

func (m *mockDb) ListChannels(workspaceID string) ([]string, error) {
	if m.Workspaces == nil {
		return nil, errors.New("no workspaces")
	}
	if m.Workspaces[workspaceID] == nil {
		return nil, errors.New("no channels")
	}
	var channels []string

	// debug
	fmt.Println(m.Workspaces[workspaceID])

	for channelID := range m.Workspaces[workspaceID] {
		channels = append(channels, channelID)
	}
	return channels, nil
}
