package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
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

func printUsage() {
	fmt.Println("SMPP Application REST Server")
	fmt.Println("\nUsage:")
	fmt.Printf("  %s [options]\n", os.Args[0])
	fmt.Println("\nOptions:")
	fmt.Println("  -c string")
	fmt.Println("        Configuration file path (default: smpp-app.yaml)")
	fmt.Println("  -server-port uint")
	fmt.Println("        REST server port (overrides port in config file)")
	fmt.Println("\nExample:")
	fmt.Println("  Start with default configuration:")
	fmt.Println("    ./rest-server -c config/smpp-app.yaml")
	fmt.Println("  Start with custom REST port:")
	fmt.Println("    ./rest-server -c config/smpp-app.yaml -server-port 8082")
	fmt.Println("\nAPI Endpoints:")
	fmt.Println("  /startLoop?tps=<number>  Start sending messages with specified TPS")
	fmt.Println("  /stopLoop                Stop sending messages")
	fmt.Println("\nMetrics:")
	fmt.Printf("  Metrics are printed every %d seconds showing:\n", MetricsInterval)
	fmt.Println("    ao: Number of messages sent")
	fmt.Println("    ao failure: Number of failed messages")
	fmt.Println("    at: Number of messages received")
	fmt.Println("    at failure: Number of failed receives")
}

func main() {
	flag.Usage = printUsage
	confPath := flag.String("c", "smpp-app.yaml", "configuration file path")
	serverPort := flag.Uint("server-port", 0, "REST server port (overrides config file)")
	flag.Parse()

	if len(os.Args) == 1 {
		printUsage()
		os.Exit(1)
	}

	ctx := context.Background()
	// get smpp app config
	conf, err := config.GetSmppConf(*confPath)
	if err != nil {
		log.WithError(err).Fatal("Failed to load configuration")
	}
	logger.SetupLogger()

	// Set log level from config
	if level, err := logrus.ParseLevel(conf.App.Log.Level); err != nil {
		log.Warnf("Invalid log level %s, using default", conf.App.Log.Level)
	} else {
		log.SetLevel(level)
		log.Debugf("Log level set to %s", level)
	}

	// Override REST port if specified in command line
	if *serverPort > 0 {
		log.WithFields(logrus.Fields{
			"config_port": conf.App.Rest.Port,
			"cmd_port":    *serverPort,
		}).Debug("Overriding REST port from command line")
		conf.App.Rest.Port = uint16(*serverPort)
	}

	log.Debug("Starting rest-server...")

	b = broker.NewBroker()
	go b.Start()
	log.Debug("Broker started")

	// start metrics
	inm := gometrics.NewInmemSink(time.Duration(MetricsInterval)*time.Second, time.Minute)
	gometrics.NewGlobal(gometrics.DefaultConfig("smpp-app"), inm)
	go printMetrics(inm)
	log.Debug("Metrics initialized")

	// init smpp handler
	handler = smppclient.ProvideService(ctx, logrus.StandardLogger(), conf.App.SmppConn, b, inm)
	if err != nil {
		log.Fatal(err)
	}
	// start smpp app one by one
	handler.Init(ctx)
	addr := conf.GetRestAddr()
	log.WithFields(logrus.Fields{
		"rest_addr": conf.App.Rest.Addr,
		"rest_port": conf.App.Rest.Port,
		"full_addr": addr,
	}).Info("REST server configuration")

	http.HandleFunc("/startLoop", startLoop)
	http.HandleFunc("/stopLoop", stopLoop)
	log.Debug("HTTP endpoints registered")
	log.Fatal(http.ListenAndServe(addr, nil))
}

func startLoop(w http.ResponseWriter, r *http.Request) {
	log.WithFields(log.Fields{
		"method": r.Method,
		"path":   r.URL.Path,
		"query":  r.URL.RawQuery,
	}).Debug("Received startLoop request")

	tps, err := strconv.Atoi(r.FormValue("tps"))
	if err != nil {
		log.WithError(err).Error("Invalid TPS parameter")
		JSONResp(w, map[string]string{"error": "Invalid TPS parameter"}, http.StatusBadRequest)
		return
	}

	log.WithFields(log.Fields{
		"tps": tps,
	}).Debug("Starting message loop")

	b.Publish(tps)
	JSONResp(w, map[string]string{"status": "started", "tps": strconv.Itoa(tps)}, http.StatusOK)
}

func stopLoop(w http.ResponseWriter, r *http.Request) {
	log.WithFields(log.Fields{
		"method": r.Method,
		"path":   r.URL.Path,
	}).Debug("Received stopLoop request")

	b.Publish(0)
	handler.Stop(context.Background())

	log.Debug("Message loop stopped")
	JSONResp(w, map[string]string{"status": "stopped"}, http.StatusOK)
}

// Send Json in http response
func JSONResp(w http.ResponseWriter, resp interface{}, code int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.WithError(err).Error("Failed to encode JSON response")
	}
}

func printMetrics(inm *gometrics.InmemSink) {
	log.Debug("Starting metrics printer")
	ticker := time.NewTicker(5 * time.Second)
	for {
		<-ticker.C
		data := inm.Data()
		var interval *gometrics.IntervalMetrics
		n := len(data)
		switch {
		case n == 0:
			log.Debug("No metrics data available")
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

		log.Info(result)
	}
}
