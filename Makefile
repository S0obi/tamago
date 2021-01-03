BINARY_NAME=tamago
BINARY_OUTPUT_FOLDER_PATH=bin
SERVICE_ENTRYPOINT=cmd/tamago/main.go

.DEFAULT_GOAL = run

.PHONY: lint
lint:
	golint ./...

.PHONY: build
build:
	go build -gcflags="-m" -ldflags '-s -w' -o bin/$(BINARY_NAME) $(SERVICE_ENTRYPOINT)

.PHONY: release
release: build zip

.PHONY: zip
zip:
	cp -r assets bin/assets
	cd bin && zip -r tamago.zip tamago assets
	rm -rf bin/assets

.PHONY: run
run:
	go run $(SERVICE_ENTRYPOINT)

.PHONY: test
test:
	go test ./...