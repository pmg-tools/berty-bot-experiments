package main

import (
	"fmt"

	sqlite "github.com/flyingtime/gorm-sqlcipher"
	"gorm.io/gorm"
)

type sqlLite struct {
	db *gorm.DB
}

func NewSqliteDB(databaseName string) (*sqlLite, error) {
	db, err := gorm.Open(sqlite.Open(databaseName))
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
	Channels []Channel
}

type Channel struct {
	ID          uint `gorm:"primary_key"`
	BertyLink   string
	Name        string
	WorkspaceID uint
	Workspace   *Workspace
}

type TeritoriKey struct {
	ID     uint `gorm:"primary_key"`
	PubKey string
	UserId uint
}

type User struct {
	ID             uint `gorm:"primary_key"`
	BertyPubKey    string
	TeritoriPubKey []TeritoriKey `gorm:"ForeignKey:UserId"`
}

// AddUser to database if it doesn't exist
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
	if tx.Error != nil {
		return tx.Error
	}

	tx = db.Create(&TeritoriKey{PubKey: teritoriPubkey, UserId: user.ID})
	if tx.Error != nil {
		return tx.Error
	}

	return nil
}

// AddChannel to database if it doesn't exist
func (s sqlLite) AddChannel(workspaceName string, channelName string, bertyGroupLink string) error {
	db := s.db

	var workspace Workspace
	tx := db.Where(Workspace{Name: workspaceName}).FirstOrCreate(&workspace)

	var channel Channel
	tx = db.Where(&Channel{
		Name:        channelName,
		BertyLink:   bertyGroupLink,
		WorkspaceID: workspace.ID,
	}).FirstOrCreate(&channel)
	if tx.Error != nil {
		return tx.Error
	}

	return nil
}

// AddWorkspace to database if it doesn't exist
func (s sqlLite) AddWorkspace(workspaceName string) error {
	db := s.db
	var workspace Workspace
	tx := db.Where(Workspace{Name: workspaceName}).FirstOrCreate(&workspace)
	if tx.Error != nil {
		return tx.Error
	}

	return nil
}

// ListWorkspaces return all channels
func (s sqlLite) ListWorkspaces() ([]string, error) {
	var workspaces []Workspace
	_ = s.db.Find(&workspaces)

	var workspaceIDs []string
	for _, workspace := range workspaces {
		workspaceIDs = append(workspaceIDs, workspace.Name)
	}

	return workspaceIDs, nil
}

// ListUsers return all users
func (s sqlLite) ListUsers() ([]User, error) {
	var users []User
	_ = s.db.Find(&users)

	return users, nil
}

// ListChannels return all channels of a workspace
func (s sqlLite) ListChannels(workspaceName string) ([]string, error) {
	var workspace Workspace
	s.db.Where("name = ?", workspaceName).Preload("Channels").First(&workspace)

	var channelIDs []string
	for _, channel := range workspace.Channels {
		channelIDs = append(channelIDs, channel.Name)
	}

	return channelIDs, nil
}

// GetChannelsInvitation get a list of channels in a workspace and return their invitation links
func (s sqlLite) GetChannelsInvitation(workspaceName string, channelsName []string) []Channel {
	var workspace Workspace
	s.db.Where("name = ?", workspaceName).Preload("Channels", "name in (?)", channelsName).First(&workspace)

	newChannel := false
	createChannel := true
	for _, v := range channelsName {
		for _, w := range workspace.Channels {
			if v == w.Name {
				createChannel = false
				break
			}
		}
		if createChannel == true {
			chanName := fmt.Sprintf("%s/#%s", workspaceName, v)
			link, err := bertyBotCreateGroup(chanName)
			if err != nil {
				return nil
			}

			err = s.AddChannel(workspaceName, v, link)
			newChannel = true
		}
		createChannel = true
	}

	if newChannel == true {
		s.db.Where("name = ?", workspaceName).Preload("Channels", "name in (?)", channelsName).First(&workspace)
	}

	return workspace.Channels
}
