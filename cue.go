package cuei

import (
	"fmt"
	gobs "github.com/futzu/gob"
)

/*
*
Cue is a SCTE35 cue.

	A Cue contains:
	    1 InfoSection
	    1 SpliceCommand
	    1 Dll  (Descriptor loop length)
	    0 or more Splice Descriptors
	    1 Crc32
	    1 packetData (if parsed from MPEGTS)

*
*/
type Cue struct {
	Bites       []byte
	InfoSection *InfoSection
	Command     *SpliceCommand
	Dll         uint16
	Descriptors []SpliceDescriptor `json:",omitempty"`
	PacketData  *packetData        `json:",omitempty"`
	Crc32       uint32
}

// Decode extracts bits for the Cue values.
func (cue *Cue) Decode(bites []byte) bool {
	var gob gobs.Gob
	gob.Load(bites)
	cue.InfoSection = &InfoSection{}
	if cue.InfoSection.Decode(&gob) {
		cue.Command = &SpliceCommand{}
		cue.Command.Decode(cue.InfoSection.SpliceCommandType, &gob)
		cue.Dll = gob.UInt16(16)
		cue.dscptrLoop(cue.Dll, &gob)
		cue.Crc32 = gob.UInt32(32)
		return true
	}
	return false
}

// DscptrLoop loops over any splice descriptors
func (cue *Cue) dscptrLoop(dll uint16, gob *gobs.Gob) {
	var i uint16
	i = 0
	l := dll
	for i < l {
		tag := gob.UInt8(8)
		i++
		length := gob.UInt16(8)
		i++
		i += length
		var sdr SpliceDescriptor
		sdr.Decode(gob, tag, uint8(length))
		cue.Descriptors = append(cue.Descriptors, sdr)
	}
}

// Show display SCTE-35 data as JSON.
func (cue *Cue) Show() {
	fmt.Println(MkJson(&cue))
}

// Encode Cue currently works for Splice Inserts and Time Signals
func (cue *Cue) Encode() []byte {
	cmdb := cue.Command.Encode()
	cmdl := len(cmdb)
	cue.InfoSection.SpliceCommandLength = uint16(cmdl)
	cue.InfoSection.SpliceCommandType = cue.Command.CommandType
	// 11 bytes for info section + command + 2 descriptor loop length
	// + descriptor loop + 4 for crc
	cue.InfoSection.SectionLength = uint16(11 + cmdl + 2 + 4)
	cue.InfoSection.Encode()
	nb := &Nbin{}
	isbits := uint(len(cue.InfoSection.Bites) << 3)
	nb.AddBytes(cue.InfoSection.Bites, isbits)
	ccbits := uint(len(cmdb) << 3)
	nb.AddBytes(cmdb, ccbits)
	nb.Add16(0, 16)
	cue.Crc32 = CRC32(nb.Bites.Bytes())
	nb.Add32(cue.Crc32, 32)
	cue.Bites = nb.Bites.Bytes()
	return nb.Bites.Bytes()
}
