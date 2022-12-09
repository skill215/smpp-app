package msggenerator

import (
	"fmt"
	"math"
	"sync"

	"github.com/skill215/go-smpp/smpp"
	"github.com/skill215/go-smpp/smpp/pdu/pdufield"
	"github.com/skill215/go-smpp/smpp/pdu/pdutext"
	"github.com/skill215/smpp-app/config"
)

type MsgGenerator struct {
	sync.Mutex
	index int
	conf  *config.MessageConfig
	stop  int
}

func New(conf *config.MessageConfig) *MsgGenerator {
	stop := conf.Send.Dst.Daddr.Start
	if stop <= conf.Send.Dst.Daddr.Start {
		stop = int(math.Pow(10, float64(conf.Send.Dst.Daddr.GenerateLen))) - 1
	}
	return &MsgGenerator{
		conf: conf,
		stop: stop,
	}
}

func (mg *MsgGenerator) GenerateMsg() *smpp.ShortMessage {
	sms := smpp.ShortMessage{
		//Dst:           mg.generateDaddr(),
		SourceAddrTON: uint8(mg.conf.Send.Src.Ton),
		SourceAddrNPI: uint8(mg.conf.Send.Src.Npi),
		DestAddrTON:   uint8(mg.conf.Send.Dst.Ton),
		DestAddrNPI:   uint8(mg.conf.Send.Dst.Npi),
		Text:          pdutext.UCS2(mg.conf.Send.Content),
	}
	if len(mg.conf.Send.Src.Oaddr) > 0 {
		sms.Src = mg.conf.Send.Src.Oaddr
	}
	if mg.conf.Send.RequireSR {
		sms.Register = pdufield.FinalDeliveryReceipt
	} else {
		sms.Register = pdufield.NoDeliveryReceipt
	}
	return &sms
}

func (mg *MsgGenerator) GenerateDaddr() string {
	mg.Lock()
	defer mg.Unlock()
	if mg.index >= mg.stop {
		mg.index = mg.conf.Send.Dst.Daddr.Start - 1
	}
	mg.index++
	daddr := fmt.Sprintf("%s%0*d%s", mg.conf.Send.Dst.Daddr.Prefix, mg.conf.Send.Dst.Daddr.GenerateLen, mg.index, mg.conf.Send.Dst.Daddr.Suffix)
	//fmt.Println("daddr is ", daddr)
	return daddr
}
