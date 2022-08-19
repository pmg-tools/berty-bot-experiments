package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"strings"

	"berty.tech/berty/v2/go/pkg/bertybot"
	"github.com/Doozers/ETH-Signature/ethsign"
)

type metaData struct {
	Pubkey    string `json:"pubkey,omitempty"`
	PrevNonce string `json:"prev_nonce,omitempty"`
	PrevSig   string `json:"prev_sig,omitempty"`
	Sig       string `json:"sig,omitempty"`
	Nonce     string `json:"nonce,omitempty"`
	Message   string `json:"message,omitempty"`
	HashType  string `json:"hash_type,omitempty"`
}

type teritoriData struct {
	Step int `json:"step"`
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
	hash := sha256.Sum256([]byte(proof))
	sig := Sign((*[64]byte)(PublicKey), hash[:])
	b64sig := base64.StdEncoding.EncodeToString(sig)
	data := teritoriData{
		Step: 1,
		Data: metaData{
			Nonce: nonceStr,
			Sig:   b64sig,
		},
	}

	fmt.Println(data)
	return &data, nil
}

func step2(d database, t teritoriData) (*teritoriData, error) {
	if t.Data.PrevNonce == "" || t.Data.PrevSig == "" || t.Data.Pubkey == "" || t.Data.Sig == "" || t.Data.HashType == "" {
		return nil, errors.New("missing arg")
	}

	prevSig, err := base64.StdEncoding.DecodeString(t.Data.PrevSig)
	if err != nil {
		return nil, err
	}

	res, ok := Verify((*[32]byte)(PrivateKey), prevSig)

	hash32 := sha256.Sum256([]byte(fmt.Sprintf("%d%ssisi", t.Data.PrevNonce, t.Data.Pubkey+"sisi")))
	hash := hash32[:]
	if !ok || bytes.Equal(res, hash) {
		return nil, errors.New("invalid previous signature")
	}

	hash32_2 := sha256.Sum256([]byte(t.Data.PrevNonce + t.Data.PrevSig))
	hash2 := hash32_2[:]

	switch t.Data.HashType {
	case "text":
		ok, err = ethsign.Verify(fmt.Sprintf("%x", hash2), t.Data.Sig, t.Data.Pubkey, ethsign.TextHash)
	case "keccak256":
		ok, err = ethsign.Verify(fmt.Sprintf("%x", hash2), t.Data.Sig, t.Data.Pubkey, ethsign.Keccak256)
	}
	if err != nil {
		return nil, err
	}

	var data teritoriData
	if /* verify signature */ ok {
		data = teritoriData{
			Step: 3,
			Data: metaData{
				Message: "sync accepted",
			},
		}

		err = d.SyncTeritoriKey(t.Data.Pubkey, "bertyPubkey")
		if err != nil {
			return nil, err
		}
	} else {
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
