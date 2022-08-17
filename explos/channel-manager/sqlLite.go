package main

import (
	"fmt"
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

	var user User
	tx := db.Where(&User{
		BertyPubKey: bertyPubKey,
	}).FirstOrInit(&user)
	if tx.Error != nil {
		return tx.Error
	}

	return nil
}

func (s sqlLite) AddChannel(workspaceName string, channelName string, bertyGroupLink string) error {
	db := s.db

	var workspace Workspace
	tx := db.Where(Workspace{Name: workspaceName}).FirstOrCreate(&workspace)

	var channel Channel
	tx = db.Where(&Channel{
		ChannelName: channelName,
		BertyLink:   bertyGroupLink,
		Wid:         workspace.ID,
	}).FirstOrCreate(&channel)
	if tx.Error != nil {
		return tx.Error
	}

	return nil
}

func (s sqlLite) AddWorkspace(workspaceName string) error {
	db := s.db
	var workspace Workspace
	tx := db.Where(Workspace{Name: workspaceName}).FirstOrCreate(&workspace)
	if tx.Error != nil {
		return tx.Error
	}

	return nil
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

	newChannel := false
	createChannel := true
	for _, v := range channelsName {
		for _, w := range channels {
			if v == w.ChannelName {
				createChannel = false
				break
			}
		}
		if createChannel == true {
			link, err := bertyBotCreateGroup(fmt.Sprintf("%s/#%s", workspaceName, v))
			if err != nil {
				return nil
			}

			err = s.AddChannel(workspaceName, v, link)
			newChannel = true
		}
		createChannel = true
	}

	if newChannel == true {
		_ = s.db.Where("name = ?", workspaceName).First(&workspace)
		_ = s.db.Where("wid = ? AND channel_name IN(?)", workspace.ID, channelsName).Find(&channels)
	}

	return channels
}
