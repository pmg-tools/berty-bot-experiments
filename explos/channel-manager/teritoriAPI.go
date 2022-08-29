package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type refreshData struct {
	Access []struct {
		Workspace string   `json:"workspace"`
		Channel   []string `json:"channel"`
	} `json:"access"`
}

func requestUserAccess(api string, pubKey string) (*refreshData, error) {
	var data refreshData
	client := http.Client{}
	req, err := http.NewRequest("GET", api, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("pubKey", pubKey)
	do, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer do.Body.Close()

	body, err := ioutil.ReadAll(do.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}
	return &data, nil
}
