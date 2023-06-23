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
		1 Dll  Descriptor loop length
	    0 or more Splice Descriptors
		1 Crc32
	    1 packetData (if parsed from MPEGTS)

*
*/
type Cue struct {
	//	Bites       []byte
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
	isecb := cue.InfoSection.Encode()
	nb := &Nbin{}
	isecbits := uint(len(isecb) << 3)
	nb.AddBytes(isecb, isecbits)
	cmdbits := uint(cmdl << 3)
	nb.AddBytes(cmdb, cmdbits)
	nb.Add16(0, 16) // descriptor loop currently disabled for encoding
	cue.Crc32 = CRC32(nb.Bites.Bytes())
	nb.Add32(cue.Crc32, 32)
	//cue.Bites = nb.Bites.Bytes()
	return nb.Bites.Bytes()
}

// Convert  Cue.Command  from a  Time Signal to a Splice Insert and return a base64 styring
func (cue *Cue) Six2Five() string {
	upidStarts := []uint16{0x34, 0x36, 0x38}
	// upidStops := []uint16{0x35, 0x37, 0x39}
	if cue.InfoSection.SpliceCommandType == 6 {

		cue.Command.CommandType = 5
		cue.Command.Name = "Six 2 Five"
		cue.InfoSection.SpliceCommandType = 5
		cue.Command.ProgramSpliceFlag = true
		cue.Command.SpliceEventCancelIndicator = false
		cue.Command.OutOfNetworkIndicator = false
		cue.Command.TimeSpecifiedFlag = false
		cue.Command.DurationFlag = false
		cue.Command.BreakAutoReturn = false
		cue.Command.SpliceImmediateFlag = false
		cue.Command.AvailNum = 0
		cue.Command.AvailExpected = 0
		if cue.Command.PTS > 0.0 {
			cue.Command.TimeSpecifiedFlag = true
			cue.Command.PTS = cue.Command.PTS
		}
		for _, dscptr := range cue.Descriptors {
			if dscptr.Tag == 2 {
				//value, _ := strconv.ParseInt(hex, 16, 64)
				cue.Command.SpliceEventID = uint32(5) //Hex2Int(dscptr.SegmentationEventID)&uint64(2^31)
				if IsIn(upidStarts, uint16(dscptr.SegmentationTypeID)) {
					if dscptr.SegmentationDurationFlag {
						cue.Command.OutOfNetworkIndicator = true
						cue.Command.DurationFlag = true
						cue.Command.BreakAutoReturn = true
						cue.Command.BreakDuration = dscptr.SegmentationDuration
					}
				}
			}

		}
	}
	fmt.Println("Six 2 Five")
	cue.Show()
	return EncB64(cue.Encode())

}
