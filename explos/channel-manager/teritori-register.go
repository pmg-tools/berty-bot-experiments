package main

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"strings"

	"berty.tech/berty/v2/go/pkg/bertybot"
)

type metaData struct {
	Pubkey    string `json:"pubkey,omitempty"`
	PrevNonce string `json:"prev_nonce,omitempty"`
	PrevSig   string `json:"prev_sig,omitempty"`
	Sig       string `json:"sig,omitempty"`
	Nonce     string `json:"nonce,omitempty"`
	Message   string `json:"message,omitempty"`
}

type teritoriData struct {
	Step int `json:"step"`
	//Data map[string]string `json:"data"`
	Data metaData
}

func step0(t teritoriData) (*teritoriData, error) {
	if t.Data.Pubkey == "" {
		return nil, errors.New("missing territoriPubKey")
	}

	// to modify
	nonce := rand.Int()
	nonceStr := fmt.Sprintf("%d", nonce)
	proof := fmt.Sprintf("%d%ssisi", nonce, t.Data.Pubkey)
	sig := Sign((*[64]byte)(PublicKey), []byte(proof))
	b64sig := base64.StdEncoding.EncodeToString(sig)
	data := teritoriData{
		Step: 1,
		Data: metaData{
			Nonce: nonceStr,
			Sig:   b64sig,
		},
	}

	fmt.Println(base64.StdEncoding.EncodeToString(Sign((*[64]byte)(PublicKey), []byte(fmt.Sprintf("%d%ssisi", nonce, t.Data.Pubkey)))))
	return &data, nil
}

func step2(d database, t teritoriData) (*teritoriData, error) {
	if t.Data.PrevNonce == "" || t.Data.PrevSig == "" || t.Data.Pubkey == "" || t.Data.Sig == "" {
		return nil, errors.New("missing arg")
	}

	prevSig, err := base64.StdEncoding.DecodeString(t.Data.PrevSig)
	if err != nil {
		return nil, err
	}

	res, ok := Verify((*[32]byte)(PrivateKey), prevSig)
	fmt.Println(res)
	if !ok || string(res) != t.Data.PrevNonce+t.Data.Pubkey+"sisi" {
		return nil, errors.New("invalid previous signature")
	}

	var data teritoriData
	if /* verify signature */ true {
		data = teritoriData{
			Step: 3,
			Data: metaData{
				Message: "sync accepted",
			},
		}
		if err != nil {
			return nil, err
		}

		err = d.SyncTeritoriKey(t.Data.Pubkey, "bertyPubkey")
		if err != nil {
			return nil, err
		}
		data = teritoriData{
			Step: 3,
			Data: metaData{
				Message: "sync rejected",
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
