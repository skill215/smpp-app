package smppapp

import (
	"context"

	"github.com/skill215/go-smpp/smpp"
)

type SmppTransmiter struct {
	conf *Smpp
	tx   *smpp.Transmitter
}

func ProvideSmppTransmitter(ctx context.Context, conf *Smpp) (*SmppTransmiter, error) {
	st := SmppTransmiter{}
	return &st, nil
}

func (st *SmppTransmiter) Init() {}

func (st *SmppTransmiter) Start() {

}

func (st *SmppTransmiter) Stop() {}
