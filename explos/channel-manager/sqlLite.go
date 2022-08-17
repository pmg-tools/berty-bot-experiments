package main

import (
	"errors"
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
	err = s.db.AutoMigrate(&User{}, &TeritoriKey{}, &Workspace{}, &Channel{})
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

type TeritoriKey struct {
	ID     uint `gorm:"primary_key"`
	PubKey string
	UserId uint
}

type User struct {
	ID              uint `gorm:"primary_key"`
	BertyPubKey     string
	TerritoriPubKey []TeritoriKey `gorm:"ForeignKey:UserId"`
}

func (s sqlLite) AddUser(bertyPubKey string) error {
	db := s.db
	tx := db.FirstOrCreate(&User{
		BertyPubKey: bertyPubKey,
	})
	if tx.Error != nil {
		return tx.Error
	}

	return nil
}

func (s sqlLite) SyncTeritoriKey(teritoriPubkey string, bertyPubkey string) error {
	db := s.db

	user := User{BertyPubKey: bertyPubkey}
	tx := db.FirstOrCreate(&user, "berty_pub_key = ?", bertyPubkey)
	db.Create(&TeritoriKey{PubKey: teritoriPubkey, UserId: user.ID})
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

func (s sqlLite) UserExist(bertyPubkey string) bool {
	var user User
	tx := s.db.Where("berty_pub_key = ?", bertyPubkey).First(&user)

	return !errors.Is(tx.Error, gorm.ErrRecordNotFound)
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
