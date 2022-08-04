BERTY_PORT ?= 4298
BERTY_GROUP ?= https://berty.tech/id\#group/8ejngpAxnMPPKdsY8DA3jkWA2G4dU3tdh5qFk6Ym9g7Q5Y5zEcPR6tx4wUxjLvqN8aLySh2r7L6UNRXfNXokC9UoN1M5J38bQE5RsEe8orupzr83E1XjHn3y4cTBsVogHm2BFo84pWSGsdWxMVhtkXH3s4yDtuRZkYdB1BVqUvwi6B1s8G9pKSmT1vPMmiMUV3MHBe6pXxVnP3j7Jx5heZrkQeFQWPBeXVQ18ZysQmiroEtZjLPoncx8zSVBqkWMPATqu6JEuhcXN7qPVcN7cxGzxhoo/name=PMG+-+bot+experiments

run: tidy
	go run . --berty-node-addr=127.0.0.1:$(BERTY_PORT) --berty-group-invite="$(BERTY_GROUP)" --debug

start-mini-companion:
	berty daemon --node.listeners=/ip4/127.0.0.1/tcp/$(BERTY_PORT)/grpc --store.inmem --p2p.ipfs-api-listeners /ip4/127.0.0.1/tcp/1$(BERTY_PORT)

install: tidy fmt
	go install .

# utilities
tidy:; go mod tidy
fmt:;  gofumpt -w .
