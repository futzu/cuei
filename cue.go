package cuei

import (
	"fmt"
)

/*
*
Cue is a SCTE35 cue.

A Cue contains:

		1 InfoSection
	   	1 Command
	   	1 Dll  Descriptor loop length
	   	0 or more Splice Descriptors
	   	1 Crc32
	   	1 packetData (if parsed from MPEGTS)

*
*/
type Cue struct {
	InfoSection *InfoSection
	Command     *Command
	Dll         uint16       `json:"DescriptorLoopLength"`
	Descriptors []Descriptor `json:",omitempty"`
	PacketData  *packetData  `json:",omitempty"`
	Crc32       uint32
}

// Decode extracts bits for the Cue values.
func (cue *Cue) Decode(bites []byte) bool {
	var bd BitDecoder
	bd.Load(bites)
	cue.InfoSection = &InfoSection{}
	if cue.InfoSection.Decode(&bd) {
		cue.Command = &Command{}
		cue.Command.Decode(cue.InfoSection.CommandType, &bd)
		cue.Dll = bd.UInt16(16)
		cue.dscptrLoop(cue.Dll, &bd)
		cue.Crc32 = bd.UInt32(32)
		return true
	}
	return false
}

// DscptrLoop loops over any splice descriptors
func (cue *Cue) dscptrLoop(dll uint16, bd *BitDecoder) {
	var i uint16
	i = 0
	l := dll
	for i < l {
		tag := bd.UInt8(8)
		i++
		length := bd.UInt16(8)
		i++
		i += length
		var sdr Descriptor
		sdr.Decode(bd, tag, uint8(length))
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
	cue.InfoSection.CommandLength = uint16(cmdl)
	cue.InfoSection.CommandType = cue.Command.CommandType
	// 11 bytes for info section + command + 2 descriptor loop length
	// + descriptor loop + 4 for crc
	cue.InfoSection.SectionLength = uint16(11 + cmdl + 2 + 4)
	isecb := cue.InfoSection.Encode()
	be := &BitEncoder{}
	isecbits := uint(len(isecb) << 3)
	be.AddBytes(isecb, isecbits)
	cmdbits := uint(cmdl << 3)
	be.AddBytes(cmdb, cmdbits)
	be.Add(0, 16) // descriptor loop currently disabled for encoding
	cue.Crc32 = CRC32(be.Bites.Bytes())
	be.Add(cue.Crc32, 32)
	return be.Bites.Bytes()
}

// used by Six2Five to convert a time signal to a splice insert
func (cue *Cue) mkSpliceInsert() {
	cue.Command.CommandType = 5
	cue.Command.Name = "Splice Insert"
	cue.InfoSection.CommandType = 5
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
}

/*
	 *

		Convert  Cue.Command  from a  Time Signal
		to a Splice Insert and return a base64 string

		Example Usage:

			package main

		import (
			"os"
			"fmt"
			"github.com/futzu/cuei"
		)

		func main() {
			args := os.Args[1:]
			for _,arg := range args {
				fmt.Printf("\nNext File: %s\n\n", arg)
				stream := cuei.NewStream()
				cues :=stream.Decode(arg)
				for _,c:= range cues {
					fmt.Println(c.Six2Five())
				}
			}
		}

*
*/
func (cue *Cue) Six2Five() string {
	segStarts := []uint16{0x22, 0x30, 0x32, 0x34, 0x36, 0x38, 0x3a, 0x3c, 0x3e, 0x44, 0x46}
	segStops := []uint16{0x23, 0x31, 0x33, 0x35, 0x37, 0x39, 0x3b, 0x3d, 0x3f, 0x45, 0x47}
	if cue.InfoSection.CommandType == 6 {
		for _, dscptr := range cue.Descriptors {
			if dscptr.Tag == 2 {
				//value, _ := strconv.ParseInt(hex, 16, 64)
				cue.Command.SpliceEventID = uint32(Hex2Int(dscptr.SegmentationEventID))
				if IsIn(segStarts, uint16(dscptr.SegmentationTypeID)) {
					if dscptr.SegmentationDurationFlag {
						cue.mkSpliceInsert()
						cue.Command.OutOfNetworkIndicator = true
						cue.Command.DurationFlag = true
						cue.Command.BreakAutoReturn = true
						cue.Command.BreakDuration = dscptr.SegmentationDuration
						return EncB64(cue.Encode())
					}
				}
				if IsIn(segStops, uint16(dscptr.SegmentationTypeID)) {
					cue.mkSpliceInsert()
					cue2 := NewCue()
					cue2.Decode(cue.Encode())
					return EncB64(cue2.Encode())
				}
			}
		}
	}
	return EncB64(cue.Encode())

}

// initialize and return a *Cue
func NewCue() *Cue {
	cue := &Cue{}
	return cue
}
