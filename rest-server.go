package smppapp

import (
	"context"
	"fmt"
	"log"
	"net/http"
)

type SmppApp interface {
	Init()
	Start()
	Stop()
}

type App struct {
	apps []SmppApp
	conf *AppConfig
}

func (a *App) Init() {}

func main() {
	ctx := context.Background()
	// get smpp app config
	conf, err := GetSmppConf()
	if err != nil {
		panic(err)
	}
	app := App{conf: conf}
	// start smpp app one by one
	app.Init()

	addr := conf.getRestAddr()

	http.HandleFunc("/startLoop", startLoop)
	http.HandleFunc("/stopLoop", stopLoop)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(addr), nil))
}
