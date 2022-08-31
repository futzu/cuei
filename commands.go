package cuei

import (
    "github.com/futzu/gob"
)

type SpliceCommand struct {

    Name                        string  
    CommandType                 uint8
	Identifier                  uint32  `json:",omitempty"`
	Bites                       []byte  `json:",omitempty"`
    SpliceEventID               string  `json:",omitempty"`
	SpliceEventCancelIndicator  bool    `json:",omitempty"`
	OutOfNetworkIndicator       bool    `json:",omitempty"`
	ProgramSpliceFlag           bool    `json:",omitempty"`
	DurationFlag                bool    `json:",omitempty"`
	BreakAutoReturn             bool    `json:",omitempty"`
	BreakDuration               float64 `json:",omitempty"`
	SpliceImmediateFlag         bool    `json:",omitempty"`
	ComponentCount              uint8   `json:",omitempty"`
	Components                  []uint8 `json:",omitempty"`
	UniqueProgramID             uint16  `json:",omitempty"`
	AvailNum                    uint8   `json:",omitempty"`
	AvailExpected               uint8   `json:",omitempty"`
	TimeSpecifiedFlag           bool    `json:",omitempty"`
	PTS                         float64 `json:",omitempty"`
    gob                          *gob.Gob
}

// CommandDecoder returns a Command by cmdtype
func (cmd *SpliceCommand) Decoder(cmdtype uint8,gob *gob.Gob) {
    cmd.gob = gob
    //cmd.CommandType = cmdtype
	cmdmap := map[uint8]func() {
	    0:cmd.SpliceNull,
	    5:cmd.SpliceInsert,
	    6:cmd.TimeSignal,
        7:cmd.BandwidthReservation,
        255:cmd.Private,
	}
    fn, ok := cmdmap[cmd.CommandType]
	if ok {
		fn()
	}
	
}


// Bandwidth Reservation
func (cmd *SpliceCommand) BandwidthReservation() {
	cmd.Name = "Bandwidth Reservation"
	cmd.gob.Forward(0)
}
// Private Command
func (cmd *SpliceCommand) Private() {
	cmd.Name = "Private Command"
	cmd.Identifier = cmd.gob.UInt32(32)
	cmd.Bites = cmd.gob.Bytes(24)
}

// Splice Null
func (cmd *SpliceCommand) SpliceNull() {
	cmd.Name = "Splice Null"
}

// Splice Insert
func (cmd *SpliceCommand) SpliceInsert() {
	cmd.Name = "Splice Insert"
	cmd.SpliceEventID = cmd.gob.Hex(32)
	cmd.SpliceEventCancelIndicator = cmd.gob.Bool()
	cmd.gob.Forward(7)
	if !cmd.SpliceEventCancelIndicator {
		cmd.OutOfNetworkIndicator = cmd.gob.Bool()
		cmd.ProgramSpliceFlag = cmd.gob.Bool()
		cmd.DurationFlag = cmd.gob.Bool()
		cmd.SpliceImmediateFlag = cmd.gob.Bool()
		cmd.gob.Forward(4)
	}
	if cmd.ProgramSpliceFlag == true {
		if !cmd.SpliceImmediateFlag {
			cmd.spliceTime()
		}
	} else {
		cmd.ComponentCount = cmd.gob.UInt8(8)
		var Components [256]uint8
		cmd.Components = Components[0:cmd.ComponentCount]
		for i := range cmd.Components {
			cmd.Components[i] = cmd.gob.UInt8(8)
		}
		if !cmd.SpliceImmediateFlag {
			cmd.spliceTime()
		}
	}
	if cmd.DurationFlag == true {
		cmd.parseBreak()
	}
	cmd.UniqueProgramID = cmd.gob.UInt16(16)
	cmd.AvailNum = cmd.gob.UInt8(8)
	cmd.AvailExpected = cmd.gob.UInt8(8)
}

func (cmd *SpliceCommand) parseBreak() {
	cmd.BreakAutoReturn = cmd.gob.Bool()
	cmd.gob.Forward(6)
	cmd.BreakDuration = cmd.gob.As90k(33)
}

func (cmd *SpliceCommand) spliceTime() {
	cmd.TimeSpecifiedFlag = cmd.gob.Bool()
	if cmd.TimeSpecifiedFlag {
		cmd.gob.Forward(6)
		cmd.PTS = cmd.gob.As90k(33)
	} else {
		cmd.gob.Forward(7)
	}
}

// Time Signal
func (cmd *SpliceCommand) TimeSignal() {
	cmd.Name = "Time Signal"
	cmd.spliceTime()
}
