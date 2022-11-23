package cuei

import (
	goober "github.com/futzu/gob"
)

type SpliceCommand struct {
	Name                       string
	CommandType                uint8
	Identifier                 uint32  `json:",omitempty"`
	Bites                      []byte  `json:",omitempty"`
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

/**
Decode returns a Command by cmdtype

    Supported Commands:
    	0x0: Splice Null,
    	0x5: Splice Insert,
    	0x6: Time Signal,
    	0x7: Bandwidth Reservation,
    	0xff: Private,
**/
func (cmd *SpliceCommand) Decode(cmdtype uint8, gob *goober.Gob) {
	cmd.CommandType = cmdtype
	switch cmdtype {
	case 0x0:
		cmd.spliceNull(gob)
	case 0x5:
		cmd.spliceInsert(gob)
	case 0x6:
		cmd.timeSignal(gob)
	case 0x7:
		cmd.bandwidthReservation(gob)
	case 0xff:
		cmd.private(gob)
	}

}

// bandwidth Reservation
func (cmd *SpliceCommand) bandwidthReservation(gob *goober.Gob) {
	cmd.Name = "Bandwidth Reservation"
	gob.Forward(0)
}

// private Command
func (cmd *SpliceCommand) private(gob *goober.Gob) {
	cmd.Name = "Private Command"
	cmd.Identifier = gob.UInt32(32)
	cmd.Bites = gob.Bytes(24)
}

// splice Null
func (cmd *SpliceCommand) spliceNull(gob *goober.Gob) {
	cmd.Name = "Splice Null"
	gob.Forward(0)
}

// splice Insert
func (cmd *SpliceCommand) spliceInsert(gob *goober.Gob) {
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

func (cmd *SpliceCommand) parseBreak(gob *goober.Gob) {
	cmd.BreakAutoReturn = gob.Flag()
	gob.Forward(6)
	cmd.BreakDuration = gob.As90k(33)
}

func (cmd *SpliceCommand) spliceTime(gob *goober.Gob) {
	cmd.TimeSpecifiedFlag = gob.Flag()
	if cmd.TimeSpecifiedFlag {
		gob.Forward(6)
		cmd.PTS = gob.As90k(33)
	} else {
		gob.Forward(7)
	}
}

// time Signal
func (cmd *SpliceCommand) timeSignal(gob *goober.Gob) {
	cmd.Name = "Time Signal"
	cmd.spliceTime(gob)
}
