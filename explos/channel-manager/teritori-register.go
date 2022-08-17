package main

import (
	"berty.tech/berty/v2/go/pkg/bertybot"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"strings"
)

type teritoriData struct {
	Step int               `json:"step"`
	Data map[string]string `json:"data"`
}

func step0(ctx bertybot.Context, t teritoriData) {
	if t.Data["PubKey"] == "" {
		ctx.ReplyString("error: missing teritoriPubKey")
		return
	}

	// to modify
	pubKey, err := os.ReadFile("public.key")
	if err != nil {
		ctx.ReplyString("error: " + err.Error())
		return
	}

	nonce := rand.Int()

	m, err := json.Marshal(teritoriData{
		Step: 1,
		Data: map[string]string{
			"nonce": fmt.Sprintf("%d", nonce),
			"sig":   base64.StdEncoding.EncodeToString(Sign((*[64]byte)(pubKey), []byte(fmt.Sprintf("%d%ssisi", nonce, t.Data["teritoriPubKey"])))),
		},
	})
	if err != nil {
		ctx.ReplyString("error: " + err.Error())
		return
	}
	ctx.ReplyString(string(m))
}

func step2(ctx bertybot.Context, d database, t teritoriData) {
	if t.Data["prev_nonce"] == "" || t.Data["prev_sig"] == "" || t.Data["PubKey"] == "" || t.Data["sig"] == "" {
		ctx.ReplyString("error: missing arg")
		return
	}

	privKey, err := os.ReadFile("private.key")
	if err != nil {
		ctx.ReplyString("error: " + err.Error())
		return
	}

	res, ok := Verify((*[32]byte)(privKey), []byte(t.Data["prev_sig"]))
	if !ok || !strings.Contains(string(res), t.Data["prev_nonce"]) {
		ctx.ReplyString("error: invalid previous signature")
		return
	}

	if /* verify signature */ true {
		if ok := d.ConfirmUser(t.Data["teritoriPubKey"], "bertyPubKey"); !ok {
			ctx.ReplyString("error: user not found")
			return
		}
		m, err := json.Marshal(teritoriData{
			Step: 3,
			Data: map[string]string{
				"message": "accepted",
			},
		})
		if err != nil {
			ctx.ReplyString("error: " + err.Error())
			return
		}
		ctx.ReplyString(string(m))
	}

	m, err := json.Marshal(teritoriData{
		Step: 3,
		Data: map[string]string{
			"message": "rejected",
		},
	})
	if err != nil {
		ctx.ReplyString("error: " + err.Error())
		return
	}
	ctx.ReplyString(string(m))
}

func teritoriAuth(d database) func(ctx bertybot.Context) {
	return func(ctx bertybot.Context) {
		data := strings.Replace(ctx.UserMessage, "/link-teritori-account ", "", 1)

		var t teritoriData
		err := json.Unmarshal([]byte(data), &t)
		if err != nil {
			ctx.ReplyString("error: " + err.Error())
		}

		switch t.Step {
		case 0:
			step0(ctx, t)
			break
		case 2:
			step2(ctx, d, t)
			break
		default:
			ctx.ReplyString("error: unknown step")
		}
	}
}
