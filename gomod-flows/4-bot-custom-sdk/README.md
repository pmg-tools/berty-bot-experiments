# 4-bot-custom-sdk

## About

This one is using the official `berty/berty` repo for most libs, excepted the bot SDK which is copied with the project.

This method makes it ultra easy to make changes in the SDK, easy to show it to anyone (no need to make secondary clone or checkout something).
But, it makes it harder to open PRs against the original repo.

---

FYI, the lib was duplicated like this:

```sh
# copy all the deps' sources in vendor/ dir
go mod vendor

# copy the interesting library
cp -rf vendor/berty.tech/berty/v2/go/pkg/bertybot .

# cleanup vendor dir
rm -rf vendor

# monkey patch files that requires berty's internals
$EDITOR ./bertybot/recipes.go
```

## How

1. `make start-mini-companion`
2. `make run`
