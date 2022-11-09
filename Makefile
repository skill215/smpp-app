# Set environment for Go command
export GOBIN			?= $(shell go env GOPATH)/bin
export GOPRIVATE		= https://github.com/skill215/smpp-app
export GOPROXY			= https://goproxy.io,direct
export GOSUMDB			= off


hello:
	echo "Hello"

build:
	GOOS=linux CGO_ENABLED=0 go build -o target/receiver rest-server.go config.go smpp-receiver.go
	cp smpp-app.yaml target/smpp-app.yaml

clean:
	rm target/*