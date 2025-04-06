build:
	@GOOS=linux GOARCH=amd64 go build  -o bin/inventory ./cmd

build-win:
	@go build -o bin/inventory.exe  ./cmd

run: build
	./bin/inventory

run-win: build-win
	./bin/inventory.exe

test:
	go test -v ./...