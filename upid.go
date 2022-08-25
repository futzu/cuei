package cuei

import (
	"fmt"
)

// Upid is the Struct for Segmentation Upida
type Upid struct{
	Name  	  string        `json:",omitempty"`
    	UpidType  uint8         `json:",omitempty"`
	Value 	  string        `json:",omitempty"`
    	TSID      uint64        `json:",omitempty"`
	Reserved  uint8         `json:",omitempty"`
	EndOfDay  uint8         `json:",omitempty"`
	UniqueFor uint64        `json:",omitempty"`
	ContentID string        `json:",omitempty"`
    	Upids []Upid            `json:",omitempty"`
    	FormatIdentifier string `json:",omitempty"`
	PrivateData      string `json:",omitempty"`

}


// UpidDecoder calls a method based on upidType
func (upid *Upid) Decoder(bitn *Bitn, upidType uint8,  upidlen uint8) {
    
    upid.UpidType = upidType
	
    switch upidType {
	case 0x01, 0x02:
        upid.Name = "Deprecated"
		upid.URI(bitn,upidlen)
	case 0x03:
		upid.Name = "AdID"
        upid.URI(bitn,upidlen)
	case 0x05, 0x06:
		upid.Name ="ISAN"
        upid.ISAN(bitn,upidlen)
	case 0x07:
		upid.Name="TID"
        upid.URI(bitn,upidlen)
	case 0x08:
		upid.Name="AiringID"
        upid.AirID(bitn,upidlen)
	case 0x09:
		upid.Name ="ADI"
        upid.URI(bitn,upidlen)
	case 0x0a:
		upid.Name= "EIDR"
        upid.EIDR(bitn,upidlen)
	case 0x0b:
		upid.Name = "ATSC"
        upid.ATSC(bitn,upidlen)
	case 0x0c:
		upid.Name= "MPU"
        upid.MPU(bitn,upidlen)
	case 0x0d:
		upid.Name= "MID"
        upid.MID(bitn,upidlen)
	case 0x0e:
		upid.Name= "ADS Info"
        upid.URI(bitn,upidlen)
	case 0x0f:
		upid.Name= "URI"
        upid.URI(bitn,upidlen)
	case 0x10:
		upid.Name= "UUID"
        upid.URI(bitn,upidlen)
	default:
		upid.Name= "UPID"
        upid.URI(bitn,upidlen)

	}

}

// Decode for AirId
func (upid *Upid) AirID(bitn *Bitn, upidlen uint8) {
	upid.Value = bitn.AsHex(uint(upidlen << 3))
}


// Decode for Isan Upid
func (upid *Upid) ISAN(bitn *Bitn, upidlen uint8) {
	upid.Value = bitn.AsAscii(uint(upidlen << 3))
}

// Decode for URI Upid
func (upid *Upid) URI(bitn *Bitn, upidlen uint8) {
	upid.Value = bitn.AsAscii(uint(upidlen) << 3)
}

// Decode for ATSC Upid
func (upid *Upid) ATSC(bitn *Bitn, upidlen uint8) {
	upid.TSID = bitn.AsUInt64(16)
	upid.Reserved = bitn.AsUInt8(2)
	upid.EndOfDay = bitn.AsUInt8(5)
	upid.UniqueFor = bitn.AsUInt64(9)
	upid.ContentID = bitn.AsAscii(uint((upidlen - 4) << 3))
}


// Decode for EIDR Upid
func (upid *Upid) EIDR(bitn *Bitn, upidlen uint8) {
	if upidlen == 12 {
		head := bitn.AsUInt64(16)
		tail := bitn.AsHex(80)
		upid.Value = fmt.Sprintf("10%v/%v", head, tail)
	}
}


// Decode for MPU Upid
func (upid *Upid) MPU(bitn *Bitn, upidlen uint8) {
	ulb := uint(upidlen) << 3
	upid.FormatIdentifier = bitn.AsHex(32)
	upid.PrivateData = bitn.AsAscii(ulb - 32)

}


// Decode for MID Upid
func (upid *Upid) MID(bitn *Bitn, upidlen uint8) {
	var i uint8
	i = 0
	for i < upidlen {
		utype := bitn.AsUInt8(8)
		i++
		ulen := bitn.AsUInt8(8)
		i++
		i += ulen
		var mupid Upid
		mupid.Decoder(bitn,utype, ulen)
		upid.Upids = append(upid.Upids, mupid)

	}
}
