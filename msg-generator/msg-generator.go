package msggenerator

import (
	"fmt"
	"math"
	"strings"
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

func (mg *MsgGenerator) GenerateMsgContent(ud string) string {
	if strings.Contains(ud, "{random url}") {
		return strings.Replace(ud, "{random url}", generateRandomURL(), 1)
	}
	return ud
}

func (mg *MsgGenerator) GenerateMsg() *smpp.ShortMessage {
	sms := smpp.ShortMessage{
		//Dst:           mg.generateDaddr(),
		SourceAddrTON: uint8(mg.conf.Send.Src.Ton),
		SourceAddrNPI: uint8(mg.conf.Send.Src.Npi),
		DestAddrTON:   uint8(mg.conf.Send.Dst.Ton),
		DestAddrNPI:   uint8(mg.conf.Send.Dst.Npi),
		//Text:          pdutext.Latin1(mg.conf.Send.Content),
		//Text: pdutext.Raw([]byte(mg.conf.Send.Content)),
		//ESMClass: 0x40,
	}
	content := mg.GenerateMsgContent(mg.conf.Send.Content)
	switch mg.conf.Send.Dcs {
	case 0:
		sms.Text = pdutext.GSM7(content)
	case 3:
		sms.Text = pdutext.Latin1(content)
	case 4:
		sms.Text = pdutext.Binary2(content)
	case 8:
		sms.Text = pdutext.UCS2(content)
	default:
		sms.Text = pdutext.Raw(content)
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

func convert8bitTo7bit(in []byte) []byte {
	if len(in) == 0 {
		return nil
	}
	if len(in) == 1 {
		return []byte{in[0] << 1}
	}
	out := []byte{}
	shift := 1
	for i := 0; i < len(in)-1; i++ {
		ou := (in[i]&0x7f)<<shift | (in[i+1]&0x7f)>>(7-shift)
		out = append(out, ou)
		shift++
		if shift == 8 {
			i++
			shift = 1
		}
	}
	if shift != 1 {
		out = append(out, in[len(in)-1]<<shift)
	}
	return out
}
