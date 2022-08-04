# 3-bot-rel-clone-go-dep

## About

A bot using a cloned berty/berty repo somewhere on the file system.

This method allows:
* patching berty/berty in real time.
* staying in the berty/berty repo to make opening a PR very easy.

Drawback -> "works for me, but not for others" (require people to make more steps)

## How

1. `git clone https://github.com/berty/berty ../berty-clone`
2. Patch it.
3. `make start-mini-companion`
4. `make run`
