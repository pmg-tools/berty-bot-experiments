package main

import (
	"berty.tech/berty/v2/go/pkg/bertybot"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
)

type RefreshData struct {
	Access []struct {
		Workspace string   `json:"workspace"`
		Channel   []string `json:"channel"`
	} `json:"access"`
}

func refreshUser(db database, api string) func(ctx bertybot.Context) {
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
		//

		var data RefreshData
		client := http.Client{}
		req, err := http.NewRequest("GET", api, nil)
		if err != nil {
			_ = ctx.ReplyString(err.Error())
			return
		}
		req.Header.Set("pubKey", pubKey)
		do, err := client.Do(req)
		if err != nil {
			_ = ctx.ReplyString(err.Error())
			return
		}

		defer do.Body.Close()

		body, err := ioutil.ReadAll(do.Body)
		if err != nil {
			return
		}

		err = json.Unmarshal(body, &data)
		if err != nil {
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
