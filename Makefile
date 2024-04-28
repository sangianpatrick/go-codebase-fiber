.PHONY: install test-dev test cover run.dev build clean

install:
	go mod download

test:
	@echo "Run unit testing ..."
		mkdir -p ./coverage && \
			go test -v -coverprofile=./coverage/coverage.out -covermode=atomic ./...

cover: test
	@echo "Generating coverprofile ..."
		go tool cover -func=./coverage/coverage.out &&\
			go tool cover -html=./coverage/coverage.out -o ./coverage/coverage.html

run.dev:
	@echo "Run in development mode ..."
		GOOGLE_APPLICATION_CREDENTIALS=/home/patrick/Documents/tsel-assessment/tsel-ticketmaster-github-action.json go run cmd/web/main.go

build:
	@echo "Building the executable file ..."
		CGO_ENABLED=1 GOOS=linux go build -tags musl -a -o bin/web cmd/web/main.go &&\
			cp bin/web /tmp/web

clean:
	@echo "Cleansing the last built ..."
		rm -rf bin