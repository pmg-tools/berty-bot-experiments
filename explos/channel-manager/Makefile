BERTY_PORT ?= 4296
BERTY_GROUP ?=https://berty.tech/id\#group/8ejngpAxnMPpDAAsp2LLAAKMEpFipB2EEGp8kMQnw1nhBavAZPiGLWju1C2ggCNuJ7RwkpGfaB2hdX8jiuyvnxwCJnEiiQPKgzP9qo5oGCiMWTfDGks9gnJnQnzkQC7XzFUgMp3MvuxtDKJqCjhy7dBKYABZkAUWjEjerujQRoJHbQto1P8tWcRiHxLFXurKX4Y9DS6moGcyMz8CmKaAvuf8D4Va8Z6cJSiHBzXFz7e158K8ruH4amk1qEMPDkT78xLycotAqajfnNaUBBd5Ua6hDAC1/name=Moi%2Bbots

run: tidy
	go run . -berty-node-addr=127.0.0.1:$(BERTY_PORT) --berty-group-invite="$(BERTY_GROUP)" --debug

start-mini-companion:
	berty daemon --node.listeners=/ip4/127.0.0.1/tcp/$(BERTY_PORT)/grpc --p2p.ipfs-api-listeners /ip4/127.0.0.1/tcp/1$(BERTY_PORT)

install: tidy fmt
	go install .

# utilities
tidy:; go mod tidy
fmt:;  gofumpt -w .
