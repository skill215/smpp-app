package smppclient

import (
	"context"

	gometrics "github.com/armon/go-metrics"
	"github.com/sirupsen/logrus"
	"github.com/skill215/smpp-app/broker"
	"github.com/skill215/smpp-app/config"
)

var (
	Interval = 5
)

type SmppClient interface {
	Init()
	Start(int)
	Stop()
}

type SmppHandler struct {
	log     *logrus.Logger
	inm     *gometrics.InmemSink
	broker  *broker.Broker
	clients []SmppClient
}

func ProvideService(ctx context.Context, log *logrus.Logger, conf []config.SmppConfig, broker *broker.Broker, inm *gometrics.InmemSink) *SmppHandler {
	handler := SmppHandler{
		log:     log,
		broker:  broker,
		inm:     inm,
		clients: []SmppClient{},
	}

	for _, c := range conf {
		handler.clients = append(handler.clients, createClient(c, log, handler.inm, broker))
	}

	log.Infof("inital %d clinets\n", len(handler.clients))
	return &handler
}

func (sh *SmppHandler) Init(ctx context.Context) {
	for _, client := range sh.clients {
		client.Init()
	}
}

func (sh *SmppHandler) Run(ctx context.Context, tps int) {
	for _, client := range sh.clients {
		client.Start(tps)
	}
}

func (sh *SmppHandler) Stop(ctx context.Context) {
	for _, client := range sh.clients {
		client.Stop()
	}
}

func createClient(conf config.SmppConfig, log *logrus.Logger, inm *gometrics.InmemSink, broker *broker.Broker) SmppClient {
	ctx := context.Background()
	log.Infof("create client with conf %+v", conf)
	switch conf.Client.Type {
	case "transceiver":
		return ProvideSmppTransceiver(ctx, conf, inm, broker, log)
	case "receiver":
		return ProvideSmppReceiver(ctx, conf, inm, broker, log)
	default:
		return ProvideSmppTransmitter(ctx, conf, inm, broker, log)
	}
}
