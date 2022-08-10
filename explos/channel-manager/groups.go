package main

import (
	"berty.tech/berty/v2/go/pkg/bertyprotocol"
	"berty.tech/berty/v2/go/pkg/messengertypes"
)

func bertyBotCreateGroup(name string) (*messengertypes.BertyLink, error) {
	g, _, err := bertyprotocol.NewGroupMultiMember()
	if err != nil {
		return nil, err
	}

	group := &messengertypes.BertyGroup{
		Group:       g,
		DisplayName: name,
	}
	link := group.GetBertyLink()
	return link, nil
}
