package cuei

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
	PrivateBytes               []byte  `json:",omitempty"`
	Identifier                 uint32  `json:",omitempty"`
	SpliceEventID              uint32  `json:",omitempty"`
	SpliceEventCancelIndicator bool    `json:",omitempty"`
	OutOfNetworkIndicator      bool    `json:",omitempty"`
	ProgramSpliceFlag          bool    `json:",omitempty"`
	DurationFlag               bool    `json:",omitempty"`
	BreakAutoReturn            bool    `json:",omitempty"`
	BreakDuration              float64 `json:",omitempty"`
	SpliceImmediateFlag        bool    `json:",omitempty"`
	UniqueProgramID            uint16  `json:",omitempty"`
	AvailNum                   uint8   `json:",omitempty"`
	AvailExpected              uint8   `json:",omitempty"`
	TimeSpecifiedFlag          bool    `json:",omitempty"`
	PTS                        float64 `json:",omitempty"`
}

// Decode a Splice Command
func (cmd *Command) Decode(cmdtype uint8, bd *BitDecoder) {
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
func (cmd *Command) decodeBandwidthReservation(bd *BitDecoder) {
	cmd.Name = "Bandwidth Reservation"
	bd.goForward(0)
}

// private Command
func (cmd *Command) decodePrivate(bd *BitDecoder) {
	cmd.Name = "Private Command"
	cmd.Identifier = bd.uInt32(32)
	cmd.PrivateBytes = bd.asBytes(24)
}

// splice Null
func (cmd *Command) decodeSpliceNull(bd *BitDecoder) {
	cmd.Name = "Splice Null"
	bd.goForward(0)
}

// splice Insert
func (cmd *Command) decodeSpliceInsert(bd *BitDecoder) {
	cmd.Name = "Splice Insert"
	cmd.SpliceEventID = bd.uInt32(32)
	cmd.SpliceEventCancelIndicator = bd.asFlag()
	bd.goForward(7)
	cmd.OutOfNetworkIndicator = bd.asFlag()
	cmd.ProgramSpliceFlag = bd.asFlag()
	cmd.DurationFlag = bd.asFlag()
	cmd.SpliceImmediateFlag = bd.asFlag()
	bd.goForward(4)
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
	be := &BitEncoder{}
	be.Add(1, 8) //bumper
	be.Add(cmd.SpliceEventID, 32)
	be.Add(cmd.SpliceEventCancelIndicator, 1)
	be.Reserve(7)
	be.Add(cmd.OutOfNetworkIndicator, 1)
	be.Add(cmd.ProgramSpliceFlag, 1)
	be.Add(cmd.DurationFlag, 1)
	be.Add(cmd.SpliceImmediateFlag, 1)
	be.Reserve(4)
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

func (cmd *Command) encodeBreak(be *BitEncoder) {
	be.Add(cmd.BreakAutoReturn, 1)
	be.Reserve(6)
	be.Add(cmd.BreakDuration, 33)
}

// encode PTS splice times
func (cmd *Command) encodeSpliceTime(be *BitEncoder) {
	be.Add(cmd.TimeSpecifiedFlag, 1)
	if cmd.TimeSpecifiedFlag == true {
		be.Reserve(6)
		be.Add(cmd.PTS, 33)
		return
	}
	be.Reserve(7)
}

func (cmd *Command) parseBreak(bd *BitDecoder) {
	cmd.BreakAutoReturn = bd.asFlag()
	bd.goForward(6)
	cmd.BreakDuration = bd.as90k(33)
}

func (cmd *Command) spliceTime(bd *BitDecoder) {
	cmd.TimeSpecifiedFlag = bd.asFlag()
	if cmd.TimeSpecifiedFlag {
		bd.goForward(6)
		cmd.PTS = bd.as90k(33)
	} else {
		bd.goForward(7)
	}
}

// decode Time Signal Splice Commands
func (cmd *Command) decodeTimeSignal(bd *BitDecoder) {
	cmd.Name = "Time Signal"
	cmd.spliceTime(bd)
}

// encode Time Signal Splice Commands
func (cmd *Command) encodeTimeSignal() []byte {
	be := &BitEncoder{}
	cmd.encodeSpliceTime(be)
	return be.Bites.Bytes()
}
