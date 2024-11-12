package msggenerator

import (
	"testing"

	"github.com/go-playground/assert/v2"
	"github.com/skill215/go-smpp/smpp/pdu/pdutext"
)

func TestConvert7to8(t *testing.T) {
	in := []byte{0x1b, 0x61}  // 0x0001 1011 0110 0001
	exp := []byte{0x37, 0x84} // 0x0011 011 110001 00
	out := convert8bitTo7bit(in)

	in2 := []byte{0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f}
	exp2 := []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}
	out2 := convert8bitTo7bit(in2)
	assert.Equal(t, exp, out)
	assert.Equal(t, exp2, out2)

	in3 := []byte{0x1b, 0x61, 0x1b, 0x69, 0x1b, 0x6f, 0x1b, 0x75}
	out3 := convert8bitTo7bit(in3)
	t.Log(out3)
}

func TestConvertlatin1ToUCS2(t *testing.T) {
	latinStr := "¡¢£¤¥¦§¨©ª«¬®¯°±²³´µ¶·¸¹º»¼½¾¿ÀÁÂÃÄÅÆÇÉÊÈËÌÍÎÏÐÑÒÓÔ{Ö×ØÙÚÛÜÝÞßàáâãäåæçèéêëìíîïðñòóôõö÷øùúûüýþÿ"
	latin := pdutext.Latin1(latinStr)
	latinBytes := latin.Encode()
	ucs2 := pdutext.UCS2(latinStr)
	ucs2Bytes := ucs2.Encode()

	assert.Equal(t, len(ucs2Bytes), 2*len(latinBytes))
	for i := 0; i < len(latinBytes); i++ {
		assert.Equal(t, ucs2Bytes[i*2+1], latinBytes[i])
		assert.Equal(t, ucs2Bytes[i*2], byte(0))
	}
}
