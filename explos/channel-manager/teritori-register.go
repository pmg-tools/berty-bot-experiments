package main

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"strings"

	"berty.tech/berty/v2/go/pkg/bertybot"
)

type teritoriData struct {
	Step int               `json:"step"`
	Data map[string]string `json:"data"`
}

func step0(t teritoriData) (*teritoriData, error) {
	if t.Data["pubkey"] == "" {
		return nil, errors.New("missing territoriPubKey")
	}

	// to modify
	pubKey, err := os.ReadFile("public.key")
	if err != nil {
		return nil, err
	}

	nonce := rand.Int()
	nonceStr := fmt.Sprintf("%d", nonce)
	proof := fmt.Sprintf("%d%ssisi", nonce, t.Data["pubkey"])
	sig := Sign((*[64]byte)(pubKey), []byte(proof))
	b64sig := base64.StdEncoding.EncodeToString(sig)
	data := teritoriData{
		Step: 1,
		Data: map[string]string{
			"nonce": nonceStr,
			"sig":   b64sig,
		},
	}
	if err != nil {
		return nil, err
	}
	fmt.Println(base64.StdEncoding.EncodeToString(Sign((*[64]byte)(pubKey), []byte(fmt.Sprintf("%d%ssisi", nonce, t.Data["pubkey"])))))
	return &data, nil
}

func step2(d database, t teritoriData) (*teritoriData, error) {
	if t.Data["prev_nonce"] == "" || t.Data["prev_sig"] == "" || t.Data["pubkey"] == "" || t.Data["sig"] == "" {
		return nil, errors.New("missing arg")
	}

	privKey, err := os.ReadFile("private.key")
	if err != nil {
		return nil, err
	}

	prevSig, err := base64.StdEncoding.DecodeString(t.Data["prev_sig"])
	if err != nil {
		return nil, err
	}

	res, ok := Verify((*[32]byte)(privKey), prevSig)
	fmt.Println(res)
	if !ok || string(res) != t.Data["prev_nonce"]+t.Data["pubkey"]+"sisi" {
		return nil, errors.New("invalid previous signature")
	}

	var data teritoriData
	if /* verify signature */ true {
		data = teritoriData{
			Step: 3,
			Data: map[string]string{
				"message": "sync accepted",
			},
		}
		if err != nil {
			return nil, err
		}

		err = d.SyncTeritoriKey(t.Data["pubkey"], "bertyPubkey")
		if err != nil {
			return nil, err
		}
		data = teritoriData{
			Step: 3,
			Data: map[string]string{
				"message": "sync rejected",
			},
		}
	}
	return &data, nil
}

func TeritoriAuth(d database) func(ctx bertybot.Context) {
	return func(ctx bertybot.Context) {
		data := strings.Replace(ctx.UserMessage, "/link-teritori-account ", "", 1)

		var t teritoriData
		err := json.Unmarshal([]byte(data), &t)
		if err != nil {
			ctx.ReplyString("error: " + err.Error())
		}

		var res *teritoriData
		switch t.Step {
		case 0:
			res, err = step0(t)

		case 2:
			res, err = step2(d, t)
			if err != nil {
				return
			}

		default:
			err = errors.New("error: unknown step")
		}
		
		if err != nil {
			ctx.ReplyString("error: " + err.Error())
			return
		}
		a, err := json.Marshal(res)
		if err != nil {
			ctx.ReplyString("error: " + err.Error())
			return
		}
		ctx.ReplyString(string(a))
	}
}
