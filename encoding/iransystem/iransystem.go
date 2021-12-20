package iransystem

import (
	"bytes"
	"io/ioutil"
	"strconv"

	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"
)

var irsysToAscii = "000102030405060708090A0B0C0D0E0F" +
	"101112131415161718191A1B1C1D1E1F" +
	"FF2122232425262728292A2B2C2D2E2F" +
	"303132333435363738393A3B3C3D3E3F" +
	"404142434445464748494A4B4C4D4E4F" +
	"505152535455565758595A5B5C5D5E5F" +
	"606162636465666768696A6B6C6D6E6F" +
	"707172737475767778797A7B7C7D7E7F" +
	"30313233343536373839A1DCBFC2C6C1" +
	"C7C7C8C88181CACACBCBCCCC8D8DCDCD" +
	"CECECFD0D1D28ED3D3D4D4D5D5D6D6D8" +
	"D9B1B2B3B4B5B6B7B8B9BABBBCBDBEBF" +
	"7B414243444546474849CACBCCCDCECF" +
	"7D4A4B4C4D4E4F505152DADBDCDDDEDF" +
	"5CDADADADADBDBDBDBDDDDDEDE989890" +
	"90E1E1E1E3E3E4E4E6E5E5E5ECECECA0"

var asciiToIrsys = "000102030405060708090A0B0C0D0E0F" +
	"101112131415161718191A1B1C1D1E1F" +
	"202122232425262728292A2B2C2D2E2F" +
	"303132333435363738393A3B3C3D3E3F" +
	"404142434445464748494A4B4C4D4E4F" +
	"505152535455565758595A5B5C5D5E5F" +
	"606162636465666768696A6B6C6D6E6F" +
	"707172737475767778797A7B7C7D7E7F" +
	"80948283848586878889A1DC8C9CA68F" +
	"EF91929394959697ED999A9B9C9D9E9F" +
	"FF8AA2A3A4A5A6A7A8A9AAABACADAEAF" +
	"B0B1B2B3B4B5B6B7B8B9BABBBCBDBE8C" +
	"C08F8DC3C4C58E9092C996989A9EA0A2" +
	"A3A4A5A7A9ABADD7AFB0E1E58BE9EBDF" +
	"E0F1E2F4F6F9F8E7E8E9EAEBFDEDEEEF" +
	"F0F1F2F3F4F5F6F7F8F9FAFBFCFDFE20"

//EncodeToString converts from iransystem bytes to a utf8 string
func EncodeToString(data []byte) string {

	buf := bytes.NewBuffer([]byte{})

	for _, b := range data {
		var x uint32 = uint32(b)
		tmp := irsysToAscii[(x * 2) : (x*2)+2]
		i, _ := strconv.ParseUint(tmp, 16, 8)
		buf.WriteByte(byte(i))
	}
	rData := buf.Bytes()
	for i, j := 0, len(rData)-1; i < j; i, j = i+1, j-1 {
		rData[i], rData[j] = rData[j], rData[i]
	}

	input := bytes.NewReader(rData)
	csentry := charmap.Windows1256
	reader := transform.NewReader(input, csentry.NewDecoder())
	output, _ := ioutil.ReadAll(reader)

	return string(output)
}

//Decode converts from a utf8 encoded string to iransystem bytes
func Decode(str string) []byte {
	csentry := charmap.Windows1256

	win1256 := make([]byte, 0)
	for _, r := range str {
		b, _ := csentry.EncodeRune(r)
		win1256 = append(win1256, b)
	}
	data := make([]byte, 0)

	for i := 0; i < len(win1256); i++ {
		b := uint32(win1256[i])
		bVal, _ := strconv.ParseUint(asciiToIrsys[b*2:b*2+2], 16, 8)
		data = append(data, byte(bVal))
	}

	for i, j := 0, len(data)-1; i < j; i, j = i+1, j-1 {
		data[i], data[j] = data[j], data[i]
	}

	return data
}
