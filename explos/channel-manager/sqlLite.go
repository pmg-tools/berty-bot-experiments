package main

import (
	"errors"
	sqlite "github.com/flyingtime/gorm-sqlcipher"
	"gorm.io/gorm"
)

type sqlLite struct {
	db *gorm.DB
}

func NewSqliteDB() (*sqlLite, error) {
	db, err := gorm.Open(sqlite.Open("test.db"))
	if err != nil {
		return nil, err
	}

	s := &sqlLite{db: db}
	err = s.db.AutoMigrate(&User{}, &Workspace{}, &Channel{})
	if err != nil {
		return nil, err
	}

	return s, nil
}

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

type teritoriKey struct {
	PubKey string `gorm:"primary_key"`
	UserId uint
}

type User struct {
	ID             uint `gorm:"primary_key"`
	BertyPubKey    string
	teritoriPubKey []teritoriKey `gorm:"ForeignKey:UserId"`
}

func (s sqlLite) AddUser(teritoriPubKey string, bertyPubKey string, nonce int) error {
	// gest user exist cases (berty and teritori pubKeys)
	db := s.db
	tx := db.Create(&User{
		BertyPubKey: bertyPubKey,
	})
	if tx.Error != nil {
		return tx.Error
	}

	return nil
}

func (s sqlLite) AddChannel(workspaceName string, channelName string, bertyGroupLink string) error {
	db := s.db

	var workspace Workspace
	tx := db.Where("name = ?", workspaceName).First(&workspace)
	if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		ws := &Workspace{Name: workspaceName}
		db.Create(ws)
		workspace.ID = ws.ID
	}

	tx = db.Create(&Channel{
		ChannelName: channelName,
		BertyLink:   bertyGroupLink,
		Wid:         workspace.ID,
	})
	if tx.Error != nil {
		return tx.Error
	}

	return nil
}

func (s sqlLite) AddWorkspace(workspaceName string) error {
	db := s.db
	tx := db.Create(&Workspace{Name: workspaceName})
	if tx.Error != nil {
		return tx.Error
	}

	return nil
}

func (s sqlLite) UserExist(pubKey string) bool {
	var user User
	tx := s.db.Where("pub_key = ?", pubKey).First(&user)

	return !errors.Is(tx.Error, gorm.ErrRecordNotFound)
}

func (s sqlLite) ConfirmUser(teritoriPubKey string, bertyPubKey string) (ok bool) {
	var user User
	if err := s.db.Where("teritori_pub_key = ? and berty_pub_key = ?", teritoriPubKey, bertyPubKey).First(&user); err.Error != nil {
		return false
	}

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
