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
	AddUser(pubKey string) bool

	AddWorkspace(workspaceName string) bool
	AddChannel(workspaceName string, channelName string, bertyGroupLink string) bool

	WorkspaceExist(workspaceName string) bool
	ChannelExist(workspaceName string, channelName string) bool

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

		if db.WorkspaceExist(workspaceName) {
			_ = ctx.ReplyString("workspace already exists")
			return
		}

		mutex.Lock()
		db.AddWorkspace(workspaceName)
		mutex.Unlock()

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

		if db.ChannelExist(workspaceName, channelName) {
			_ = ctx.ReplyString("workspace or channel already exist")
			return
		}

		link, err := bertyBotCreateGroup(channelName)
		if err != nil {
			_ = ctx.ReplyString(err.Error())
		}

		mutex.Lock()
		db.AddChannel(workspaceName, channelName, link)
		mutex.Unlock()

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
