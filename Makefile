build: server client

server: clean
	GOOS=linux GOARCH=amd64 go build server.go
	zip server.zip server

client: clean
	go build client.go

clean:
	-rm server client
