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
	"github.com/skill215/smpp-app/limiter"
	msggenerator "github.com/skill215/smpp-app/msg-generator"
)

type SmppTransceiver struct {
	log          *logrus.Logger
	conf         *config.SmppConfig
	tr           []chan interface{}
	inm          *gometrics.InmemSink
	broker       *broker.Broker
	msgGenerator *msggenerator.MsgGenerator
}

func ProvideSmppTransceiver(ctx context.Context, conf config.SmppConfig, inm *gometrics.InmemSink, broker *broker.Broker, log *logrus.Logger) *SmppTransceiver {
	tr := SmppTransceiver{
		log:          log,
		conf:         &conf,
		inm:          inm,
		broker:       broker,
		tr:           []chan interface{}{},
		msgGenerator: msggenerator.New(&conf.Message),
	}
	return &tr
}

func (st *SmppTransceiver) Init() {
	st.log.Infof("transceiver init conf %+v", st.conf)
	for i := 0; i < int(st.conf.Client.Count); i++ {
		tr := &smpp.Transceiver{
			Addr:   fmt.Sprintf("%s:%d", st.conf.Server.Addr, st.conf.Server.Port),
			User:   st.conf.Server.User,
			Passwd: st.conf.Server.Password,
		}

		msgCh := st.broker.Subscribe()
		st.tr = append(st.tr, msgCh)
		st.bind(tr, msgCh)
	}
}

func (st *SmppTransceiver) bind(tc *smpp.Transceiver, msgCh chan interface{}) {
	conn := tc.Bind()
	tc.Handler = st.handleAT
	limiter := limiter.Limiter{}
	limiter.Set(0, time.Second)

	// goroutine to reconnect
	go func() {
		for {
			status := <-conn
			if status.Error() != nil || status.Status().String() != "Connected" {
				time.Sleep(5 * time.Second)
				conn = tc.Bind()
			}
		}
	}()

	// go routine to handle traffic control
	go func() {
		for {
			msg := <-msgCh
			tps := msg.(int)
			// every second allow tps, token bucket contains 1

			limiter.Set(tps, time.Second)
		}
	}()

	// goroutine to submit sm
	go func() {
		for {
			if limiter.Allow() {
				msg := st.msgGenerator.GenerateMsg()
				msg.Dst = st.msgGenerator.GenerateDaddr()
				// for USC2 encoding
				smlist, err := st.submitMsg(tc, msg)
				if err != nil {
					time.Sleep(50 * time.Microsecond)
				} else {
					for _, sm := range smlist {
						st.inm.IncrCounter([]string{"ao"}, 1)
						if sm.Resp().Header().Status != 0x00000000 {
							st.inm.IncrCounter([]string{"ao failure"}, 1)
						}
					}
				}
			} else {
				// not allowed in this second, just sleep
				time.Sleep(10 * time.Millisecond)
			}
		}
	}()

}

func (st *SmppTransceiver) Start(tps int) {

}

func (st *SmppTransceiver) Stop() {

}

func (st *SmppTransceiver) handleAT(p pdu.Body) {
	st.log.Debugf("receive AT, ID: %s, Status: %s", p.Header().ID.String(), p.Header().Status.Error())
	if p.Header().Status != 0x00000000 {
		st.inm.IncrCounter([]string{"at failure"}, 1)
	}
	st.inm.IncrCounter([]string{"at"}, 1)
}

func (st *SmppTransceiver) submitMsg(tc *smpp.Transceiver, msg *smpp.ShortMessage) ([]smpp.ShortMessage, error) {
	if len(msg.Text.Encode()) <= 132 {
		if sm, err := tc.Submit(msg); err != nil {
			return []smpp.ShortMessage{}, err
		} else {
			return []smpp.ShortMessage{*sm}, nil
		}
	} else {
		// concatenated message
		return tc.SubmitLongMsg(msg)
	}
}
