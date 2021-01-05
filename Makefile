BINARY_NAME=tamago
SERVICE_ENTRYPOINT=cmd/tamago/main.go

.DEFAULT_GOAL = run

.PHONY: lint
lint:
	golint ./...

build: $(SERVICE_ENTRYPOINT)
	go build -gcflags="-m" -ldflags '-s -w' -o bin/$(BINARY_NAME) $(SERVICE_ENTRYPOINT)

.PHONY: release
release: build zip clean

zip: assets bin/$(BINARY_NAME)
	cp -r assets bin/assets
	cd bin && zip -r tamago.zip $(BINARY_NAME) assets

.PHONY: clean
clean:
	rm bin/$(BINARY_NAME)
	rm -rf bin/assets

.PHONY: run
run:
	go run $(SERVICE_ENTRYPOINT)

.PHONY: test
test:
	go test ./...