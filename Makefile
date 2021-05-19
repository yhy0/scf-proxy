build: server client

server: clean
	GOOS=linux GOARCH=amd64 go build -trimpath server.go
	zip server.zip server

client: clean
	go build -trimpath client.go

fmt:
	go fmt ./...

clean:
	-rm server server.zip client
