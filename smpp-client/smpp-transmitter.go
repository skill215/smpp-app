package smppclient

import (
	"context"

	"github.com/skill215/go-smpp/smpp"
	"github.com/skill215/smpp-app/config"
)

type SmppTransmiter struct {
	conf *config.SmppConfig
	tx   *smpp.Transmitter
}

func ProvideSmppTransmitter(ctx context.Context, conf *config.SmppConfig) (*SmppTransmiter, error) {
	st := SmppTransmiter{}
	return &st, nil
}

func (st *SmppTransmiter) Init() {}

func (st *SmppTransmiter) Start() {

}

func (st *SmppTransmiter) Stop() {}
