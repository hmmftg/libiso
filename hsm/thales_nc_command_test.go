package hsm

import (
	"encoding/hex"
	"fmt"
	"github.com/hmmftg/libiso/net"
	"strings"
	"testing"
)

func init() {
	hsm := NewThalesHsm("127.0.0.1", 1500, AsciiEncoding)
	go hsm.Start()

}

func Test_Thales_NC(t *testing.T) {

	cmdStr := "303030303030303030303032;4e43;"
	cmdStr = strings.Replace(cmdStr, ";", "", -1)
	msgData, _ := hex.DecodeString(cmdStr)

	fmt.Println(hex.Dump(msgData))

	hsmClient := net.NewNetCatClient("127.0.0.1:1500", net.Mli2e)
	err := hsmClient.OpenConnection()
	failOnErr(t, err)
	defer hsmClient.Close()
	err = hsmClient.Write(msgData)

	failOnErr(t, err)

	responseData, err := hsmClient.ReadNextPacket()
	failOnErr(t, err)
	fmt.Println(hex.Dump(responseData))

}
