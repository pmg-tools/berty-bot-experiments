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

type territoriData struct {
	Step int               `json:"step"`
	Data map[string]string `json:"data"`
}

func step0(ctx bertybot.Context, t territoriData) {
	if t.Data["pubkey"] == "" {
		ctx.ReplyString("error: missing territoriPubKey")
		return
	}

	// to modify
	pubKey, err := os.ReadFile("public.key")
	if err != nil {
		ctx.ReplyString("error: " + err.Error())
		return
	}

	nonce := rand.Int()

	m, err := json.Marshal(territoriData{
		Step: 1,
		Data: map[string]string{
			"nonce": fmt.Sprintf("%d", nonce),
			"sig":   base64.StdEncoding.EncodeToString(Sign((*[64]byte)(pubKey), []byte(fmt.Sprintf("%d%ssisi", nonce, t.Data["pubkey"])))),
		},
	})
	if err != nil {
		ctx.ReplyString("error: " + err.Error())
		return
	}
	fmt.Println(base64.StdEncoding.EncodeToString(Sign((*[64]byte)(pubKey), []byte(fmt.Sprintf("%d%ssisi", nonce, t.Data["pubkey"])))))
	ctx.ReplyString(string(m))
}

func step2(ctx bertybot.Context, d database, t territoriData) {
	if t.Data["prev_nonce"] == "" || t.Data["prev_sig"] == "" || t.Data["pubkey"] == "" || t.Data["sig"] == "" {
		ctx.ReplyString("error: missing arg")
		return
	}

	privKey, err := os.ReadFile("private.key")
	if err != nil {
		ctx.ReplyString("error: " + err.Error())
		return
	}

	prevSig, err := base64.StdEncoding.DecodeString(t.Data["prev_sig"])
	if err != nil {
		ctx.ReplyString("error: " + err.Error())
		return
	}

	res, ok := Verify((*[32]byte)(privKey), prevSig)
	fmt.Println(res)
	if !ok || string(res) != t.Data["prev_nonce"]+t.Data["pubkey"]+"sisi" {
		ctx.ReplyString("error: invalid previous signature")
		return
	}

	if /* verify signature */ true {
		m, err := json.Marshal(territoriData{
			Step: 3,
			Data: map[string]string{
				"message": "sync accepted",
			},
		})
		if err != nil {
			ctx.ReplyString("error: " + err.Error())
			return
		}

		err = d.SyncTeritoriKey(t.Data["pubkey"], "bertyPubkey")
		if err != nil {
			ctx.ReplyString("error: " + err.Error())
			return
		}
		ctx.ReplyString(string(m))
		return
	}

	m, err := json.Marshal(territoriData{
		Step: 3,
		Data: map[string]string{
			"message": "sync rejected",
		},
	})
	if err != nil {
		ctx.ReplyString("error: " + err.Error())
		return
	}
	ctx.ReplyString(string(m))
}

func TerritoriAuth(d database) func(ctx bertybot.Context) {
	return func(ctx bertybot.Context) {
		data := strings.Replace(ctx.UserMessage, "/link-territori-account ", "", 1)

		var t territoriData
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
