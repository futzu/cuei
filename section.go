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
