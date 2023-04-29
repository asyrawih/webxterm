run:
	@echo "Running Server"
	@go run ./cmd/xterm_server/main.go

build:
	env GOOS=linux   go build -o server_linux ./cmd/xterm_server/ 
	env GOOS=darwin  go build -o server_darwin ./cmd/xterm_server/ 
