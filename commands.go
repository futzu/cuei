package cuei

import (
	gobs "github.com/futzu/gob"
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
func (cmd *Command) Decode(cmdtype uint8, gob *gobs.Gob) {
	cmd.CommandType = cmdtype
	switch cmdtype {
	case 0x0:
		cmd.decodeSpliceNull(gob)
	case 0x5:
		cmd.decodeSpliceInsert(gob)
	case 0x6:
		cmd.decodeTimeSignal(gob)
	case 0x7:
		cmd.decodeBandwidthReservation(gob)
	case 0xff:
		cmd.decodePrivate(gob)
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
func (cmd *Command) decodeBandwidthReservation(gob *gobs.Gob) {
	cmd.Name = "Bandwidth Reservation"
	gob.Forward(0)
}

// private Command
func (cmd *Command) decodePrivate(gob *gobs.Gob) {
	cmd.Name = "Private Command"
	cmd.Identifier = gob.UInt32(32)
	cmd.PrivateBytes = gob.Bytes(24)
}

// splice Null
func (cmd *Command) decodeSpliceNull(gob *gobs.Gob) {
	cmd.Name = "Splice Null"
	gob.Forward(0)
}

// splice Insert
func (cmd *Command) decodeSpliceInsert(gob *gobs.Gob) {
	cmd.Name = "Splice Insert"
	cmd.SpliceEventID = gob.UInt32(32)
	cmd.SpliceEventCancelIndicator = gob.Flag()
	gob.Forward(7)
	cmd.OutOfNetworkIndicator = gob.Flag()
	cmd.ProgramSpliceFlag = gob.Flag()
	cmd.DurationFlag = gob.Flag()
	cmd.SpliceImmediateFlag = gob.Flag()
	gob.Forward(4)
	if cmd.SpliceImmediateFlag == false {
		cmd.spliceTime(gob)
	}
	if cmd.DurationFlag == true {
		cmd.parseBreak(gob)
	}
	cmd.UniqueProgramID = gob.UInt16(16)
	cmd.AvailNum = gob.UInt8(8)
	cmd.AvailExpected = gob.UInt8(8)
}

// encode Splice Insert Splice Command
func (cmd *Command) encodeSpliceInsert() []byte {
	nb := &Nbin{}
	nb.Add8(1, 8) //bumper
	nb.Add32(cmd.SpliceEventID, 32)
	nb.AddFlag(cmd.SpliceEventCancelIndicator)
	nb.Reserve(7)
	nb.AddFlag(cmd.OutOfNetworkIndicator)
	nb.AddFlag(cmd.ProgramSpliceFlag)
	nb.AddFlag(cmd.DurationFlag)
	nb.AddFlag(cmd.SpliceImmediateFlag)
	nb.Reserve(4)
	if cmd.SpliceImmediateFlag == false {
		cmd.encodeSpliceTime(nb)
	}
	if cmd.DurationFlag == true {
		cmd.encodeBreak(nb)
	}
	nb.Add16(cmd.UniqueProgramID, 16)
	nb.Add8(cmd.AvailNum, 8)
	nb.Add8(cmd.AvailExpected, 8)
	return nb.Bites.Bytes()[1:] // drop Bytes[0] it's just a bumper to allow leading zero values

}

func (cmd *Command) encodeBreak(nb *Nbin) {
	nb.AddFlag(cmd.BreakAutoReturn)
	nb.Reserve(6)
	nb.Add90k(cmd.BreakDuration, 33)
}

// encode PTS splice times
func (cmd *Command) encodeSpliceTime(nb *Nbin) {
	nb.AddFlag(cmd.TimeSpecifiedFlag)
	if cmd.TimeSpecifiedFlag == true {
		nb.Reserve(6)
		nb.Add90k(cmd.PTS, 33)
		return
	}
	nb.Reserve(7)
}

func (cmd *Command) parseBreak(gob *gobs.Gob) {
	cmd.BreakAutoReturn = gob.Flag()
	gob.Forward(6)
	cmd.BreakDuration = gob.As90k(33)
}

func (cmd *Command) spliceTime(gob *gobs.Gob) {
	cmd.TimeSpecifiedFlag = gob.Flag()
	if cmd.TimeSpecifiedFlag {
		gob.Forward(6)
		cmd.PTS = gob.As90k(33)
	} else {
		gob.Forward(7)
	}
}

// decode Time Signal Splice Commands
func (cmd *Command) decodeTimeSignal(gob *gobs.Gob) {
	cmd.Name = "Time Signal"
	cmd.spliceTime(gob)
}

// encode Time Signal Splice Commands
func (cmd *Command) encodeTimeSignal() []byte {
	nb := &Nbin{}
	cmd.encodeSpliceTime(nb)
	return nb.Bites.Bytes()
}
