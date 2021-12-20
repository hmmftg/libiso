package iso8583

import (
	"encoding/hex"
	"errors"
	_ "fmt"
	"log"

	"github.com/hmmftg/libiso/encoding/ebcdic"
	"github.com/hmmftg/libiso/encoding/iransystem"
)

type FieldData struct {
	fieldData []byte
	fieldDef  IsoField
	bmpDef    *BitMap
}

func (fldData *FieldData) BitmapDef() *BitMap {
	return fldData.bmpDef
}

func (fldData *FieldData) Def() IsoField {
	return fldData.fieldDef
}

var ErrUnsupportedEncoding = errors.New("libiso: Invalid/Unsupported field encoding")

//SetData sets field data as per the encoding
//additional padding will be applied if required
func (fldData *FieldData) SetData(value string) error {

	switch fldData.fieldDef.getDataEncoding() {
	case asciiEncoding:
		{
			switch fldData.fieldDef.(type) {
			case *FixedFieldDef:
				{
					data := []byte(value)
					fldData.setTruncatePad(data)
					break
				}
			default:
				{
					fldData.fieldData = []byte(value)
				}
			}

		}
	case iransystemEncoding:
		{
			data := iransystem.Decode(value)
			switch fldData.fieldDef.(type) {
			case *FixedFieldDef:
				{
					fldData.setTruncatePad(data)
					break
				}
			default:
				{
					fldData.fieldData = data
				}
			}

		}
	case ebcdicEncoding:
		{
			data := ebcdic.Decode(value)
			switch fldData.fieldDef.(type) {
			case *FixedFieldDef:
				{
					fldData.setTruncatePad(data)
					break
				}
			default:
				{
					fldData.fieldData = data
				}
			}

		}
	case binaryEncoding:
		fallthrough
	case bcdEncoding:
		{
			var err error

			data, err := hex.DecodeString(value)
			if err != nil {
				return errors.New("libiso: Invalid value for a binary/bcd field")
			}
			switch fldData.fieldDef.(type) {
			case *FixedFieldDef:
				{
					fldData.setTruncatePad(data)
					break
				}
			default:
				{
					fldData.fieldData = data
				}
			}

		}
	default:
		return ErrUnsupportedEncoding
	}

	return nil

}

func (fldData *FieldData) setTruncatePad(data []byte) {

	defObj := fldData.fieldDef.(*FixedFieldDef)
	padByte := byte(0x00)
	switch defObj.getDataEncoding() {
	case asciiEncoding:
	case iransystemEncoding:
		padByte = 0x20
	case ebcdicEncoding:
		padByte = 0x40
	}

	if len(data) == defObj.dataSize {
		fldData.fieldData = data
	} else if len(data) > defObj.dataSize {
		//truncate
		fldData.fieldData = data[:]
	} else {

		fldData.fieldData = make([]byte, defObj.dataSize)
		for i := range fldData.fieldData {
			fldData.fieldData[i] = padByte
		}
		copy(fldData.fieldData, data)
	}
}

//make a copy of FieldData
func (fldData *FieldData) copy() *FieldData {

	newFldData := new(FieldData)
	newFldData.fieldData = make([]byte, len(fldData.fieldData))
	copy(newFldData.fieldData, fldData.fieldData)
	newFldData.fieldDef = fldData.fieldDef

	return newFldData
}

func (fldData FieldData) String() string {

	if fldData.bmpDef != nil {
		return hex.EncodeToString(fldData.bmpDef.Bytes())
	}

	switch fldData.fieldDef.getDataEncoding() {
	case asciiEncoding:
		return string(fldData.fieldData)
	case ebcdicEncoding:
		encoded := ebcdic.EncodeToString(fldData.fieldData)
		log.Println("encoded - ", encoded, "hex ", hex.EncodeToString(fldData.fieldData))
		return encoded
	case iransystemEncoding:
		encoded := iransystem.EncodeToString(fldData.fieldData)
		log.Println("encoded - ", encoded, "hex ", hex.EncodeToString(fldData.fieldData))
		return encoded
	case binaryEncoding:
		fallthrough
	case bcdEncoding:
		return hex.EncodeToString(fldData.fieldData)
	default:
		log.Printf("unsupported encoding: %d", fldData.fieldDef.getDataEncoding())
	}

	return ""

}

//return the raw data associated with this field
//it will also include any ll portions for a variable field
func (fldData FieldData) Bytes() []byte {

	if fldData.bmpDef != nil {
		//if it's a bmp field, just return the data
		return fldData.bmpDef.Bytes()
	}

	if fldData.fieldDef.IsFixed() {
		dataLen := fldData.fieldDef.DataLength()
		if len(fldData.fieldData) > dataLen {
			log.Printf("Warning: field [%s] length exceeds defined length, will be truncated", fldData.Def().String())
			return fldData.fieldData[0:dataLen]
		} else if len(fldData.fieldData) < dataLen {
			//add default padding
			newFldData := make([]byte, dataLen)
			copy(newFldData, fldData.fieldData)
			return newFldData
		}
		return fldData.fieldData[0:dataLen]

	} else {
		//variable fields should have length indicators
		dataLen := len(fldData.fieldData)
		ll := fldData.fieldDef.EncodedLength(dataLen)
		llData := append(ll, fldData.fieldData...)
		return llData

	}

}
