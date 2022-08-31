package cuei

import (
//    "github.com/futzu/gob"
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

// CommandDecoder returns a Command by cmdtype
func (cmd *SpliceCommand) Decoder(cmdtype uint8, gob *Gob) {
	cmd.CommandType = cmdtype
	switch cmdtype {
	case 0:
		cmd.SpliceNull(gob)
	case 5:
		cmd.SpliceInsert(gob)
	case 6:
		cmd.TimeSignal(gob)
	case 7:
		cmd.BandwidthReservation(gob)
	case 255:
		cmd.Private(gob)
	}

}

// Bandwidth Reservation
func (cmd *SpliceCommand) BandwidthReservation(gob *Gob) {
	cmd.Name = "Bandwidth Reservation"
	gob.Forward(0)
}

// Private Command
func (cmd *SpliceCommand) Private(gob *Gob) {
	cmd.Name = "Private Command"
	cmd.Identifier = gob.uint32(32)
	cmd.Bites = gob.Bytes(24)
}

// Splice Null
func (cmd *SpliceCommand) SpliceNull(gob *Gob) {
	cmd.Name = "Splice Null"
	gob.Forward(0)
}

// Splice Insert
func (cmd *SpliceCommand) SpliceInsert(gob *Gob) {
	cmd.Name = "Splice Insert"
	cmd.SpliceEventID = gob.Hex(32)
	cmd.SpliceEventCancelIndicator = gob.Bool()
	gob.Forward(7)
	if !cmd.SpliceEventCancelIndicator {
		cmd.OutOfNetworkIndicator = gob.Bool()
		cmd.ProgramSpliceFlag = gob.Bool()
		cmd.DurationFlag = gob.Bool()
		cmd.SpliceImmediateFlag = gob.Bool()
		gob.Forward(4)
	}
	if cmd.ProgramSpliceFlag == true {
		if !cmd.SpliceImmediateFlag {
			cmd.spliceTime(gob)
		}
	} else {
		cmd.ComponentCount = gob.uint8(8)
		var Components [256]uint8
		cmd.Components = Components[0:cmd.ComponentCount]
		for i := range cmd.Components {
			cmd.Components[i] = gob.uint8(8)
		}
		if !cmd.SpliceImmediateFlag {
			cmd.spliceTime(gob)
		}
	}
	if cmd.DurationFlag == true {
		cmd.parseBreak(gob)
	}
	cmd.UniqueProgramID = gob.uint16(16)
	cmd.AvailNum = gob.uint8(8)
	cmd.AvailExpected = gob.uint8(8)
}

func (cmd *SpliceCommand) parseBreak(gob *Gob) {
	cmd.BreakAutoReturn = gob.Bool()
	gob.Forward(6)
	cmd.BreakDuration = gob.As90k(33)
}

func (cmd *SpliceCommand) spliceTime(gob *Gob) {
	cmd.TimeSpecifiedFlag = gob.Bool()
	if cmd.TimeSpecifiedFlag {
		gob.Forward(6)
		cmd.PTS = gob.As90k(33)
	} else {
		gob.Forward(7)
	}
}

// Time Signal
func (cmd *SpliceCommand) TimeSignal(gob *Gob) {
	cmd.Name = "Time Signal"
	cmd.spliceTime(gob)
}
