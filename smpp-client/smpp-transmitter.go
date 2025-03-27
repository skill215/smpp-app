package smppclient

import (
	"context"
	"fmt"
	"net"
	"time"

	gometrics "github.com/armon/go-metrics"
	"github.com/sirupsen/logrus"
	"github.com/skill215/go-smpp/smpp"
	"github.com/skill215/smpp-app/broker"
	"github.com/skill215/smpp-app/config"
	"github.com/skill215/smpp-app/limiter"
	msggenerator "github.com/skill215/smpp-app/msg-generator"
)

type SmppTransmiter struct {
	log          *logrus.Logger
	conf         *config.SmppConfig
	tx           []chan interface{}
	inm          *gometrics.InmemSink
	broker       *broker.Broker
	msgGenerator *msggenerator.MsgGenerator
}

func ProvideSmppTransmitter(ctx context.Context, conf config.SmppConfig, inm *gometrics.InmemSink, broker *broker.Broker, log *logrus.Logger) *SmppTransmiter {
	st := SmppTransmiter{
		log:          log,
		conf:         &conf,
		inm:          inm,
		broker:       broker,
		tx:           []chan interface{}{},
		msgGenerator: msggenerator.New(&conf.Message),
	}
	return &st
}

func (st *SmppTransmiter) Init() {
	st.log.Infof("transmitter init %+v", st.conf)
	for i := 0; i < int(st.conf.Client.Count); i++ {
		tx := smpp.Transmitter{
			Addr:   fmt.Sprintf("%s:%d", st.conf.Server.Addr, st.conf.Server.Port),
			User:   st.conf.Server.User,
			Passwd: st.conf.Server.Password,
		}

		msgCh := st.broker.Subscribe()
		st.tx = append(st.tx, msgCh)
		st.bind(&tx, msgCh)
	}
}

func (st *SmppTransmiter) bind(tx *smpp.Transmitter, msgCh chan interface{}) {
	conn := tx.Bind()
	st.log.WithFields(logrus.Fields{
		"addr":     tx.Addr,
		"user":     tx.User,
		"type":     st.conf.Client.Type,
		"conn_num": st.conf.Client.Count,
	}).Info("Starting SMPP bind")

	limiter := limiter.Limiter{}
	limiter.Set(0, time.Second)

	// goroutine to reconnect
	go func() {
		var lastStatus string
		for {
			status := <-conn
			currentStatus := status.Status().String()

			if status.Error() != nil {
				st.log.WithFields(logrus.Fields{
					"addr":       tx.Addr,
					"user":       tx.User,
					"type":       st.conf.Client.Type,
					"error":      status.Error(),
					"error_type": fmt.Sprintf("%T", status.Error()),
					"status":     currentStatus,
					"raw_status": fmt.Sprintf("%+v", status),
				}).Error("SMPP bind failed")

				// Add network diagnostic logs
				if netErr, ok := status.Error().(*net.OpError); ok {
					st.log.WithFields(logrus.Fields{
						"network":     netErr.Net,
						"source":      netErr.Source,
						"address":     netErr.Addr,
						"timeout":     netErr.Timeout(),
						"temporary":   netErr.Temporary(),
						"error_phase": netErr.Op,
					}).Error("Network operation error details")
				}

				time.Sleep(5 * time.Second)
				st.log.WithFields(logrus.Fields{
					"addr":    tx.Addr,
					"user":    tx.User,
					"attempt": "reconnect",
				}).Debug("Attempting to rebind...")
				conn = tx.Bind()
			} else if currentStatus != "Connected" {
				st.log.WithFields(logrus.Fields{
					"addr":       tx.Addr,
					"user":       tx.User,
					"type":       st.conf.Client.Type,
					"status":     currentStatus,
					"prev_error": status.Error(),
					"raw_status": fmt.Sprintf("%+v", status),
				}).Warn("SMPP connection status changed")
				time.Sleep(5 * time.Second)
				st.log.WithFields(logrus.Fields{
					"addr":    tx.Addr,
					"user":    tx.User,
					"attempt": "reconnect",
				}).Debug("Attempting to rebind...")
				conn = tx.Bind()
			} else if lastStatus != "Connected" {
				// Only print once when transitioning from non-Connected to Connected status
				st.log.WithFields(logrus.Fields{
					"addr":   tx.Addr,
					"user":   tx.User,
					"type":   st.conf.Client.Type,
					"status": currentStatus,
				}).Info("SMPP bind successful")
			}
			lastStatus = currentStatus
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
				// Generate a new message each time before sending
				msg := st.msgGenerator.GenerateMsg()
				msg.Dst = st.msgGenerator.GenerateDaddr()
				// for USC2 encoding
				smlist, err := st.submitMsg(tx, msg)
				if err != nil {
					st.log.WithFields(logrus.Fields{
						"addr":           tx.Addr,
						"user":           tx.User,
						"dst":            msg.Dst,
						"error":          err,
						"content_length": len(msg.Text.Encode()),
					}).Debug("Failed to submit message")
					time.Sleep(50 * time.Microsecond)
				} else {
					for _, sm := range smlist {
						st.inm.IncrCounter([]string{"ao"}, 1)
						if sm.Resp().Header().Status != 0x00000000 {
							st.log.WithFields(logrus.Fields{
								"addr":   tx.Addr,
								"user":   tx.User,
								"dst":    msg.Dst,
								"status": sm.Resp().Header().Status,
							}).Debug("Message submission got non-zero status")
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

func (st *SmppTransmiter) Start(tps int) {

}

func (st *SmppTransmiter) Stop() {

}

func (st *SmppTransmiter) submitMsg(tx *smpp.Transmitter, msg *smpp.ShortMessage) ([]smpp.ShortMessage, error) {
	if len(msg.Text.Encode()) <= 132 {
		if sm, err := tx.Submit(msg); err != nil {
			return []smpp.ShortMessage{}, err
		} else {
			return []smpp.ShortMessage{*sm}, nil
		}
	} else {
		// concatenated message
		return tx.SubmitLongMsg(msg)
	}
}
