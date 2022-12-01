package smppapp

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"
	"github.com/skill215/smpp-app/config"
	smppclient "github.com/skill215/smpp-app/smpp-client"
)

var (
	event = make(chan int)
)

func main() {
	ctx := context.Background()
	// get smpp app config
	conf, err := config.GetSmppConf()
	if err != nil {
		panic(err)
	}

	handler, err := smppclient.ProvideService(ctx, log, &conf.App.SmppConn, event)
	if err != nil {
		log.Fatal(err)
	}
	// start smpp app one by one
	handler.Init(ctx)
	handler.Run(ctx)
	addr := conf.GetRestAddr()

	http.HandleFunc("/startLoop", startLoop)
	http.HandleFunc("/stopLoop", stopLoop)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(addr), nil))
}

func startLoop(w http.ResponseWriter, r *http.Request) {
	event <- 1000
}

func stopLoop(w http.ResponseWriter, r *http.Request) {
	event <- 0
}

// Send Json in http response
func JSONResp(w http.ResponseWriter, resp interface{}, code int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(resp) //nolint:errcheck
}
