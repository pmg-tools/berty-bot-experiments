package main

import (
	"berty.tech/berty/v2/go/pkg/bertybot"
	"berty.tech/berty/v2/go/pkg/messengertypes"
	"errors"
	"fmt"
	"github.com/gogo/protobuf/proto"
	"strings"
	"sync"
)

type database interface {
	AddUser(bertyPubKey string) error
	SyncTeritoriKey(teritoriPubkey string, bertyPubkey string) error

	AddWorkspace(workspaceName string) error
	AddChannel(workspaceName string, channelName string, bertyGroupLink string) error

	GetChannelsInvitation(workspaceName string, channelsName []string) []Channel

	// debug functions
	ListWorkspaces() ([]string, error)
	ListChannels(workspaceName string) ([]string, error)
	ListUsers() ([]User, error)
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

func bertyBotRefresh(db database, api string) func(ctx bertybot.Context) {
	return func(ctx bertybot.Context) {
		// TODO: upgrade sdk to avoid it
		if ctx.IsReplay || !ctx.IsNew {
			return
		}

		pubKey, err := func(ctx bertybot.Context) (string, error) {
			if len(ctx.CommandArgs) == 2 {
				return ctx.CommandArgs[1], nil
			}
			return "", errors.New("bad arguments")
		}(ctx)
		if err != nil {
			_ = ctx.ReplyString(err.Error())
			return
		}

		data, err := requestUserAccess(api, pubKey)
		if err != nil {
			_ = ctx.ReplyString(err.Error())
			return
		}

		var channels []Channel
		for _, v := range data.Access {
			channels = db.GetChannelsInvitation(v.Workspace, v.Channel)
			for _, w := range channels {
				_ = ctx.ReplyString(w.BertyLink)
			}
		}
	}
}

func bertyBotRefreshAll(db database, api string) func(ctx bertybot.Context) {
	return func(ctx bertybot.Context) {
		users, err := db.ListUsers()
		if err != nil {
			_ = ctx.ReplyString("error: " + err.Error())
		}

		var data *RefreshData
		var channels []Channel
		for _, v := range users {
			data, err = requestUserAccess(api, v.BertyPubKey)
			for _, w := range data.Access {
				channels = db.GetChannelsInvitation(w.Workspace, w.Channel)
				for _, x := range channels {
					userMessage, err := proto.Marshal(&messengertypes.AppMessage_UserMessage{Body: x.BertyLink})
					if err != nil {
						_ = ctx.ReplyString("error: " + err.Error())
					}
					_, err = ctx.Client.Interact(ctx.Context, &messengertypes.Interact_Request{
						Type:                  messengertypes.AppMessage_TypeUserMessage,
						Payload:               userMessage,
						ConversationPublicKey: v.BertyPubKey,
					})
					if err != nil {
						_ = ctx.ReplyString("error: " + err.Error())
					}
				}
			}
		}
	}
}
