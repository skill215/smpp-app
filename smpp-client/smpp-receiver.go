package smppclient

import (
	"context"
	"fmt"
	"time"

	gometrics "github.com/armon/go-metrics"
	"github.com/sirupsen/logrus"
	"github.com/skill215/go-smpp/smpp"
	"github.com/skill215/go-smpp/smpp/pdu"
	"github.com/skill215/smpp-app/broker"
	"github.com/skill215/smpp-app/config"
)

type SmppReceiver struct {
	log    *logrus.Logger
	conf   *config.SmppConfig
	rc     []smpp.Receiver
	inm    *gometrics.InmemSink
	broker *broker.Broker
}

func ProvideSmppReceiver(ctx context.Context, conf config.SmppConfig, inm *gometrics.InmemSink, broker *broker.Broker, log *logrus.Logger) *SmppReceiver {
	sr := SmppReceiver{
		conf:   &conf,
		inm:    inm,
		log:    log,
		broker: broker,
		rc:     []smpp.Receiver{},
	}
	return &sr
}

func (sr *SmppReceiver) Init() {
	sr.log.Infof("smpp receiver init")
	for i := 0; i < int(sr.conf.Client.Count); i++ {
		rc := smpp.Receiver{
			Addr:   fmt.Sprintf("%s:%d", sr.conf.Server.Addr, sr.conf.Server.Port),
			User:   sr.conf.Server.User,
			Passwd: sr.conf.Server.Password,
		}
		sr.rc = append(sr.rc, rc)
		sr.bind(&rc)
	}
}

func (sr *SmppReceiver) bind(rc *smpp.Receiver) {
	conn := rc.Bind()
	rc.Handler = sr.handleAT

	// goroutine to reconnect
	go func() {
		for {
			status := <-conn
			if status.Error() != nil || status.Status().String() != "Connected" {
				time.Sleep(5 * time.Second)
				conn = rc.Bind()
			}
		}
	}()
}

func (sr *SmppReceiver) Start(tps int) {

}

func (sr *SmppReceiver) Stop() {

}

func (sr *SmppReceiver) handleAT(p pdu.Body) {
	sr.log.Debugf("receive AT, ID: %s, Status: %s", p.Header().ID.String(), p.Header().Status.Error())
	sr.inm.IncrCounter([]string{"at"}, 1)
	if p.Header().Status != 0x00000000 {
		sr.inm.IncrCounter([]string{"at failure"}, 1)
	}
}
