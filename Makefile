build:
	go build -o conv-relay cmd/main.go
run:
	./conv-relay
ngrok:
	~/ngrok http 8080