package cuei

import (
	"encoding/json"
	"fmt"
)

/*
Command

	These Splice Command types are consolidated into Command.

	     0x0: Splice Null,
	     0x5: Splice Insert,
	     0x6: Time Signal,
	     0x7: Bandwidth Reservation,
	     0xff: Private,
*/
type Command struct {
	Name                       string
	CommandType                uint8
	PrivateBytes               []byte
	Identifier                 uint32
	SpliceEventID              uint32
	SpliceEventCancelIndicator bool
	EventIDComplianceFlag      bool
	OutOfNetworkIndicator      bool
	ProgramSpliceFlag          bool
	DurationFlag               bool
	BreakAutoReturn            bool
	BreakDuration              float64
	SpliceImmediateFlag        bool
	UniqueProgramID            uint16
	AvailNum                   uint8
	AvailExpected              uint8
	TimeSpecifiedFlag          bool
	PTS                        float64
}

// only show TimeSignal values in JSON, used by cmd.MarshalJSON()
func (cmd *Command) jsonTimeSignal() ([]byte, error) {
	return json.Marshal(struct {
		Name              string
		CommandType       uint8
		TimeSpecifiedFlag bool
		PTS               float64
	}{
		Name:              cmd.Name,
		CommandType:       cmd.CommandType,
		TimeSpecifiedFlag: cmd.TimeSpecifiedFlag,
		PTS:               cmd.PTS})

}

// only show SpliceInsert values in JSON, used by cmd.MarshalJSON()
func (cmd *Command) jsonSpliceInsert() ([]byte, error) {
	return json.Marshal(struct {
		Name                       string
		CommandType                uint8
		SpliceEventID              uint32
		SpliceEventCancelIndicator bool
		OutOfNetworkIndicator      bool
		ProgramSpliceFlag          bool
		DurationFlag               bool
		BreakDuration              float64
		BreakAutoReturn            bool
		SpliceImmediateFlag        bool
		EventIDComplianceFlag      bool
		UniqueProgramID            uint16
		AvailNum                   uint8
		AvailExpected              uint8
		PTS                        float64
	}{
		Name:                       cmd.Name,
		CommandType:                cmd.CommandType,
		SpliceEventID:              cmd.SpliceEventID,
		SpliceEventCancelIndicator: cmd.SpliceEventCancelIndicator,
		OutOfNetworkIndicator:      cmd.OutOfNetworkIndicator,
		ProgramSpliceFlag:          cmd.ProgramSpliceFlag,
		DurationFlag:               cmd.DurationFlag,
		BreakDuration:              cmd.BreakDuration,
		BreakAutoReturn:            cmd.BreakAutoReturn,
		SpliceImmediateFlag:        cmd.SpliceImmediateFlag,
		EventIDComplianceFlag:      cmd.EventIDComplianceFlag,
		UniqueProgramID:            cmd.UniqueProgramID,
		AvailNum:                   cmd.AvailNum,
		AvailExpected:              cmd.AvailExpected,
		PTS:                        cmd.PTS})
}

func (cmd *Command) MarshalJSON() ([]byte, error) {
	switch cmd.CommandType {
	case 0x5:
		return cmd.jsonSpliceInsert()
	case 0x6:
		return cmd.jsonTimeSignal()
	}
	type Funk Command
	return json.Marshal(&struct{ *Funk }{(*Funk)(cmd)})
}

// Return Command as JSON
func (cmd *Command) Json() string {
	stuff, err := cmd.MarshalJSON()
	chk(err)
	return string(stuff)

}

// Print Command as JSON
func (cmd *Command) Show() {
	fmt.Printf(cmd.Json())
}

// Decode a Splice Command
func (cmd *Command) decode(cmdtype uint8, bd *bitDecoder) {
	cmd.CommandType = cmdtype
	switch cmdtype {
	case 0x0:
		cmd.decodeSpliceNull(bd)
	case 0x5:
		cmd.decodeSpliceInsert(bd)
	case 0x6:
		cmd.decodeTimeSignal(bd)
	case 0x7:
		cmd.decodeBandwidthReservation(bd)
	case 0xff:
		cmd.decodePrivate(bd)
	}

}

// Encode a Splice Command and return the bytes
// mostly used by cuei.Cue
func (cmd *Command) encode() []byte {
	blank := []byte{}
	switch cmd.CommandType {
	case 0x5:
		return cmd.encodeSpliceInsert()

	case 0x6:
		return cmd.encodeTimeSignal()
	}
	return blank

}

// bandwidth Reservation
func (cmd *Command) decodeBandwidthReservation(bd *bitDecoder) {
	cmd.Name = "Bandwidth Reservation"
	bd.goForward(0)
}

// private Command
func (cmd *Command) decodePrivate(bd *bitDecoder) {
	cmd.Name = "Private Command"
	cmd.Identifier = bd.uInt32(32)
	cmd.PrivateBytes = bd.asBytes(24)
}

// splice Null
func (cmd *Command) decodeSpliceNull(bd *bitDecoder) {
	cmd.Name = "Splice Null"
	bd.goForward(0)
}

// splice Insert
func (cmd *Command) decodeSpliceInsert(bd *bitDecoder) {
	cmd.Name = "Splice Insert"
	cmd.SpliceEventID = bd.uInt32(32)
	cmd.SpliceEventCancelIndicator = bd.asFlag()
	bd.goForward(7)
	cmd.OutOfNetworkIndicator = bd.asFlag()
	cmd.ProgramSpliceFlag = bd.asFlag()
	cmd.DurationFlag = bd.asFlag()
	cmd.SpliceImmediateFlag = bd.asFlag()
	cmd.EventIDComplianceFlag = bd.asFlag()
	bd.goForward(3)
	if !cmd.SpliceImmediateFlag {
		cmd.spliceTime(bd)
	}
	if cmd.DurationFlag == true {
		cmd.parseBreak(bd)
	}
	cmd.UniqueProgramID = bd.uInt16(16)
	cmd.AvailNum = bd.uInt8(8)
	cmd.AvailExpected = bd.uInt8(8)
}

// encode Splice Insert Splice Command
func (cmd *Command) encodeSpliceInsert() []byte {
	be := &bitEncoder{}
	be.Add(1, 8) //bumper
	be.Add(cmd.SpliceEventID, 32)
	be.Add(cmd.SpliceEventCancelIndicator, 1)
	be.Reserve(7)
	be.Add(cmd.OutOfNetworkIndicator, 1)
	be.Add(cmd.ProgramSpliceFlag, 1)
	be.Add(cmd.DurationFlag, 1)
	be.Add(cmd.SpliceImmediateFlag, 1)
	be.Add(cmd.EventIDComplianceFlag, 1)
	be.Reserve(3)
	if !cmd.SpliceImmediateFlag {
		cmd.encodeSpliceTime(be)
	}
	if cmd.DurationFlag {
		cmd.encodeBreak(be)
	}
	be.Add(cmd.UniqueProgramID, 16)
	be.Add(cmd.AvailNum, 8)
	be.Add(cmd.AvailExpected, 8)
	return be.Bites.Bytes()[1:] // drop Bytes[0] it's just a bumper to allow leading zero values

}

func (cmd *Command) encodeBreak(be *bitEncoder) {
	be.Add(cmd.BreakAutoReturn, 1)
	be.Reserve(6)
	be.Add(cmd.BreakDuration, 33)
}

// encode PTS splice times
func (cmd *Command) encodeSpliceTime(be *bitEncoder) {
	be.Add(cmd.TimeSpecifiedFlag, 1)
	if cmd.TimeSpecifiedFlag == true {
		be.Reserve(6)
		be.Add(cmd.PTS, 33)
		return
	}
	be.Reserve(7)
}

func (cmd *Command) parseBreak(bd *bitDecoder) {
	cmd.BreakAutoReturn = bd.asFlag()
	bd.goForward(6)
	cmd.BreakDuration = bd.as90k(33)
}

func (cmd *Command) spliceTime(bd *bitDecoder) {
	cmd.TimeSpecifiedFlag = bd.asFlag()
	if cmd.TimeSpecifiedFlag {
		bd.goForward(6)
		cmd.PTS = bd.as90k(33)
	} else {
		bd.goForward(7)
	}
}

// decode Time Signal Splice Commands
func (cmd *Command) decodeTimeSignal(bd *bitDecoder) {
	cmd.Name = "Time Signal"
	cmd.spliceTime(bd)
}

// encode Time Signal Splice Commands
func (cmd *Command) encodeTimeSignal() []byte {
	be := &bitEncoder{}
	cmd.encodeSpliceTime(be)
	return be.Bites.Bytes()
}
