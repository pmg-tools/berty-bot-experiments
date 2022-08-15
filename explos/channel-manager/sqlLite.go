package main

import (
	"gorm.io/gorm"
)

type Workspace struct {
	ID       uint `gorm:"primary_key"`
	Name     string
	Channels []Channel `gorm:"ForeignKey:Wid"`
}

type Channel struct {
	ID          uint `gorm:"primary_key"`
	BertyLink   string
	ChannelName string
	Wid         uint
}

type User struct {
	ID              uint `gorm:"primary_key"`
	TerritoriPubKey string
	BertyPubKey     string
	Status          int
	Nonce           int
}

type sqlLite struct {
	db *gorm.DB
}

func (s sqlLite) AddUser(territoriPubKey string, bertyPubKey string, nonce int) (ok bool) {
	// gest user exist cases (berty and territori pubKeys)
	db := s.db
	db.Create(&User{
		TerritoriPubKey: territoriPubKey,
		BertyPubKey:     bertyPubKey,
		Status:          0,
		Nonce:           nonce,
	})

	return true
}

func (s sqlLite) AddChannel(workspaceName string, channelName string, bertyGroupLink string) (ok bool) {
	db := s.db

	var workspace Workspace
	_ = db.Where("name = ?", workspaceName).First(&workspace)
	if workspace.Name == "" {
		ws := &Workspace{Name: workspaceName}
		db.Create(ws)
		workspace.ID = ws.ID
	}

	_ = db.Create(&Channel{
		ChannelName: channelName,
		BertyLink:   bertyGroupLink,
		Wid:         workspace.ID,
	})

	return true
}

func (s sqlLite) AddWorkspace(workspaceName string) (ok bool) {
	db := s.db
	_ = db.Create(&Workspace{Name: workspaceName})

	return true
}

func (s sqlLite) UserExist(pubKey string) bool {
	var user User
	_ = s.db.Where("pub_key = ?", pubKey).First(&user)

	return user.TerritoriPubKey != ""
}

func (s sqlLite) ConfirmUser(territoriPubKey string, bertyPubKey string) (ok bool) {
	var user User
	if err := s.db.Where("territori_pub_key = ? and berty_pub_key = ?", territoriPubKey, bertyPubKey).First(&user); err.Error != nil {
		return false
	}

	user.Status = 1
	user.Nonce = 0

	if err := s.db.Save(&user); err.Error != nil {
		return false
	}

	return true
}

func (s sqlLite) ChannelExist(workspaceName string, channelName string) bool {
	var workspace Workspace
	_ = s.db.Where("name = ?", workspaceName).First(&workspace)
	if workspace.Name == "" {
		return false
	}

	var channel Channel
	_ = s.db.Where("wid = ? AND channel_name = ?", workspace.ID, channelName).First(&channel)

	return channel.ChannelName != ""
}

func (s sqlLite) WorkspaceExist(workspaceName string) bool {
	var workspace Workspace
	_ = s.db.Where("name = ?", workspaceName).First(&workspace)

	return workspace.Name != ""
}

func (s sqlLite) ListWorkspaces() ([]string, error) {
	var workspaces []Workspace
	_ = s.db.Find(&workspaces)

	var workspaceIDs []string
	for _, workspace := range workspaces {
		workspaceIDs = append(workspaceIDs, workspace.Name)
	}

	return workspaceIDs, nil
}

func (s sqlLite) ListChannels(workspaceName string) ([]string, error) {
	var workspace Workspace
	_ = s.db.Where("name = ?", workspaceName).First(&workspace)

	var channels []Channel
	_ = s.db.Where("wid = ?", workspace.ID).Find(&channels)

	var channelIDs []string
	for _, channel := range channels {
		channelIDs = append(channelIDs, channel.ChannelName)
	}

	return channelIDs, nil
}

func (s sqlLite) GetChannelsInvitation(workspaceName string, channelsName []string) []Channel {

	var workspace Workspace
	_ = s.db.Where("name = ?", workspaceName).First(&workspace)

	var channels []Channel
	_ = s.db.Where("wid = ? AND channel_name IN(?)", workspace.ID, channelsName).Find(&channels)
	return channels
}
