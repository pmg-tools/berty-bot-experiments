package main

import (
	"berty.tech/berty/v2/go/pkg/bertyprotocol"
	"berty.tech/berty/v2/go/pkg/messengertypes"
)

func bertyBotCreateGroup(name string) (error, *messengertypes.BertyLink) {
	// TODO: upgrade sdk to avoid it
	g, _, err := bertyprotocol.NewGroupMultiMember()
	if err != nil {
		return err, nil
	}
	group := &messengertypes.BertyGroup{
		Group:       g,
		DisplayName: name,
	}
	link := group.GetBertyLink()
	return nil, link
}
