package cuei

import (
	"encoding/json"
	"fmt"
)

// Tag, Length , Name and Identifier for Descriptors
type TagLenNameId struct {
	Tag        uint8
	Length     uint8
	Name       string
	Identifier string
}

// audioCmpt is a struct for audioDscptr Components
type AudioCmpt struct {
	ComponentTag  uint8
	ISOCode       uint32
	BitstreamMode uint8
	NumChannels   uint8
	FullSrvcAudio bool
}

// avail Descriptor
type AvailDescriptor struct {
	ProviderAvailID uint32
}

// DTMF Descriptor
type DTMFDescriptor struct {
	PreRoll   uint8
	DTMFCount uint8
	DTMFChars uint64
}

// Segmentation Descriptor
type SegmentationDescriptor struct {
	SegmentationEventID                    string
	SegmentationEventCancelIndicator       bool
	SegmentationEventIDComplianceIndicator bool
	ProgramSegmentationFlag                bool
	SegmentationDurationFlag               bool
	DeliveryNotRestrictedFlag              bool
	WebDeliveryAllowedFlag                 bool
	NoRegionalBlackoutFlag                 bool
	ArchiveAllowedFlag                     bool
	DeviceRestrictions                     string
	SegmentationDuration                   float64
	SegmentationMessage                    string
	SegmentationUpidType                   uint8
	SegmentationUpidLength                 uint8
	SegmentationUpid                       *Upid
	SegmentationTypeID                     uint8
	SegmentNum                             uint8
	SegmentsExpected                       uint8
	SubSegmentNum                          uint8
	SubSegmentsExpected                    uint8
}

/*
*

	Descriptor is the combination of all the descriptors
	this is to maintain dot notation in the Cue struct.

*
*/
type Descriptor struct {
	TagLenNameId
	AvailDescriptor
	DTMFDescriptor
	SegmentationDescriptor
	AudioComponents []AudioCmpt
	TAISeconds      uint64
	TAINano         uint32
	UTCOffset       uint16
}

func (dscptr *Descriptor) jsonAvailDescriptor() ([]byte, error) {
	return json.Marshal(&struct {
		TagLenNameId
		AvailDescriptor
	}{
		TagLenNameId:    dscptr.TagLenNameId,
		AvailDescriptor: dscptr.AvailDescriptor,
	})
}

func (dscptr *Descriptor) jsonDTMFDescriptor() ([]byte, error) {
	return json.Marshal(&struct {
		TagLenNameId
		DTMFDescriptor
	}{
		TagLenNameId:   dscptr.TagLenNameId,
		DTMFDescriptor: dscptr.DTMFDescriptor,
	})
}

func (dscptr *Descriptor) jsonSegmentationDescriptor() ([]byte, error) {
	return json.Marshal(&struct {
		TagLenNameId
		SegmentationDescriptor
	}{
		TagLenNameId:           dscptr.TagLenNameId,
		SegmentationDescriptor: dscptr.SegmentationDescriptor,
	})
}

/*
	 *
	    Custom MarshalJSON
	        Marshal a Descriptor into

	        0x0: AvailDescriptor,
		    0x1: DTMFDescriptor,
		    0x2: SegmentationDescriptor

	        or just return the Descriptor

*
*/
func (dscptr *Descriptor) MarshalJSON() ([]byte, error) {
	switch dscptr.Tag {
	case 0x0:
		return dscptr.jsonAvailDescriptor()
	case 0x1:
		return dscptr.jsonDTMFDescriptor()
	case 0x2:
		return dscptr.jsonSegmentationDescriptor()
	}
	type Funk Descriptor
	return json.Marshal(&struct{ *Funk }{(*Funk)(dscptr)})
}

// Return Descriptor as JSON
func (dscptr *Descriptor) Json() string {
	stuff, err := dscptr.MarshalJSON()
	chk(err)
	return string(stuff)
}

// Print Descriptor as JSON
func (dscptr *Descriptor) Show() {
	fmt.Printf(dscptr.Json())
}

/*
*
Decode returns a Splice Descriptor by tag.

	The following Splice Descriptors are recognized.

	    0x0: Avail Descriptor,
	    0x1: DTMF Descriptor,
	    0x2: Segmentation Descriptor,
	    0x3: Time Descriptor,
	    0x4: Audio Descrioptor,

*
*/
func (dscptr *Descriptor) decode(bd *bitDecoder, tag uint8, length uint8) {
	switch tag {
	case 0x0:
		dscptr.Tag = 0x0
		dscptr.decodeAvailDescriptor(bd, tag, length)
	case 0x1:
		dscptr.Tag = 0x1
		dscptr.decodeDTMFDescriptor(bd, tag, length)
	case 0x2:
		dscptr.Tag = 0x2
		dscptr.decodeSegmentationDescriptor(bd, tag, length)
	case 0x3:
		dscptr.Tag = 0x3
		dscptr.decodeTimeDescriptor(bd, tag, length)
	case 0x4:
		dscptr.Tag = 0x4
		dscptr.decodeAudioDescriptor(bd, tag, length)
	}
}

func (dscptr *Descriptor) decodeAudioDescriptor(bd *bitDecoder, tag uint8, length uint8) {
	dscptr.Tag = tag
	dscptr.Length = length
	dscptr.Identifier = bd.asAscii(32)
	ccount := bd.uInt8(4)
	bd.goForward(4)
	for ccount > 0 {
		ccount--
		ct := bd.uInt8(8)
		iso := bd.uInt32(24)
		bsm := bd.uInt8(3)
		nc := bd.uInt8(4)
		fsa := bd.asFlag()
		dscptr.AudioComponents = append(dscptr.AudioComponents, AudioCmpt{ct, iso, bsm, nc, fsa})
	}
}

// Decode for  Avail Descriptors
func (dscptr *Descriptor) decodeAvailDescriptor(bd *bitDecoder, tag uint8, length uint8) {
	dscptr.Tag = tag
	dscptr.Length = length
	dscptr.Identifier = bd.asAscii(32)
	dscptr.Name = "Avail Descriptor"
	dscptr.ProviderAvailID = bd.uInt32(32)
}

// Decode for DTMF Splice Descriptor
func (dscptr *Descriptor) decodeDTMFDescriptor(bd *bitDecoder, tag uint8, length uint8) {
	dscptr.Tag = tag
	dscptr.Length = length
	dscptr.Identifier = bd.asAscii(32)
	dscptr.Name = "DTMF Descriptor"
	dscptr.PreRoll = bd.uInt8(8)
	dscptr.DTMFCount = bd.uInt8(3)
	//bd.goForward(5)
	dscptr.DTMFChars = bd.uInt64(uint(8 * dscptr.DTMFCount))

}

// Decode for the Time Descriptor
func (dscptr *Descriptor) decodeTimeDescriptor(bd *bitDecoder, tag uint8, length uint8) {
	dscptr.Tag = tag
	dscptr.Length = length
	dscptr.Identifier = bd.asAscii(32)
	dscptr.Name = "Time Descriptor"
	dscptr.TAISeconds = bd.uInt64(48)
	dscptr.TAINano = bd.uInt32(32)
	dscptr.UTCOffset = bd.uInt16(16)
}

// Decode for the Segmentation Descriptor
func (dscptr *Descriptor) decodeSegmentationDescriptor(bd *bitDecoder, tag uint8, length uint8) {
	dscptr.Tag = tag
	dscptr.Length = length
	dscptr.Identifier = bd.asAscii(32)
	dscptr.Name = "Segmentation Descriptor"
	dscptr.SegmentationEventID = bd.asHex(32)
	dscptr.SegmentationEventCancelIndicator = bd.asFlag()
	dscptr.SegmentationEventIDComplianceIndicator = bd.asFlag()
	bd.goForward(6)
	if !dscptr.SegmentationEventCancelIndicator {
		dscptr.decodeSegFlags(bd)
		dscptr.decodeSegmentation(bd)
	}
}

func (dscptr *Descriptor) decodeSegFlags(bd *bitDecoder) {
	dscptr.ProgramSegmentationFlag = bd.asFlag()
	dscptr.SegmentationDurationFlag = bd.asFlag()
	dscptr.DeliveryNotRestrictedFlag = bd.asFlag()
	if !dscptr.DeliveryNotRestrictedFlag {
		dscptr.WebDeliveryAllowedFlag = bd.asFlag()
		dscptr.NoRegionalBlackoutFlag = bd.asFlag()
		dscptr.ArchiveAllowedFlag = bd.asFlag()
		dscptr.DeviceRestrictions = table20[bd.uInt8(2)] // 8
	} else {
		bd.goForward(5)
	}
}

func (dscptr *Descriptor) decodeSegmentation(bd *bitDecoder) {
	if dscptr.SegmentationDurationFlag {
		dscptr.SegmentationDuration = bd.as90k(40)
	}
	dscptr.SegmentationUpidType = bd.uInt8(8)
	dscptr.SegmentationUpidLength = bd.uInt8(8)
	if dscptr.SegmentationUpidLength > 0 {
		dscptr.SegmentationUpid = &Upid{}
		dscptr.SegmentationUpid.decode(bd, dscptr.SegmentationUpidType, dscptr.SegmentationUpidLength)
	}
	dscptr.SegmentationTypeID = bd.uInt8(8)
	mesg, ok := table22[dscptr.SegmentationTypeID]
	if ok {
		dscptr.SegmentationMessage = mesg
	}
	dscptr.SegmentNum = bd.uInt8(8)
	dscptr.SegmentsExpected = bd.uInt8(8)
	subSegIDs := []uint16{0x30, 0x32, 0x34, 0x36, 0x38, 0x3A, 0x44, 0x46}
	if IsIn(subSegIDs, uint16(dscptr.SegmentationTypeID)) {
		dscptr.SubSegmentNum = bd.uInt8(8)
		dscptr.SubSegmentsExpected = bd.uInt8(8)
		//dscptr.SubSegmentNum = 0
		//dscptr.SubSegmentsExpected = 0
	}
}

func (dscptr *Descriptor) encode(be *bitEncoder) {
	switch dscptr.Tag {
	case 0x0:
		dscptr.encodeAvailDescriptor(be)
	case 0x2:
		dscptr.encodeSegmentationDescriptor(be)
	}
}

// Encode for Avail Descriptors
func (dscptr *Descriptor) encodeAvailDescriptor(be *bitEncoder) {
	be.Add(uint32(dscptr.ProviderAvailID), 32)
}

// Encode a segmentation descriptor
func (dscptr *Descriptor) encodeSegmentationDescriptor(be *bitEncoder) {
	be.AddHex64(dscptr.SegmentationEventID, 32)
	be.Add(dscptr.SegmentationEventCancelIndicator, 1)
	be.Add(dscptr.SegmentationEventIDComplianceIndicator, 1)
	be.Reserve(6)
	if !dscptr.SegmentationEventCancelIndicator {
		dscptr.encodeFlags(be)
		dscptr.encodeSegmentation(be)
	}
}

func (dscptr *Descriptor) encodeFlags(be *bitEncoder) {
	be.Add(dscptr.ProgramSegmentationFlag, 1)
	be.Add(dscptr.SegmentationDurationFlag, 1)
	be.Add(dscptr.DeliveryNotRestrictedFlag, 1)
	if !dscptr.DeliveryNotRestrictedFlag {
		be.Add(dscptr.WebDeliveryAllowedFlag, 1)
		be.Add(dscptr.NoRegionalBlackoutFlag, 1)
		be.Add(dscptr.ArchiveAllowedFlag, 1)
		//   a_key = k_by_v(table20, dscptr.device_restrictions)
		//     nbin.add_int(a_key, 2)
		be.Add(3, 2) //  dscptr.device_restrictions
	} else {
		be.Reserve(5)
	}
}

func (dscptr *Descriptor) encodeSegmentation(be *bitEncoder) {
	if dscptr.SegmentationDurationFlag {
		be.Add(float64(dscptr.SegmentationDuration), 40)
	}
	be.Add(dscptr.SegmentationUpidType, 8)
	be.Add(dscptr.SegmentationUpidLength, 8)
	//be.Reserve(int(dscptr.SegmentationUpidLength <<3))
	if dscptr.SegmentationUpidLength > 0 {
		dscptr.SegmentationUpid.encode(be, dscptr.SegmentationUpidType)
	}
	be.Add(dscptr.SegmentationTypeID, 8)
	dscptr.encodeSegments(be)
}

func (dscptr *Descriptor) encodeSegments(be *bitEncoder) {
	be.Add(dscptr.SegmentNum, 8)
	be.Add(dscptr.SegmentsExpected, 8)
	subSegIDs := []uint16{0x30, 0x32, 0x34, 0x36, 0x38, 0x3A, 0x44, 0x46}
	if IsIn(subSegIDs, uint16(dscptr.SegmentationTypeID)) {
		be.Add(dscptr.SubSegmentNum, 8)
		be.Add(dscptr.SubSegmentsExpected, 8)
	}

}
