build:
	go build -o conv-relay cmd/main.go
run:
	./conv-relay
ngrok:
	~/ngrok-2 http 8080