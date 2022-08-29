package cuei

import(
    "github.com/futzu/gob"
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

// Decode splice info section values.
func (infosec *InfoSection) Decode(gob *gob.Gob) bool {
	infosec.Name = "Splice Info Section"
	infosec.TableID = gob.Hex(8)
	if infosec.TableID != "0xfc" {
		return false
	}
	infosec.SectionSyntaxIndicator = gob.Bool()
	if infosec.SectionSyntaxIndicator {
		return false
	}
	infosec.Private = gob.Bool()
	infosec.Reserved = gob.Hex(2)
	infosec.SectionLength = gob.UInt16(12)
	infosec.ProtocolVersion = gob.UInt8(8)
	if infosec.ProtocolVersion != 0 {
		return false
	}
	infosec.EncryptedPacket = gob.Bool()
	infosec.EncryptionAlgorithm = gob.UInt8(6)
	infosec.PtsAdjustment = gob.As90k(33)
	infosec.CwIndex = gob.Hex(8)
	infosec.Tier = gob.Hex(12)
	infosec.SpliceCommandLength = gob.UInt16(12)
	infosec.SpliceCommandType = gob.UInt8(8)
	return true
}
