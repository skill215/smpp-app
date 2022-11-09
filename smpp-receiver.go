package smppapp

import (
	"context"
	"fmt"
	"log"
	"sync/atomic"

	"github.com/skill215/go-smpp/smpp"
	"github.com/skill215/go-smpp/smpp/pdu"
)

var total atomic.Uint64

type SmppReceiver struct {
	conf *Smpp
	rc   *smpp.Receiver
}

func ProvideSmppReceiver(ctx context.Context, conf *Smpp) (*SmppReceiver, error) {
	sr := SmppReceiver{
		conf: conf,
		rc: &smpp.Receiver{
			Addr:   fmt.Sprintf("%s:%d", conf.Server.Addr, conf.Server.Port),
			User:   conf.Server.User,
			Passwd: conf.Server.Password,
		},
	}
	return &sr, nil
}

func (sr *SmppReceiver) bind() error {
	sr.rc.Handler = handleAt
	conn := sr.rc.Bind()
	if status := <-conn; status.Error() != nil {
		log.Fatalf("unable to connect to smpp server. err %v", status.Error())
	}
	return nil
}

func handleAt(p pdu.Body) {
	total.Add(1)
}

func (sr *SmppReceiver) Init() {}

func (sr *SmppReceiver) Start() {

}

func (sr *SmppReceiver) Stop() {

}
