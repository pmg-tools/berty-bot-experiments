package main

import (
	"berty.tech/berty/v2/go/pkg/bertybot"
	"errors"
	"fmt"
	"strings"
	"sync"
)

type database interface {
	UserExist(pubKey string) bool
	AddUser(territoriPubKey string, bertyPubKey string, nonce int) error
	ConfirmUser(territoriPubKey string, bertyPubKey string) bool

	AddWorkspace(workspaceName string) error
	AddChannel(workspaceName string, channelName string, bertyGroupLink string) error

	WorkspaceExist(workspaceName string) bool
	ChannelExist(workspaceName string, channelName string) bool

	GetChannelsInvitation(workspaceName string, channelsName []string) []Channel

	// debug functions
	ListWorkspaces() ([]string, error)
	ListChannels(workspaceName string) ([]string, error)
}

func bertyBotAddWorkspace(db database, mutex *sync.Mutex) func(ctx bertybot.Context) {
	return func(ctx bertybot.Context) {
		// TODO: upgrade sdk to avoid it
		if ctx.IsReplay || !ctx.IsNew {
			return
		}

		workspaceName, err := func(ctx bertybot.Context) (string, error) {
			if len(ctx.CommandArgs) == 2 {
				return ctx.CommandArgs[1], nil
			}
			return "", errors.New("bad arguments")
		}(ctx)
		if err != nil {
			_ = ctx.ReplyString(err.Error())
			return
		}
		//

		mutex.Lock()
		if db.WorkspaceExist(workspaceName) {
			_ = ctx.ReplyString("workspace already exists")
			return
		}

		err = db.AddWorkspace(workspaceName)
		mutex.Unlock()

		if err != nil {
			_ = ctx.ReplyString(err.Error())
			return
		}

		_ = ctx.ReplyString("workspace added")
	}
}

func bertyBotAddChannel(db database, mutex *sync.Mutex) func(ctx bertybot.Context) {
	return func(ctx bertybot.Context) {
		// TODO: upgrade sdk to avoid it
		if ctx.IsReplay || !ctx.IsNew {
			return
		}

		workspaceName, channelName, err := func(ctx bertybot.Context) (string, string, error) {
			if len(ctx.CommandArgs) == 3 {
				return ctx.CommandArgs[1], ctx.CommandArgs[2], nil
			}
			return "", "", errors.New("bad arguments")
		}(ctx)
		if err != nil {
			_ = ctx.ReplyString(err.Error())
			return
		}
		//

		mutex.Lock()
		if db.ChannelExist(workspaceName, channelName) {
			_ = ctx.ReplyString("workspace or channel already exist")
			return
		}

		link, err := bertyBotCreateGroup(fmt.Sprintf("%s/#%s", workspaceName, channelName))
		if err != nil {
			_ = ctx.ReplyString(err.Error())
			return
		}

		err = db.AddChannel(workspaceName, channelName, link)
		mutex.Unlock()

		if err != nil {
			_ = ctx.ReplyString(err.Error())
			return
		}

		_ = ctx.ReplyString(link)
	}
}

func bertyBotListWorkspaces(db database) func(ctx bertybot.Context) {
	return func(ctx bertybot.Context) {
		// TODO: upgrade sdk to avoid it
		if ctx.IsReplay || !ctx.IsNew {
			return
		}

		workspaces, err := db.ListWorkspaces()
		if err != nil {
			_ = ctx.ReplyString("error: " + err.Error())
			return
		}
		_ = ctx.ReplyString("workspaces: " + strings.Join(workspaces, ", "))
	}
}

func bertyBotListChannels(db database) func(ctx bertybot.Context) {
	return func(ctx bertybot.Context) {
		// TODO: upgrade sdk to avoid it
		if ctx.IsReplay || !ctx.IsNew {
			return
		}

		workspaceID, err := func(ctx bertybot.Context) (string, error) {
			if len(ctx.CommandArgs) == 2 {
				return ctx.CommandArgs[1], nil
			}
			return "", errors.New("bad arguments")
		}(ctx)
		if err != nil {
			_ = ctx.ReplyString("error: " + err.Error())
			return
		}
		//

		channels, err := db.ListChannels(workspaceID)
		if err != nil {
			_ = ctx.ReplyString("error: " + err.Error())
			return
		}
		_ = ctx.ReplyString("channels: " + strings.Join(channels, ", "))
	}
}

func bertyBotRefreshAll() func(ctx bertybot.Context) {
	return func(ctx bertybot.Context) {
		_ = ctx.ReplyString("Not implemented yet!")
		fmt.Println("Not implemented yet!")
	}
}
