package main

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

func TestBertyBotCreateGroup(t *testing.T) {
	name := "group-test"
	link, err := bertyBotCreateGroup(name)
	require.NoError(t, err)
	fmt.Println(link)
	prefix := "https://berty.tech/id#group/"
	result := strings.HasPrefix(link, prefix)
	require.Equal(t, true, result, "The link should start with the prefix defined above")
}
