package cuei

import (
	bitter "github.com/futzu/bitter"
)

/*
*
Command

	These Splice Command types are consolidated into Command.

	     0x0: Splice Null,
	     0x5: Splice Insert,
	     0x6: Time Signal,
	     0x7: Bandwidth Reservation,
	     0xff: Private,

*
*/
type Command struct {
	Name                       string
	CommandType                uint8
	PrivateBytes               []byte  //`json:",omitempty"`
	Identifier                 uint32  // `json:",omitempty"`
	SpliceEventID              uint32  // `json:",omitempty"`
	SpliceEventCancelIndicator bool    //   `json:",omitempty"`
	OutOfNetworkIndicator      bool    // `json:",omitempty"`
	ProgramSpliceFlag          bool    //`json:",omitempty"`
	DurationFlag               bool    // `json:",omitempty"`
	BreakAutoReturn            bool    //  `json:",omitempty"`
	BreakDuration              float64 // `json:",omitempty"`
	SpliceImmediateFlag        bool    //  `json:",omitempty"`
	UniqueProgramID            uint16  // `json:",omitempty"`
	AvailNum                   uint8   //`json:",omitempty"`
	AvailExpected              uint8   // `json:",omitempty"`
	TimeSpecifiedFlag          bool    //   `json:",omitempty"`
	PTS                        float64 // `json:",omitempty"`
}

// Decode a Splice Command
func (cmd *Command) Decode(cmdtype uint8, bd *bitter.Decoder) {
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
func (cmd *Command) Encode() []byte {
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
func (cmd *Command) decodeBandwidthReservation(bd *bitter.Decoder) {
	cmd.Name = "Bandwidth Reservation"
	bd.Forward(0)
}

// private Command
func (cmd *Command) decodePrivate(bd *bitter.Decoder) {
	cmd.Name = "Private Command"
	cmd.Identifier = bd.UInt32(32)
	cmd.PrivateBytes = bd.Bytes(24)
}

// splice Null
func (cmd *Command) decodeSpliceNull(bd *bitter.Decoder) {
	cmd.Name = "Splice Null"
	bd.Forward(0)
}

// splice Insert
func (cmd *Command) decodeSpliceInsert(bd *bitter.Decoder) {
	cmd.Name = "Splice Insert"
	cmd.SpliceEventID = bd.UInt32(32)
	cmd.SpliceEventCancelIndicator = bd.Flag()
	bd.Forward(7)
	cmd.OutOfNetworkIndicator = bd.Flag()
	cmd.ProgramSpliceFlag = bd.Flag()
	cmd.DurationFlag = bd.Flag()
	cmd.SpliceImmediateFlag = bd.Flag()
	bd.Forward(4)
	if cmd.SpliceImmediateFlag == false {
		cmd.spliceTime(bd)
	}
	if cmd.DurationFlag == true {
		cmd.parseBreak(bd)
	}
	cmd.UniqueProgramID = bd.UInt16(16)
	cmd.AvailNum = bd.UInt8(8)
	cmd.AvailExpected = bd.UInt8(8)
}

// encode Splice Insert Splice Command
func (cmd *Command) encodeSpliceInsert() []byte {
	be := &bitter.Encoder{}
	be.Add8(1, 8) //bumper
	be.Add32(cmd.SpliceEventID, 32)
	be.AddFlag(cmd.SpliceEventCancelIndicator)
	be.Reserve(7)
	be.AddFlag(cmd.OutOfNetworkIndicator)
	be.AddFlag(cmd.ProgramSpliceFlag)
	be.AddFlag(cmd.DurationFlag)
	be.AddFlag(cmd.SpliceImmediateFlag)
	be.Reserve(4)
	if cmd.SpliceImmediateFlag == false {
		cmd.encodeSpliceTime(be)
	}
	if cmd.DurationFlag == true {
		cmd.encodeBreak(be)
	}
	be.Add16(cmd.UniqueProgramID, 16)
	be.Add8(cmd.AvailNum, 8)
	be.Add8(cmd.AvailExpected, 8)
	return be.Bites.Bytes()[1:] // drop Bytes[0] it's just a bumper to allow leading zero values

}

func (cmd *Command) encodeBreak(be *bitter.Encoder) {
	be.AddFlag(cmd.BreakAutoReturn)
	be.Reserve(6)
	be.Add90k(cmd.BreakDuration, 33)
}

// encode PTS splice times
func (cmd *Command) encodeSpliceTime(be *bitter.Encoder) {
	be.AddFlag(cmd.TimeSpecifiedFlag)
	if cmd.TimeSpecifiedFlag == true {
		be.Reserve(6)
		be.Add90k(cmd.PTS, 33)
		return
	}
	be.Reserve(7)
}

func (cmd *Command) parseBreak(bd *bitter.Decoder) {
	cmd.BreakAutoReturn = bd.Flag()
	bd.Forward(6)
	cmd.BreakDuration = bd.As90k(33)
}

func (cmd *Command) spliceTime(bd *bitter.Decoder) {
	cmd.TimeSpecifiedFlag = bd.Flag()
	if cmd.TimeSpecifiedFlag {
		bd.Forward(6)
		cmd.PTS = bd.As90k(33)
	} else {
		bd.Forward(7)
	}
}

// decode Time Signal Splice Commands
func (cmd *Command) decodeTimeSignal(bd *bitter.Decoder) {
	cmd.Name = "Time Signal"
	cmd.spliceTime(bd)
}

// encode Time Signal Splice Commands
func (cmd *Command) encodeTimeSignal() []byte {
	be := &bitter.Encoder{}
	cmd.encodeSpliceTime(be)
	return be.Bites.Bytes()
}
