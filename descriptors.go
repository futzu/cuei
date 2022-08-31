package cuei

//import "github.com/futzu/gob"

// AudioCmpt is a struct for AudioDscptr Components
type AudioCmpt struct {
	ComponentTag  uint8
	ISOCode       uint32
	BitstreamMode uint8
	NumChannels   uint8
	FullSrvcAudio bool
}

// SegCmpt Segmentation Descriptor Component
type SegCmpt struct {
	ComponentTag uint8
	PtsOffset    float64
}

type SpliceDescriptor struct {
	Tag                              uint8       `json:",omitempty"`
	Length                           uint8       `json:",omitempty"`
	Identifier                       string      `json:",omitempty"`
	Name                             string      `json:",omitempty"`
	AudioComponents                  []AudioCmpt `json:",omitempty"`
	ProviderAvailID                  uint32      `json:",omitempty"`
	PreRoll                          uint8       `json:",omitempty"`
	DTMFCount                        uint8       `json:",omitempty"`
	DTMFChars                        uint64      `json:",omitempty"`
	TAISeconds                       uint64      `json:",omitempty"`
	TAINano                          uint32      `json:",omitempty"`
	UTCOffset                        uint16      `json:",omitempty"`
	SegmentationEventID              string      `json:",omitempty"`
	SegmentationEventCancelIndicator bool        `json:",omitempty"`
	ProgramSegmentationFlag          bool        `json:",omitempty"`
	SegmentationDurationFlag         bool        `json:",omitempty"`
	DeliveryNotRestrictedFlag        bool        `json:",omitempty"`
	WebDeliveryAllowedFlag           bool        `json:",omitempty"`
	NoRegionalBlackoutFlag           bool        `json:",omitempty"`
	ArchiveAllowedFlag               bool        `json:",omitempty"`
	DeviceRestrictions               string      `json:",omitempty"`
	Components                       []SegCmpt   `json:",omitempty"`
	SegmentationDuration             float64     `json:",omitempty"`
	SegmentationMessage              string      `json:",omitempty"`
	SegmentationUpidType             uint8       `json:",omitempty"`
	SegmentationUpidLength           uint8       `json:",omitempty"`
	SegmentationUpid                 *Upid       `json:",omitempty"`
	SegmentationTypeID               uint8       `json:",omitempty"`
	SegmentNum                       uint8       `json:",omitempty"`
	SegmentsExpected                 uint8       `json:",omitempty"`
	SubSegmentNum                    uint8       `json:",omitempty"`
	SubSegmentsExpected              uint8       `json:",omitempty"`
}

// DescriptorDecoder returns a Descriptor by tag
func (dscptr *SpliceDescriptor) Decoder(gob *Gob, tag uint8, length uint8) {
	switch tag {
	case 0:
		dscptr.Tag = 0
		dscptr.Avail(gob, tag, length)
	case 1:
		dscptr.Tag = 1
		dscptr.DTMF(gob, tag, length)
	case 2:
		dscptr.Tag = 2
		dscptr.Segmentation(gob, tag, length)
	case 3:
		dscptr.Tag = 3
		dscptr.Time(gob, tag, length)
	case 4:
		dscptr.Tag = 4
		dscptr.Audio(gob, tag, length)
	}
}

func (dscptr *SpliceDescriptor) Audio(gob *Gob, tag uint8, length uint8) {
	dscptr.Tag = tag
	dscptr.Length = length
	dscptr.Identifier = gob.Ascii(32)
	ccount := gob.uint8(4)
	gob.Forward(4)
	for ccount > 0 {
		ccount--
		ct := gob.uint8(8)
		iso := gob.uint32(24)
		bsm := gob.uint8(3)
		nc := gob.uint8(4)
		fsa := gob.Bool()
		dscptr.AudioComponents = append(dscptr.AudioComponents, AudioCmpt{ct, iso, bsm, nc, fsa})
	}
}

// Decode for the Avail
func (dscptr *SpliceDescriptor) Avail(gob *Gob, tag uint8, length uint8) {
	dscptr.Tag = tag
	dscptr.Length = length
	dscptr.Identifier = gob.Ascii(32)
	dscptr.Name = "Avail Descriptor"
	dscptr.ProviderAvailID = gob.uint32(32)
}

//  DTMF Splice Descriptor
func (dscptr *SpliceDescriptor) DTMF(gob *Gob, tag uint8, length uint8) {
	dscptr.Tag = tag
	dscptr.Length = length
	dscptr.Identifier = gob.Ascii(32)
	dscptr.Name = "DTMF Descriptor"
	dscptr.PreRoll = gob.uint8(8)
	dscptr.DTMFCount = gob.uint8(3)
	//gob.Forward(5)
	dscptr.DTMFChars = gob.uint64(uint(8 * dscptr.DTMFCount))

}

// Decode for the Time Descriptor
func (dscptr *SpliceDescriptor) Time(gob *Gob, tag uint8, length uint8) {
	dscptr.Tag = tag
	dscptr.Length = length
	dscptr.Identifier = gob.Ascii(32)
	dscptr.Name = "Time Descriptor"
	dscptr.TAISeconds = gob.uint64(48)
	dscptr.TAINano = gob.uint32(32)
	dscptr.UTCOffset = gob.uint16(16)
}

// Decode for the Segmentation Descriptor
func (dscptr *SpliceDescriptor) Segmentation(gob *Gob, tag uint8, length uint8) {
	dscptr.Tag = tag
	dscptr.Length = length
	dscptr.Identifier = gob.Ascii(32)
	dscptr.Name = "Segmentation Descriptor"
	dscptr.SegmentationEventID = gob.Hex(32)
	dscptr.SegmentationEventCancelIndicator = gob.Bool()
	gob.Forward(7)
	if !dscptr.SegmentationEventCancelIndicator {
		dscptr.decodeSegFlags(gob)
		if !dscptr.ProgramSegmentationFlag {
			dscptr.decodeSegCmpnts(gob)
		}
		dscptr.decodeSegmentation(gob)
	}
}

func (dscptr *SpliceDescriptor) decodeSegFlags(gob *Gob) {
	dscptr.ProgramSegmentationFlag = gob.Bool()
	dscptr.SegmentationDurationFlag = gob.Bool()
	dscptr.DeliveryNotRestrictedFlag = gob.Bool()
	if dscptr.DeliveryNotRestrictedFlag == false {
		dscptr.WebDeliveryAllowedFlag = gob.Bool()
		dscptr.NoRegionalBlackoutFlag = gob.Bool()
		dscptr.ArchiveAllowedFlag = gob.Bool()
		dscptr.DeviceRestrictions = table20[gob.uint8(2)]
		return
	}
	gob.Forward(5)
}

func (dscptr *SpliceDescriptor) decodeSegCmpnts(gob *Gob) {
	ccount := gob.uint8(8)
	for ccount > 0 { // 6 bytes each
		ccount--
		ct := gob.uint8(8)
		gob.Forward(7)
		po := gob.As90k(33)
		dscptr.Components = append(dscptr.Components, SegCmpt{ct, po})
	}
}

func (dscptr *SpliceDescriptor) decodeSegmentation(gob *Gob) {
	if dscptr.SegmentationDurationFlag == true {
		dscptr.SegmentationDuration = gob.As90k(40)
	}
	dscptr.SegmentationUpidType = gob.uint8(8)
	dscptr.SegmentationUpidLength = gob.uint8(8)
	dscptr.SegmentationUpid = &Upid{}
	dscptr.SegmentationUpid.Decoder(gob, dscptr.SegmentationUpidType, dscptr.SegmentationUpidLength)
	dscptr.SegmentationTypeID = gob.uint8(8)

	mesg, ok := table22[dscptr.SegmentationTypeID]
	if ok {
		dscptr.SegmentationMessage = mesg
	}
	dscptr.SegmentNum = gob.uint8(8)
	dscptr.SegmentsExpected = gob.uint8(8)
	subSegIDs := []uint8{0x34, 0x36, 0x38, 0x3a}
	if isIn8(subSegIDs, dscptr.SegmentationTypeID) {
		dscptr.SubSegmentNum = gob.uint8(8)
		dscptr.SubSegmentsExpected = gob.uint8(8)
	}
}
