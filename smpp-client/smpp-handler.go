package smppclient

import (
	"context"
	"time"

	gometrics "github.com/armon/go-metrics"
	"github.com/sirupsen/logrus"
	"github.com/skill215/smpp-app/config"
)

var (
	Interval = 5
)

type SmppHandler struct {
	log       *logrus.Logger
	conf      []config.SmppConfig
	inm       *gometrics.InmemSink
	sender    []SmppTransmiter // all sender for smpp
	receiver  []SmppReceiver   // all receiver for smpp
	eventChan chan int
}

func ProvideService(ctx context.Context, log *logrus.Logger, conf []config.SmppConfig, event chan int) (*SmppHandler, error) {
	handler := SmppHandler{
		log:       log,
		conf:      conf,
		eventChan: event,
	}

	inm := gometrics.NewInmemSink(time.Duration(Interval)*time.Second, time.Minute)
	// sig := metrics.DefaultInmemSignal(inm)
	gometrics.NewGlobal(gometrics.DefaultConfig("smpp-app"), inm)
	handler.inm = inm

	return &handler, nil
}

func (sh *SmppHandler) Init(ctx context.Context) error {
	return nil
}

func (sh *SmppHandler) Run(ctx context.Context) error {
	return nil
}
