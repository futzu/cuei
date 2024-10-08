package cuei

import (
	"fmt"
)

var uriUpids = map[uint8]string{
	0x01: "Deprecated",
	0x02: "Deprecated",
	0x03: "AdID",
	0x07: "TID",
	0x09: "ADI",
	0x10: "UUID",
	0x11: "ACR",
	0x0e: "ADS Info",
	0x0f: "URI",
}

/*
Upid is the Struct for Segmentation Upids

Non-standard UPID types are returned as bytes.
*/
type Upid struct {
	Name             string `json:",omitempty"`
	UpidType         uint8  `json:",omitempty"`
	Value            string `json:",omitempty"`
	TSID             uint16 `json:",omitempty"`
	Reserved         uint8  `json:",omitempty"`
	EndOfDay         uint8  `json:",omitempty"`
	UniqueFor        uint16 `json:",omitempty"`
	ContentID        []byte `json:",omitempty"`
	Upids            []Upid `json:",omitempty"`
	FormatIdentifier string `json:",omitempty"`
	PrivateData      []byte `json:",omitempty"`
}

// Decode Upids
func (upid *Upid) decode(bd *bitDecoder, upidType uint8, upidlen uint8) {

	upid.UpidType = upidType

	name, ok := uriUpids[upidType]
	if ok {
		upid.Name = name
		upid.uri(bd, upidlen)
	} else {

		switch upidType {
		case 0x05, 0x06:
			upid.Name = "ISAN"
			upid.isan(bd, upidlen)
		case 0x08:
			upid.Name = "AiringID"
			upid.airid(bd, upidlen)
		case 0x0a:
			upid.Name = "EIDR"
			upid.eidr(bd, upidlen)
		case 0x0b:
			upid.Name = "ATSC"
			upid.atsc(bd, upidlen)
		case 0x0c:
			upid.Name = "MPU"
			upid.mpu(bd, upidlen)
		case 0x0d:
			upid.Name = "MID"
			upid.mid(bd, upidlen)
		default:
			upid.Name = "UPID"
			upid.uri(bd, upidlen)
		}
	}
}

// Decode for AirId
func (upid *Upid) airid(bd *bitDecoder, upidlen uint8) {
	upid.Value = bd.asHex(uint(upidlen << 3))
}

// Decode for Isan Upid
func (upid *Upid) isan(bd *bitDecoder, upidlen uint8) {
	upid.Value = bd.asAscii(uint(upidlen << 3))
}

// Decode for URI Upid
func (upid *Upid) uri(bd *bitDecoder, upidlen uint8) {
	upid.Value = bd.asAscii(uint(upidlen) << 3)
}

// Decode for ATSC Upid
func (upid *Upid) atsc(bd *bitDecoder, upidlen uint8) {
	upid.TSID = bd.uInt16(16)
	upid.Reserved = bd.uInt8(2)
	upid.EndOfDay = bd.uInt8(5)
	upid.UniqueFor = bd.uInt16(9)
	upid.ContentID = bd.asBytes(uint((upidlen - 4) << 3))
}

// Decode for EIDR Upid
func (upid *Upid) eidr(bd *bitDecoder, upidlen uint8) {
	head := bd.uInt16(16)
	// Switching to Compact Binary Format
	var astring string
	nibbles := 20
	for i := 0; i < nibbles; i++ {
		astring = fmt.Sprintf("%v%x", astring, bd.uInt8(4))
	}
	upid.Value = fmt.Sprintf("0x%x%v", head, astring)
}

// Decode for MPU Upid
func (upid *Upid) mpu(bd *bitDecoder, upidlen uint8) {
	ulb := uint(upidlen) << 3
	upid.FormatIdentifier = bd.asHex(32)
	upid.PrivateData = bd.asBytes(ulb - 32)
}

// Decode for MID Upid
func (upid *Upid) mid(bd *bitDecoder, upidlen uint8) {
	var i uint8
	i = 0
	for i < upidlen {
		utype := bd.uInt8(8)
		i++
		ulen := bd.uInt8(8)
		i++
		i += ulen
		var mupid Upid
		upid.decode(bd, utype, ulen)
		upid.Upids = append(upid.Upids, mupid)
	}
}

// Encode Upids
func (upid *Upid) encode(be *bitEncoder, upidType uint8) {
	switch upid.UpidType {
	case 0x05, 0x06:
		upid.encodeIsan(be)
	case 0x08:
		upid.encodeAirId(be)
	case 0x0a:
		upid.encodeEidr(be)
	default:
		upid.encodeUri(be)
	}
}

// encode for Uri Upids
func (upid *Upid) encodeUri(be *bitEncoder) {
	if len(upid.Value) > 0 {
		be.AddBytes([]byte(upid.Value), uint(len(upid.Value)<<3))
	}
}

// encode for AirId
func (upid *Upid) encodeAirId(be *bitEncoder) {
	if len(upid.Value) > 0 {
		be.AddHex64(upid.Value, uint(len(upid.Value)<<3))
	}
}

// encode for Isan Upid
func (upid *Upid) encodeIsan(be *bitEncoder) {
	if len(upid.Value) > 0 {
		be.AddBytes([]byte(upid.Value), uint(len(upid.Value)<<3))
	}
}

// encode for Eidr Upid
func (upid *Upid) encodeEidr(be *bitEncoder) {
	be.AddHex64(upid.Value[:6], 16)
	substring := upid.Value[6:]
	for _, c := range substring {
		hexed := fmt.Sprintf("0x%s", string(c))
		be.AddHex64(hexed, 4)
	}
}
