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
	    0 or more Splice Descriptors
	    1 packetData (if parsed from MPEGTS)

*
*/
type Cue struct {
	Bites       []byte
	InfoSection *InfoSection
	Command     *SpliceCommand
	Descriptors []SpliceDescriptor `json:",omitempty"`
	PacketData  *packetData        `json:",omitempty"`
}

// Decode extracts bits for the Cue values.
func (cue *Cue) Decode(bites []byte) bool {
	var gob gobs.Gob
	gob.Load(bites)
	cue.InfoSection = &InfoSection{}
	if cue.InfoSection.Decode(&gob) {
		cue.Command = &SpliceCommand{}
		cue.Command.Decode(cue.InfoSection.SpliceCommandType, &gob)
		cue.InfoSection.DescriptorLoopLength = gob.UInt16(16)
		cue.dscptrLoop(&gob)
		return true
	}
	return false
}

// DscptrLoop loops over any splice descriptors
func (cue *Cue) dscptrLoop(gob *gobs.Gob) {
	var i uint16
	i = 0
	l := cue.InfoSection.DescriptorLoopLength
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

func (cue *Cue) Encode() []byte {
	//cmdl := len(cue.Command.Bites)
	//cue.Command.MkSpliceInsert("77",11.012344,30.1,true)
	cmdb := cue.Command.Encode()
	cmdl := len(cmdb)
	//fmt.Println(cmdb)

	cue.InfoSection.Defaults()
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
	dll := uint16(0)
	nb.Add16(dll, 16)
	/**     cue.Bites += int.to_bytes(
	            self.info_section.descriptor_loop_length, 2, byteorder="big"
	        )
	        self.bites += dscptr_bites

	**/
	crc32 := CRC32(nb.Bites.Bytes())
	nb.Add32(crc32, 32)
	cue.Bites = nb.Bites.Bytes()
	return nb.Bites.Bytes()
}

/*
  - Six2Five converts a Cue with a Time Sgnal Command
    and a Segmentation Descriptor with a
    type id of 0x34,0x35,0x36,0x37,0x38. or 0x39
    into a Cue with a Splice Insert Command.

*
*/
func (cue *Cue) Six2Five() {
	upidStarts := []uint16{0x34, 0x36, 0x38}
	upidStops := []uint16{0x35, 0x37, 0x39}
	eventid := "0x0"
	pts := 0.0
	duration := float64(0.0)
	out := false
	if cue.Command.CommandType == 6 {
		cue.Command.CommandType = 5

		cue.Command.Name = "Six2Five'd Splice Insert"
		if cue.Command.PTS > 0.0 {
			pts = cue.Command.PTS
		}
		for _, dscptr := range cue.Descriptors {
			if dscptr.Tag == 2 {
				//value, _ := strconv.ParseInt(hex, 16, 64)
				eventid = fmt.Sprintf("%v", Hex2Int(dscptr.SegmentationEventID)&uint64(2^31))
				if IsIn(upidStarts, uint16(dscptr.SegmentationTypeID)) {
					if dscptr.SegmentationDurationFlag {
						duration = dscptr.SegmentationDuration
						out = true
					}
				}
				if IsIn(upidStops, uint16(dscptr.SegmentationTypeID)) {
					out = false
				}
			}
			cue.Command.MkSpliceInsert(eventid, pts, duration, out)
			cue.Command.Name = "Six2Five'd Splice Insert"

			cue.Encode()
			//cue.Decode(cue.Encode())
			cue.Show()
			fmt.Println("Six 2 Five")
		}
	}
}
