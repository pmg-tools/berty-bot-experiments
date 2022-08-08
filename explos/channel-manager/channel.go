package main

import (
	"berty.tech/berty/v2/go/pkg/bertybot"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"
)

type database interface {
	AddChannel(channelID string) error
	ListChannels() ([]string, error)
	GetChannelMessages(channelID string) []message
	ChannelExist(channelID string) bool

	AddMessage(channelID string, msg message)
}

func bertyBotAddChannel(db database, mutex *sync.Mutex) func(ctx bertybot.Context) {
	return func(ctx bertybot.Context) {
		// TODO: upgrade sdk to avoid it
		channelID, err := func(ctx bertybot.Context) (string, error) {
			if len(ctx.CommandArgs) == 2 {
				return ctx.CommandArgs[1], nil
			}
			return "", errors.New("missing channel ID")
		}(ctx)
		if err != nil {
			_ = ctx.ReplyString(err.Error())
			return
		}

		mutex.Lock()
		err = db.AddChannel(channelID)
		mutex.Unlock()
		if err != nil {
			_ = ctx.ReplyString("error: " + err.Error())
			return
		}
		_ = ctx.ReplyString("channel added")
	}
}

func bertyBotListChannels(db database) func(ctx bertybot.Context) {
	return func(ctx bertybot.Context) {
		channels, err := db.ListChannels()
		if err != nil {
			_ = ctx.ReplyString("error: " + err.Error())
			return
		}
		_ = ctx.ReplyString("channels: " + strings.Join(channels, ", "))
	}
}

func bertyBotAddMessage(db database, mutex *sync.Mutex) func(ctx bertybot.Context) {
	return func(ctx bertybot.Context) {
		// TODO: upgrade sdk to avoid it
		channelID, messageText, err := func(ctx bertybot.Context) (string, string, error) {
			if len(ctx.CommandArgs) == 3 {
				return ctx.CommandArgs[1], ctx.CommandArgs[2], nil
			}
			return "", "", errors.New("missing channel ID")
		}(ctx)
		if err != nil {
			_ = ctx.ReplyString(err.Error())
			return
		}

		if !db.ChannelExist(channelID) {
			_ = ctx.ReplyString("error: channel does not exist")
			return
		}

		msg := message{
			Text:      messageText,
			ChannelID: channelID,
			hour:      time.Now(),
		}

		mutex.Lock()
		db.AddMessage(channelID, msg)
		mutex.Unlock()

		_ = ctx.ReplyString("message added")
	}
}

func FormatMessage(msg message) string {
	return fmt.Sprintf("%s: %s", msg.hour.Format("15:04"), msg.Text)
}

func bertyBotGetMessages(db database) func(ctx bertybot.Context) {
	return func(ctx bertybot.Context) {
		// TODO: upgrade sdk to avoid it
		channelID, err := func(ctx bertybot.Context) (string, error) {
			if len(ctx.CommandArgs) == 2 {
				return ctx.CommandArgs[1], nil
			}
			return "", errors.New("missing channel ID")
		}(ctx)
		if err != nil {
			_ = ctx.ReplyString(err.Error())
			return
		}

		if !db.ChannelExist(channelID) {
			_ = ctx.ReplyString("error: channel does not exist")
			return
		}

		messages := db.GetChannelMessages(channelID)

		var finalMessage string
		for _, msg := range messages {
			finalMessage += FormatMessage(msg) + "\n"
		}
		_ = ctx.ReplyString(finalMessage)
	}
}
