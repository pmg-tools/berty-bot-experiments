package main

import (
	"errors"
	"time"
)

type message struct {
	Text      string
	ChannelID string
	hour      time.Time
}

/*type channel struct {
	ID   string
	discordLink string
}*/

type mockDb struct {
	Channels map[string][]message
	// TODO: add a way to link channels to discord
	// Channels map[channel][]message
}

func (m *mockDb) AddChannel(channelID string) error {
	if m.Channels == nil {
		m.Channels = make(map[string][]message)
	}
	if m.Channels[channelID] == nil {
		m.Channels[channelID] = make([]message, 0)
		return nil
	}
	return errors.New("channel already exists")
}

func (m *mockDb) ListChannels() ([]string, error) {
	if m.Channels == nil {
		return nil, errors.New("no channels")
	}
	var channels []string
	for channelID := range m.Channels {
		channels = append(channels, channelID)
	}
	return channels, nil
}

func (m *mockDb) ChannelExist(channelID string) bool {
	if m.Channels == nil {
		return false
	}
	return m.Channels[channelID] != nil
}

func (m *mockDb) AddMessage(channelID string, msg message) {
	m.Channels[channelID] = append(m.Channels[channelID], msg)
}

func (m *mockDb) GetChannelMessages(channelID string) []message {
	return m.Channels[channelID]
}
