package main

import (
	"errors"
)

type channel struct {
	ID        string
	BertyLink string
}

type Channels map[string]*channel

type mockDb struct {
	Workspaces map[string]Channels
	// TODO: add a way to link channels to discord
	// Channels map[channel][]message
}

func (m *mockDb) AddWorkspace(workspaceID string) {
	if m.Workspaces == nil {
		m.Workspaces = make(map[string]Channels)
	}
	if m.Workspaces[workspaceID] == nil {
		m.Workspaces[workspaceID] = make(map[string]*channel)
	}
}

func (m *mockDb) AddChannel(workspaceID string, channelID string) {
	if m.Workspaces == nil || m.Workspaces[workspaceID] == nil {
		m.AddWorkspace(workspaceID)
	}
	if m.Workspaces[workspaceID][channelID] == nil {
		m.Workspaces[workspaceID][channelID] = &channel{
			ID:        channelID,
			BertyLink: "", // TODO: create group chat link
		}
	}
}

func (m *mockDb) WorkspaceExist(workspaceID string) bool {
	return m.Workspaces[workspaceID] != nil
}

func (m *mockDb) ChannelExist(workspaceID string, channelID string) bool {
	if m.Workspaces == nil {
		return false
	}
	return m.Workspaces[workspaceID][channelID] != nil
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
	for c := range m.Workspaces[workspaceID] {
		channels = append(channels, c)
	}
	return channels, nil
}
