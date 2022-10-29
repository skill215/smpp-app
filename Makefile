hello:
	echo "Hello"

build:
	GOOS=linux CGO_ENABLED=0 go build -o target/rest4smpp server.go smpplog.go config.go smppapp.go broker.go model.go counter.go stat.go smppappv2.go wsp.go bcd.go ref.go
	cp conf.yaml target/conf.yaml
	cp webbing_apn.wbxml target/webbing_apn.wbxml

mac:
	GOOS=darwin CGO_ENABLED=0 go build -o target/rest4smpp_mac server.go smpplog.go config.go smppapp.go broker.go model.go counter.go stat.go smppappv2.go wsp.go bcd.go ref.go
	cp conf.yaml target/conf.yaml
run:
	go run smppapp.go

clean:
	rm -rf target/*
	rm -rf rest4smpp.log

test:
	go test -timeout 30s -tags unit,integration smppapp -v -count=1