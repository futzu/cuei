package cuei

import (
	goober "github.com/futzu/gob"
)

// audioCmpt is a struct for audioDscptr Components
type audioCmpt struct {
	ComponentTag  uint8
	ISOCode       uint32
	BitstreamMode uint8
	NumChannels   uint8
	FullSrvcAudio bool
}

// segCmpt Segmentation Descriptor Component
type segCmpt struct {
	ComponentTag uint8
	PtsOffset    float64
}

type SpliceDescriptor struct {
	Tag                              uint8       `json:",omitempty"`
	Length                           uint8       `json:",omitempty"`
	Identifier                       string      `json:",omitempty"`
	Name                             string      `json:",omitempty"`
	AudioComponents                  []audioCmpt `json:",omitempty"`
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
	Components                       []segCmpt   `json:",omitempty"`
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
func (dscptr *SpliceDescriptor) Decoder(gob *goober.Gob, tag uint8, length uint8) {
	switch tag {
	case 0:
		dscptr.Tag = 0
		dscptr.availDescriptor(gob, tag, length)
	case 1:
		dscptr.Tag = 1
		dscptr.dtmfDescriptor(gob, tag, length)
	case 2:
		dscptr.Tag = 2
		dscptr.segmentationDescriptor(gob, tag, length)
	case 3:
		dscptr.Tag = 3
		dscptr.timeDescriptor(gob, tag, length)
	case 4:
		dscptr.Tag = 4
		dscptr.audioDescriptor(gob, tag, length)
	}
}

func (dscptr *SpliceDescriptor) audioDescriptor(gob *goober.Gob, tag uint8, length uint8) {
	dscptr.Tag = tag
	dscptr.Length = length
	dscptr.Identifier = gob.Ascii(32)
	ccount := gob.UInt8(4)
	gob.Forward(4)
	for ccount > 0 {
		ccount--
		ct := gob.UInt8(8)
		iso := gob.UInt32(24)
		bsm := gob.UInt8(3)
		nc := gob.UInt8(4)
		fsa := gob.Flag()
		dscptr.AudioComponents = append(dscptr.AudioComponents, audioCmpt{ct, iso, bsm, nc, fsa})
	}
}

// Decode for the avail Splice Descriptors
func (dscptr *SpliceDescriptor) availDescriptor(gob *goober.Gob, tag uint8, length uint8) {
	dscptr.Tag = tag
	dscptr.Length = length
	dscptr.Identifier = gob.Ascii(32)
	dscptr.Name = "Avail Descriptor"
	dscptr.ProviderAvailID = gob.UInt32(32)
}

// DTMF Splice Descriptor
func (dscptr *SpliceDescriptor) dtmfDescriptor(gob *goober.Gob, tag uint8, length uint8) {
	dscptr.Tag = tag
	dscptr.Length = length
	dscptr.Identifier = gob.Ascii(32)
	dscptr.Name = "DTMF Descriptor"
	dscptr.PreRoll = gob.UInt8(8)
	dscptr.DTMFCount = gob.UInt8(3)
	//gob.Forward(5)
	dscptr.DTMFChars = gob.UInt64(uint(8 * dscptr.DTMFCount))

}

// Decode for the Time Descriptor
func (dscptr *SpliceDescriptor) timeDescriptor(gob *goober.Gob, tag uint8, length uint8) {
	dscptr.Tag = tag
	dscptr.Length = length
	dscptr.Identifier = gob.Ascii(32)
	dscptr.Name = "Time Descriptor"
	dscptr.TAISeconds = gob.UInt64(48)
	dscptr.TAINano = gob.UInt32(32)
	dscptr.UTCOffset = gob.UInt16(16)
}

// Decode for the Segmentation Descriptor
func (dscptr *SpliceDescriptor) segmentationDescriptor(gob *goober.Gob, tag uint8, length uint8) {
	dscptr.Tag = tag
	dscptr.Length = length
	dscptr.Identifier = gob.Ascii(32)
	dscptr.Name = "Segmentation Descriptor"
	dscptr.SegmentationEventID = gob.Hex(32)
	dscptr.SegmentationEventCancelIndicator = gob.Flag()
	gob.Forward(7)
	if !dscptr.SegmentationEventCancelIndicator {
		dscptr.decodeSegFlags(gob)
		if !dscptr.ProgramSegmentationFlag {
			dscptr.decodeSegCmpnts(gob)
		}
		dscptr.decodeSegmentation(gob)
	}
}

func (dscptr *SpliceDescriptor) decodeSegFlags(gob *goober.Gob) {
	dscptr.ProgramSegmentationFlag = gob.Flag()
	dscptr.SegmentationDurationFlag = gob.Flag()
	dscptr.DeliveryNotRestrictedFlag = gob.Flag()
	if !dscptr.DeliveryNotRestrictedFlag {
		dscptr.WebDeliveryAllowedFlag = gob.Flag()
		dscptr.NoRegionalBlackoutFlag = gob.Flag()
		dscptr.ArchiveAllowedFlag = gob.Flag()
		dscptr.DeviceRestrictions = table20[gob.UInt8(2)]
		return
	}
	gob.Forward(5)
}

func (dscptr *SpliceDescriptor) decodeSegCmpnts(gob *goober.Gob) {
	ccount := gob.UInt8(8)
	for ccount > 0 { // 6 bytes each
		ccount--
		ct := gob.UInt8(8)
		gob.Forward(7)
		po := gob.As90k(33)
		dscptr.Components = append(dscptr.Components, segCmpt{ct, po})
	}
}
