package main

import (
	"berty.tech/berty/v2/go/pkg/bertylinks"
	"berty.tech/berty/v2/go/pkg/bertyprotocol"
	"berty.tech/berty/v2/go/pkg/messengertypes"
)

func bertyBotCreateGroup(name string) (string, error) {
	g, _, err := bertyprotocol.NewGroupMultiMember()
	if err != nil {
		return "", err
	}

	group := &messengertypes.BertyGroup{
		Group:       g,
		DisplayName: name,
	}

	link := group.GetBertyLink()
	_, web, err := bertylinks.MarshalLink(link)
	if err != nil {
		return "", err
	}

	return web, nil
}
