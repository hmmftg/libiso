package iransystem

import (
	"bytes"
	"encoding/hex"
	"testing"

	"github.com/hmmftg/libiso/encoding/iransystem"
)

func Test_IranSystem(t *testing.T) {

	//t.Log(ebcdic_to_ascii)
	buf := bytes.NewBuffer([]byte{})
	buf.WriteByte(0x2)
	buf.Bytes()

	data, _ := hex.DecodeString("f0f1f2f3f42020202090919293949596")
	str := iransystem.EncodeToString(data)
	t.Log(str, "\n")

	data = iransystem.Decode("حمید سلام")
	t.Log(hex.EncodeToString(data), "\n")

	fromBytes := iransystem.EncodeToString([]byte{0xF0, 0xF1, 0xF0, 0xF0})
	t.Log(fromBytes)

	fromBytes = iransystem.EncodeToString([]byte("حمید"))
	t.Log(fromBytes)
	t.Log(hex.EncodeToString([]byte("AGNS")))

}
