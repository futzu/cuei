package cuei

import (
	"fmt"
	gobs "github.com/futzu/gob"
)

/*
*
Upid is the Struct for Segmentation Upids

	    These UPID types are recognized.

            0x01: "Deprecated", 
            0x02: "Deprecated",
            0x03: "AdID",
            0x05: "ISAN"
            0x06: "ISAN"
            0x07: "TID",
            0x08: "AiringID",
            0x09: "ADI",
            0x10: "UUID",
            0x11: "ACR",
            0x0a: "EIDR",
            0x0b: "ATSC",
            0x0c: "MPU",
            0x0d: "MID",
            0x0e: "ADS Info",
            0x0f: "URI",

	    Non-standard UPID types are returned as bytes.

*
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

// Decode Decodes Segmentation UPIDs
func (upid *Upid) Decode(gob *gobs.Gob, upidType uint8, upidlen uint8) {

	upid.UpidType = upidType

	var uri_upids = map[uint8]string{
		0x01: "Deprecated",
		0x02: "Deprecated",
		0x03: "AdID",
		0x07: "TID",
		0x08: "AiringID",
		0x09: "ADI",
		0x10: "UUID",
		0x11: "ACR",
		0x0a: "EIDR",
		0x0b: "ATSC",
		0x0c: "MPU",
		0x0d: "MID",
		0x0e: "ADS Info",
		0x0f: "URI",
	}

	name, ok := uri_upids[upidType]
	if ok {
		upid.Name = name
		upid.uri(gob, upidlen)
	} else {

		switch upidType {

		case 0x05, 0x06:
			upid.Name = "ISAN"
			upid.isan(gob, upidlen)
		case 0x08:
			upid.Name = "AiringID"
			upid.airid(gob, upidlen)
		case 0x0a:
			upid.Name = "EIDR"
			upid.eidr(gob, upidlen)
		case 0x0b:
			upid.Name = "ATSC"
			upid.atsc(gob, upidlen)
		case 0x0c:
			upid.Name = "MPU"
			upid.mpu(gob, upidlen)
		case 0x0d:
			upid.Name = "MID"
			upid.mid(gob, upidlen)
		default:
			upid.Name = "UPID"
			upid.uri(gob, upidlen)
		}
	}
}

// Decode for AirId
func (upid *Upid) airid(gob *gobs.Gob, upidlen uint8) {
	upid.Value = gob.Hex(uint(upidlen << 3))
}

// Decode for Isan Upid
func (upid *Upid) isan(gob *gobs.Gob, upidlen uint8) {
	upid.Value = gob.Ascii(uint(upidlen << 3))
}

// Decode for URI Upid
func (upid *Upid) uri(gob *gobs.Gob, upidlen uint8) {
	upid.Value = gob.Ascii(uint(upidlen) << 3)
}

// Decode for ATSC Upid
func (upid *Upid) atsc(gob *gobs.Gob, upidlen uint8) {
	upid.TSID = gob.UInt16(16)
	upid.Reserved = gob.UInt8(2)
	upid.EndOfDay = gob.UInt8(5)
	upid.UniqueFor = gob.UInt16(9)
	upid.ContentID = gob.Bytes(uint((upidlen - 4) << 3))
}

// Decode for EIDR Upid
func (upid *Upid) eidr(gob *gobs.Gob, upidlen uint8) {
	if upidlen == 12 {
		head := gob.UInt64(16)
		tail := gob.Hex(80)
		upid.Value = fmt.Sprintf("10%v/%v", head, tail)
	}
}

// Decode for MPU Upid
func (upid *Upid) mpu(gob *gobs.Gob, upidlen uint8) {
	ulb := uint(upidlen) << 3
	upid.FormatIdentifier = gob.Hex(32)
	upid.PrivateData = gob.Bytes(ulb - 32)
}

// Decode for MID Upid
func (upid *Upid) mid(gob *gobs.Gob, upidlen uint8) {
	var i uint8
	i = 0
	for i < upidlen {
		utype := gob.UInt8(8)
		i++
		ulen := gob.UInt8(8)
		i++
		i += ulen
		var mupid Upid
		upid.Decode(gob, utype, ulen)
		upid.Upids = append(upid.Upids, mupid)
	}
}
