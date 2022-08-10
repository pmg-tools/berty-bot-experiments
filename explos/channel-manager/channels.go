package main

import (
	"berty.tech/berty/v2/go/pkg/bertybot"
	"errors"
	"strings"
	"sync"
)

type database interface {
	AddWorkspace(workspaceID string)
	AddChannel(workspaceID string, channelID string)

	WorkspaceExist(workspaceID string) bool
	ChannelExist(workspaceID string, channelID string) bool

	// debug functions
	ListWorkspaces() ([]string, error)
	ListChannels(workspaceID string) ([]string, error)
}

func bertyBotAddWorkspace(db database, mutex *sync.Mutex) func(ctx bertybot.Context) {
	return func(ctx bertybot.Context) {
		// TODO: upgrade sdk to avoid it
		if ctx.IsReplay || !ctx.IsNew {
			return
		}

		workspaceID, err := func(ctx bertybot.Context) (string, error) {
			if len(ctx.CommandArgs) == 2 {
				return ctx.CommandArgs[1], nil
			}
			return "", errors.New("missing workspace ID")
		}(ctx)
		if err != nil {
			_ = ctx.ReplyString(err.Error())
			return
		}
		//

		mutex.Lock()
		db.AddWorkspace(workspaceID)
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

		workspaceID, channelID, err := func(ctx bertybot.Context) (string, string, error) {
			if len(ctx.CommandArgs) == 3 {
				return ctx.CommandArgs[1], ctx.CommandArgs[2], nil
			}
			return "", "", errors.New("missing channel ID")
		}(ctx)
		if err != nil {
			_ = ctx.ReplyString(err.Error())
			return
		}
		//

		mutex.Lock()
		db.AddChannel(workspaceID, channelID)
		mutex.Unlock()

		_ = ctx.ReplyString("channel added")
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
			return "", errors.New("missing workspace ID")
		}(ctx)
		if err != nil {
			_ = ctx.ReplyString(err.Error())
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
