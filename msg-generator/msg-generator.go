package msggenerator

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/skill215/go-smpp/smpp"
	"github.com/skill215/go-smpp/smpp/pdu/pdufield"
	"github.com/skill215/go-smpp/smpp/pdu/pdutext"
	"github.com/skill215/smpp-app/config"
)

type MsgGenerator struct {
	sync.Mutex
	index        int
	conf         *config.MessageConfig
	stop         int
	textContents []string
	urlContents  []string
	useRandom    bool
	rnd          *rand.Rand
}

func New(conf *config.MessageConfig) *MsgGenerator {
	stop := conf.Send.Dst.Daddr.Start
	if stop <= conf.Send.Dst.Daddr.Start {
		stop = int(math.Pow(10, float64(conf.Send.Dst.Daddr.GenerateLen))) - 1
	}

	mg := &MsgGenerator{
		conf:      conf,
		stop:      stop,
		useRandom: false,
		rnd:       rand.New(rand.NewSource(time.Now().UnixNano())),
	}

	// Load file contents
	mg.loadFileContents()

	return mg
}

func (mg *MsgGenerator) loadFileContents() {
	// Read text file
	if content, err := os.ReadFile(mg.conf.Send.TextFile); err != nil {
		log.Printf("Error reading text file: %v, will use random mode", err)
		mg.useRandom = true
	} else {
		// Split by lines and skip empty lines
		for _, line := range strings.Split(string(content), "\n") {
			if trimmed := strings.TrimSpace(line); trimmed != "" {
				mg.textContents = append(mg.textContents, trimmed)
			}
		}
	}

	// Read url file
	if content, err := os.ReadFile(mg.conf.Send.UrlFile); err != nil {
		log.Printf("Error reading url file: %v, will use random mode", err)
		mg.useRandom = true
	} else {
		// Split by lines and skip empty lines
		for _, line := range strings.Split(string(content), "\n") {
			if trimmed := strings.TrimSpace(line); trimmed != "" {
				mg.urlContents = append(mg.urlContents, trimmed)
			}
		}
	}
}

func (mg *MsgGenerator) getPreDefinedContent() string {
	if len(mg.textContents) == 0 || len(mg.urlContents) == 0 {
		// If either file content is empty, return default content
		return "default message content"
	}

	// Randomly select text and URL
	text := mg.textContents[mg.rnd.Intn(len(mg.textContents))]
	url := mg.urlContents[mg.rnd.Intn(len(mg.urlContents))]

	return text + " " + url
}

func (mg *MsgGenerator) GenerateMsgContent(ud string) string {
	// If forced to use random mode or configured as random mode
	if mg.useRandom || mg.conf.Send.ContentMode == "random" {
		if strings.Contains(ud, "{random url}") {
			return strings.Replace(ud, "{random url}", generateRandomURL(), 1)
		}
		return ud
	}

	// If pre-defined mode, directly use file content
	if mg.conf.Send.ContentMode == "pre-defined" {
		return mg.getPreDefinedContent()
	}

	// Mixed mode
	if mg.conf.Send.ContentMode == "mixed" {
		// Generate random number between 0-1
		if mg.rnd.Float64() < mg.conf.Send.PreDefinedContentRatio {
			return mg.getPreDefinedContent()
		}
		// Use random mode
		if strings.Contains(ud, "{random url}") {
			return strings.Replace(ud, "{random url}", generateRandomURL(), 1)
		}
		return ud
	}

	// Unknown mode, use random
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
