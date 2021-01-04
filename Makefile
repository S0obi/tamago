BINARY_NAME=tamago
BINARY_OUTPUT_FOLDER_PATH=bin
SERVICE_ENTRYPOINT=cmd/tamago/main.go

.DEFAULT_GOAL = run

.PHONY: lint
lint:
	golint ./...

build: $(SERVICE_ENTRYPOINT)
	go build -gcflags="-m" -ldflags '-s -w' -o bin/$(BINARY_NAME) $(SERVICE_ENTRYPOINT)

.PHONY: release
release: build zip clean

zip: assets bin/tamago
	cp -r assets bin/assets
	cd bin && zip -r tamago.zip tamago assets

.PHONY: clean
clean:
	rm bin/tamago
	rm -rf bin/assets

.PHONY: run
run:
	go run $(SERVICE_ENTRYPOINT)

.PHONY: test
test:
	go test ./...