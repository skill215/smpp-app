hello:
	echo "Hello"

build:
	GOOS=linux CGO_ENABLED=0 go build -o target/receiver rest-server.go config.go smpp-receiver.go
	cp smpp-app.yaml target/smpp-app.yaml
