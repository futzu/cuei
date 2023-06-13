package cuei

import (
	gobs "github.com/futzu/gob"
)

// InfoSection is the splice info section of the SCTE 35 cue.
type InfoSection struct {
	Name                   string
	TableID                string
	SectionSyntaxIndicator bool
	Private                bool
	Reserved               string
	SectionLength          uint16
	ProtocolVersion        uint8
	EncryptedPacket        bool
	EncryptionAlgorithm    uint8
	PtsAdjustment          float64
	CwIndex                string
	Tier                   string
	SpliceCommandLength    uint16
	SpliceCommandType      uint8
	DescriptorLoopLength   uint16
	CRC                    uint32
}

// Decode Splice Info Section values.
func (infosec *InfoSection) Decode(gob *gobs.Gob) bool {
	infosec.Name = "Splice Info Section"
	infosec.TableID = gob.Hex(8)
	if infosec.TableID != "0xfc" {
		return false
	}
	infosec.SectionSyntaxIndicator = gob.Flag()
	infosec.Private = gob.Flag()
	infosec.Reserved = gob.Hex(2)
	infosec.SectionLength = gob.UInt16(12)
	infosec.ProtocolVersion = gob.UInt8(8)
	infosec.EncryptedPacket = gob.Flag()
	infosec.EncryptionAlgorithm = gob.UInt8(6)
	infosec.PtsAdjustment = gob.As90k(33)
	infosec.CwIndex = gob.Hex(8)
	infosec.Tier = gob.Hex(12)
	infosec.SpliceCommandLength = gob.UInt16(12)
	infosec.SpliceCommandType = gob.UInt8(8)
	return true
}

// Defaults sets default InfoSection values for encoding
func (infosec *InfoSection) Defaults() {
	infosec.Name = "Splice Info Section"
	infosec.TableID = "0xfc"
	infosec.SectionSyntaxIndicator = false
	infosec.Private = false
	infosec.Reserved = "0x3"
	infosec.SectionLength = 11
	infosec.ProtocolVersion = 0
	infosec.EncryptedPacket = false
	infosec.EncryptionAlgorithm = 0
	infosec.PtsAdjustment = 0.0
	infosec.CwIndex = "0x0"
	infosec.Tier = "0xfff"
	infosec.SpliceCommandLength = 0
	infosec.SpliceCommandType = 0
	infosec.DescriptorLoopLength = 0
}

// Encode Splice Info Section
func (infosec *InfoSection) Encode() []byte {
	infosec.Defaults()
	nb := &Nbin{}
	nb.AddHex64(infosec.TableID, 8)
	nb.AddFlag(infosec.SectionSyntaxIndicator)
	nb.AddFlag(infosec.Private)
	nb.Reserve(2)
	nb.Add16(infosec.SectionLength, 12)
	nb.Add8(infosec.ProtocolVersion, 8)
	nb.AddFlag(infosec.EncryptedPacket)
	nb.Add8(infosec.EncryptionAlgorithm, 6)
	nb.Add90k(infosec.PtsAdjustment, 33)
	nb.AddHex64(infosec.CwIndex, 8)
	nb.AddHex64(infosec.Tier, 12)
	nb.Add16(infosec.SpliceCommandLength, 12)
	nb.Add8(infosec.SpliceCommandType, 8)
	return nb.Bites.Bytes()
}
