package main

import (
	"berty.tech/berty/v2/go/pkg/bertybot"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
)

type territoriData struct {
	Step string            `json:"step"`
	Data map[string]string `json:"data"`
}

func TerritoriAuth(d database) func(ctx bertybot.Context) {
	return func(ctx bertybot.Context) {
		// TODO: upgrade sdk to avoid it
		if ctx.IsReplay || !ctx.IsNew {
			return
		}

		fmt.Println(ctx.CommandArgs)
		data, err := func(ctx bertybot.Context) (string, error) {
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

		var t territoriData
		err = json.Unmarshal([]byte(data), &t)
		if err != nil {
			ctx.ReplyString("error: " + err.Error())
		}

		switch t.Step {
		case "auth":
			if t.Data["territoriPubKey"] == "" {
				ctx.ReplyString("error: missing territoriPubKey")
				return
			}

			nonce := rand.Int()
			d.AddUser(t.Data["territoriPubKey"], "bertyPubKey", nonce)
			ctx.ReplyString("auth")
			break
		case "confirm":
			if /* verify signature */ true {
				if ok := d.ConfirmUser(t.Data["territoriPubKey"], "bertyPubKey"); !ok {
					ctx.ReplyString("error: user not found")
					return
				}
				ctx.ReplyString("setup confirmed !")
			}
			break
		default:
			ctx.ReplyString("error: unknown step")
		}

	}
}
