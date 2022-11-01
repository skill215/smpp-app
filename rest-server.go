package main

import (
	"context"
)

func main() {
	ctx := context.Background()
	conf, err := GetSmppConf()
	if err != nil {
		panic(err)
	}
	receiver, err := ProvideSmppReceiver(ctx, conf)
	if err != nil {
		panic(err)
	}
	if err = receiver.bind(); err != nil {
		panic(err)
	}
	for i := 0; i < 10; i++ {
		go receiver.bind()
	}
	<-ctx.Done()
}
