package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	gometrics "github.com/armon/go-metrics"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	"github.com/skill215/smpp-app/broker"
	"github.com/skill215/smpp-app/config"
	"github.com/skill215/smpp-app/logger"
	smppclient "github.com/skill215/smpp-app/smpp-client"
)

var (
	handler         *smppclient.SmppHandler
	b               *broker.Broker
	MetricsInterval = 5
)

func main() {
	ctx := context.Background()
	// get smpp app config
	conf, err := config.GetSmppConf()
	if err != nil {
		panic(err)
	}
	logger.SetupLogger()
	b = broker.NewBroker()
	go b.Start()

	// start metrics
	inm := gometrics.NewInmemSink(time.Duration(MetricsInterval)*time.Second, time.Minute)
	gometrics.NewGlobal(gometrics.DefaultConfig("smpp-app"), inm)
	go printMetrics(inm)

	// init smpp handler
	handler = smppclient.ProvideService(ctx, logrus.StandardLogger(), conf.App.SmppConn, b, inm)
	if err != nil {
		log.Fatal(err)
	}
	// start smpp app one by one
	handler.Init(ctx)
	addr := conf.GetRestAddr()

	http.HandleFunc("/startLoop", startLoop)
	http.HandleFunc("/stopLoop", stopLoop)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(addr), nil))
}

func startLoop(w http.ResponseWriter, r *http.Request) {
	tps, _ := strconv.Atoi(r.FormValue("tps"))
	b.Publish(tps)
}

func stopLoop(w http.ResponseWriter, r *http.Request) {
	b.Publish(0)
	handler.Stop(context.Background())
}

// Send Json in http response
func JSONResp(w http.ResponseWriter, resp interface{}, code int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(resp) //nolint:errcheck
}

func printMetrics(inm *gometrics.InmemSink) {
	ticker := time.NewTicker(5 * time.Second)
	for {
		<-ticker.C
		data := inm.Data()
		var interval *gometrics.IntervalMetrics
		n := len(data)
		switch {
		case n == 0:
			return
		case n == 1:
			// Show the current interval if it's all we have
			interval = data[0]
		default:
			// Show the most recent finished interval if we have one
			interval = data[n-2]
		}

		result := interval.Interval.String()
		output := map[string]int{}
		for _, counter := range interval.Counters {
			output[counter.Name] = counter.Count
		}
		for _, m := range []string{"ao", "ao failure", "at", "at failure"} {
			val, ok := output[m]
			if !ok {
				val = 0
			}
			result += fmt.Sprintf(" %s:%d ", m, val)
		}

		fmt.Println(result)
	}
}
