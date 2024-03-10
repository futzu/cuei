package cuei

import (
	"fmt"
)

// InfoSection is the splice info section of the SCTE 35 cue.
type InfoSection struct {
	Name                   string
	TableID                string
	SectionSyntaxIndicator bool
	Private                bool
	SapType		       uint8
	SapDetails             string
	SectionLength          uint16
	ProtocolVersion        uint8
	EncryptedPacket        bool
	EncryptionAlgorithm    uint8
	PtsAdjustment          float64
	CwIndex                string
	Tier                   string
	CommandLength          uint16
	CommandType            uint8
}

// decode Splice Info Section values.
func (infosec *InfoSection) decode(bd *bitDecoder) bool {
	infosec.Name = "Splice Info Section"
	infosec.TableID = bd.asHex(8)
	if infosec.TableID != "0xfc" {
		return false
	}
	infosec.SectionSyntaxIndicator = bd.asFlag()
	infosec.Private = bd.asFlag()
	infosec.SapType = bd.uInt8(2)
	details, ok := table6[infosec.SapType]
	if ok {
		infosec.SapDetails = details
	}
	infosec.SectionLength = bd.uInt16(12)
	infosec.ProtocolVersion = bd.uInt8(8)
	infosec.EncryptedPacket = bd.asFlag()
	infosec.EncryptionAlgorithm = bd.uInt8(6)
	infosec.PtsAdjustment = bd.as90k(33)
	infosec.CwIndex = bd.asHex(8)
	infosec.Tier = bd.asHex(12)
	infosec.CommandLength = bd.uInt16(12)
	infosec.CommandType = bd.uInt8(8)

	return true
}

// defaults sets default InfoSection values for encoding
func (infosec *InfoSection) defaults() {
	infosec.Name = "Splice Info Section"
	infosec.TableID = "0xfc"
	infosec.SectionSyntaxIndicator = false
	infosec.Private = false
	infosec.SapType = 0x3
	//infosec.SectionLength = 17
	infosec.ProtocolVersion = 0
	infosec.EncryptedPacket = false
	infosec.EncryptionAlgorithm = 0
	infosec.PtsAdjustment = 0.0
	infosec.CwIndex = "0x0"
	infosec.Tier = "0xfff"
	infosec.CommandLength = 0
	infosec.CommandType = 0
}

/*
Encode Splice Info Section
Encodes the InfoSection variables to bytes.
*/
func (infosec *InfoSection) encode() []byte {
	be := &bitEncoder{}
	be.Add(uint16(0xfc), 16)
	be.Add(infosec.SectionSyntaxIndicator,1)
	be.Add(infosec.Private,1)
	be.Add(infosec.SapType,2)
	be.Add(uint16(infosec.SectionLength), 12)
	be.Add(infosec.ProtocolVersion, 8)
	be.Add(infosec.EncryptedPacket, 1)
	be.Add(infosec.EncryptionAlgorithm, 6)
	be.Add(infosec.PtsAdjustment, 33)
	be.AddHex64(infosec.CwIndex, 8)
	be.AddHex64(infosec.Tier, 12)
	be.Add(infosec.CommandLength, 12)
	be.Add(infosec.CommandType, 8)
	return be.Bites.Bytes()

}

// Return InfoSection as JSON
func (infosec *InfoSection) Json() string {
	return mkJson(infosec)
}

// Print InfoSection as JSON
func (infosec *InfoSection) Show() {
	fmt.Printf(infosec.Json())
}
