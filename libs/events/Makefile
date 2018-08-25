.PHONY: docs
REPO:=github.com/blockchainworkers/conch/libs/events

docs:
	@go get github.com/davecheney/godoc2md
	godoc2md $(REPO) > README.md

test:
	go test -v ./...
