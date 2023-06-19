package cuei

import (
	gobs "github.com/futzu/gob"
)

/*
*
SpliceCommand

	These Splice Command types are consolidated into SpliceCommand.

	     0x0: Splice Null,
	     0x5: Splice Insert,
	     0x6: Time Signal,
	     0x7: Bandwidth Reservation,
	     0xff: Private,

*
*/
type SpliceCommand struct {
	Name        string
	CommandType uint8
	Identifier  uint32 `json:",omitempty"`
	Bites       []byte `json:",omitempty"`
	//*NBin
	SpliceEventID              string  `json:",omitempty"`
	SpliceEventCancelIndicator bool    `json:",omitempty"`
	OutOfNetworkIndicator      bool    `json:",omitempty"`
	ProgramSpliceFlag          bool    `json:",omitempty"`
	DurationFlag               bool    `json:",omitempty"`
	BreakAutoReturn            bool    `json:",omitempty"`
	BreakDuration              float64 `json:",omitempty"`
	SpliceImmediateFlag        bool    `json:",omitempty"`
	ComponentCount             uint8   `json:",omitempty"`
	Components                 []uint8 `json:",omitempty"`
	UniqueProgramID            uint16  `json:",omitempty"`
	AvailNum                   uint8   `json:",omitempty"`
	AvailExpected              uint8   `json:",omitempty"`
	TimeSpecifiedFlag          bool    `json:",omitempty"`
	PTS                        float64 `json:",omitempty"`
}

// Decode returns a Command by type
func (cmd *SpliceCommand) Decode(cmdtype uint8, gob *gobs.Gob) {
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

// Encode returns a SpliceCommand values as bytes
func (cmd *SpliceCommand) Encode() []byte {
	blank := []byte{}
	switch cmd.CommandType {
	case 0x5:
		return cmd.encodeSpliceInsert()
	}
	return blank
}

// bandwidth Reservation
func (cmd *SpliceCommand) decodeBandwidthReservation(gob *gobs.Gob) {
	cmd.Name = "Bandwidth Reservation"
	gob.Forward(0)
}

// private Command
func (cmd *SpliceCommand) decodePrivate(gob *gobs.Gob) {
	cmd.Name = "Private Command"
	cmd.Identifier = gob.UInt32(32)
	cmd.Bites = gob.Bytes(24)
}

// splice Null
func (cmd *SpliceCommand) decodeSpliceNull(gob *gobs.Gob) {
	cmd.Name = "Splice Null"
	gob.Forward(0)
}

// splice Insert
func (cmd *SpliceCommand) decodeSpliceInsert(gob *gobs.Gob) {
	cmd.Name = "Splice Insert"
	cmd.SpliceEventID = gob.Hex(32)
	cmd.SpliceEventCancelIndicator = gob.Flag()
	gob.Forward(7)
	if !cmd.SpliceEventCancelIndicator {
		cmd.OutOfNetworkIndicator = gob.Flag()
		cmd.ProgramSpliceFlag = gob.Flag()
		cmd.DurationFlag = gob.Flag()
		cmd.SpliceImmediateFlag = gob.Flag()
		gob.Forward(4)
	}
	if cmd.ProgramSpliceFlag == true {
		if !cmd.SpliceImmediateFlag {
			cmd.spliceTime(gob)
		}
	} else {
		cmd.ComponentCount = gob.UInt8(8)
		var Components [256]uint8
		cmd.Components = Components[0:cmd.ComponentCount]
		for i := range cmd.Components {
			cmd.Components[i] = gob.UInt8(8)
		}
		if !cmd.SpliceImmediateFlag {
			cmd.spliceTime(gob)
		}
	}
	if cmd.DurationFlag == true {
		cmd.parseBreak(gob)
	}
	cmd.UniqueProgramID = gob.UInt16(16)
	cmd.AvailNum = gob.UInt8(8)
	cmd.AvailExpected = gob.UInt8(8)
}

func (cmd *SpliceCommand) parseBreak(gob *gobs.Gob) {
	cmd.BreakAutoReturn = gob.Flag()
	gob.Forward(6)
	cmd.BreakDuration = gob.As90k(33)
}

func (cmd *SpliceCommand) spliceTime(gob *gobs.Gob) {
	cmd.TimeSpecifiedFlag = gob.Flag()
	if cmd.TimeSpecifiedFlag {
		gob.Forward(6)
		cmd.PTS = gob.As90k(33)
	} else {
		gob.Forward(7)
	}
}

// MkSpliceInsert easy Splice Insert Splice Command creation
/**


    The args set the SpliceInsert vars.

    SpliceEventID = eventid

    if pts is 0 (default):
        SpliceImmediateFlag      True
        timeSpecifiedFlag      False

    if pts > 0:
        spliceImmediateFlag      False
        timeSpecifiedFlag       True
        PTS                  pts

    If duration is 0(default)
        DurationFlag              False

    if duration IS set:
        OutOfNetworkIndicator   True
        durationFlag           True
        BreakAutoReturn         True
        BreakDuration             duration
        PTS                  pts

    if out :
        OutOfNetworkIndicator  True

    if out is False (default):
        OutOfNetworkIndicator   False

**/
func (cmd *SpliceCommand) MkSpliceInsert(eventid string, pts float64, duration float64, out bool) []byte {
	cmd.Name = "Splice Insert"
	cmd.CommandType = 5
	cmd.SpliceEventID = eventid
	cmd.SpliceEventCancelIndicator = false
	cmd.ProgramSpliceFlag = false
	cmd.SpliceImmediateFlag = true
	cmd.DurationFlag = false
	cmd.AvailNum = 0
	cmd.AvailExpected = 0
	cmd.TimeSpecifiedFlag = false
	if pts > 0 {
		cmd.ProgramSpliceFlag = true
		cmd.SpliceImmediateFlag = false
		cmd.TimeSpecifiedFlag = true
		cmd.PTS = pts
	} else {
		cmd.ComponentCount = 0
	}
	if duration > 0 {
		cmd.DurationFlag = true
		cmd.BreakDuration = duration
		cmd.BreakAutoReturn = true
		cmd.OutOfNetworkIndicator = true
	}
	return cmd.Encode()
}

// encode Splice Insert Splice Command
func (cmd *SpliceCommand) encodeSpliceInsert() []byte {
	nb := &Nbin{}
	nb.AddHex64(cmd.SpliceEventID, 32)
	nb.AddFlag(cmd.SpliceEventCancelIndicator)
	nb.Reserve(7)
	if !cmd.SpliceEventCancelIndicator {
		nb.AddFlag(cmd.OutOfNetworkIndicator)
		nb.AddFlag(cmd.ProgramSpliceFlag)
		nb.AddFlag(cmd.DurationFlag)
		nb.AddFlag(cmd.SpliceImmediateFlag)
		nb.Reserve(4)
	}
	if cmd.ProgramSpliceFlag == true {
		if !cmd.SpliceImmediateFlag {
			cmd.encodeSpliceTime(nb)
		}
	} else {
		nb.Add8(cmd.ComponentCount, 8)
		for i := range cmd.Components {
			nb.Add8(cmd.Components[i], 8)
		}
		if !cmd.SpliceImmediateFlag {
			cmd.encodeSpliceTime(nb)
		}
	}
	if cmd.DurationFlag == true {
		cmd.encodeBreak(nb)
	}
	nb.Add16(cmd.UniqueProgramID, 16)
	nb.Add8(cmd.AvailNum, 8)
	nb.Add8(cmd.AvailExpected, 8)
	cmd.Bites = nb.Bites.Bytes()
	return nb.Bites.Bytes()
}

func (cmd *SpliceCommand) encodeBreak(nb *Nbin) {
	nb.AddFlag(cmd.BreakAutoReturn)
	nb.Reserve(6)
	nb.Add90k(cmd.BreakDuration, 33)
}

// encode PTS splice times
func (cmd *SpliceCommand) encodeSpliceTime(nb *Nbin) {
	nb.AddFlag(cmd.TimeSpecifiedFlag)
	if cmd.TimeSpecifiedFlag {
		nb.Reserve(6)
		nb.Add90k(cmd.PTS, 33)
	} else {
		nb.Reserve(7)
	}
}

// decode Time Signal Splice Commands
func (cmd *SpliceCommand) decodeTimeSignal(gob *gobs.Gob) {
	cmd.Name = "Time Signal"
	cmd.spliceTime(gob)
}

// encode Time Signal Splice Commands
func (cmd *SpliceCommand) encodeTimeSignal() []byte {
	nb := &Nbin{}
	cmd.encodeSpliceTime(nb)
	cmd.Bites = nb.Bites.Bytes()
	return nb.Bites.Bytes()
}
