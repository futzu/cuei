
 `Splice Insert Splice Command as defined by GOTS`

<details> <summary> Line count 186</summary>

```go
/ spliceInsert is a struct that represents a splice insert command in SCTE35
type spliceInsert struct {
	eventID               uint32
	eventCancelIndicator  bool
	outOfNetworkIndicator bool

	isProgramSplice bool
	spliceImmediate bool

	hasPTS bool
	pts    gots.PTS

	components []Component

	hasDuration     bool
	duration        gots.PTS
	autoReturn      bool
	uniqueProgramId uint16
	availNum        uint8
	availsExpected  uint8
}

// CommandType returns the signal's splice command type value.
func (c *spliceInsert) CommandType() SpliceCommandType {
	return SpliceInsert
}

// parseSpliceInsert extracts a splice_insert() command from a bytes buffer.
// It returns a spliceInsert describing the command.
func parseSpliceInsert(buf *bytes.Buffer) (*spliceInsert, error) {
	cmd := &spliceInsert{}
	if err := cmd.parse(buf); err != nil {
		return nil, err
	}
	return cmd, nil
}

// parse will parse bytes in the form of bytes.Buffer into a splice insert struct
func (c *spliceInsert) parse(buf *bytes.Buffer) error {
	baseFields := buf.Next(5)
	if len(baseFields) < 5 { // length of required fields
		return gots.ErrInvalidSCTE35Length
	}
	c.eventID = binary.BigEndian.Uint32(baseFields[:4])
	// splice_event_cancel_indicator 1
	// reserved 7
	c.eventCancelIndicator = baseFields[4]&0x80 == 0x80
	if c.eventCancelIndicator {
		return nil
	}
	// out_of_network_indicator 1
	// program_splice_flag 1
	// duration_flag 1
	// splice_immediate_flag 1
	// reserved 4
	flags, err := buf.ReadByte()
	if err != nil {
		return gots.ErrInvalidSCTE35Length
	}
	c.outOfNetworkIndicator = flags&0x80 == 0x80
	c.isProgramSplice = flags&0x40 == 0x40
	c.hasDuration = flags&0x20 == 0x20
	c.spliceImmediate = flags&0x10 == 0x10

	if c.isProgramSplice && !c.spliceImmediate {
		hasPTS, pts, err := parseSpliceTime(buf)
		if err != nil {
			return err
		}
		if !hasPTS {
			return gots.ErrSCTE35UnsupportedSpliceCommand
		}
		c.hasPTS = hasPTS
		c.pts = pts
	}
	if !c.isProgramSplice {
		cc, err := buf.ReadByte()
		if err != nil {
			return gots.ErrInvalidSCTE35Length
		}
		// read components
		for ; cc > 0; cc-- {
			// component_tag
			tag, err := buf.ReadByte()
			if err != nil {
				return gots.ErrInvalidSCTE35Length
			}
			comp := &component{componentTag: tag}
			if !c.spliceImmediate {
				hasPts, pts, err := parseSpliceTime(buf)
				if err != nil {
					return err
				}
				comp.hasPts = hasPts
				comp.pts = pts
			}
			c.components = append(c.components, comp)
		}
	}
	if c.hasDuration {
		data := buf.Next(5)
		if len(data) < 5 {
			return gots.ErrInvalidSCTE35Length
		}
		// break_duration() structure:
		c.autoReturn = data[0]&0x80 == 0x80
		c.duration = uint40(data) & 0x01ffffffff
	}
	progInfo := buf.Next(4)
	if len(progInfo) < 4 {
		return gots.ErrInvalidSCTE35Length
	}
	c.uniqueProgramId = binary.BigEndian.Uint16(progInfo[:2])
	c.availNum = progInfo[2]
	c.availsExpected = progInfo[3]
	return nil
}

// EventID returns the event id
func (c *spliceInsert) EventID() uint32 {
	return c.eventID
}

// IsOut returns the value of the out of network indicator
func (c *spliceInsert) IsOut() bool {
	return c.outOfNetworkIndicator
}

// IsEventCanceled returns the event cancel indicator
func (c *spliceInsert) IsEventCanceled() bool {
	return c.eventCancelIndicator
}

// HasPTS returns true if there is a pts time on the command.
func (c *spliceInsert) HasPTS() bool {
	return c.hasPTS
}

// PTS returns the PTS time of the command, not including adjustment.
func (c *spliceInsert) PTS() gots.PTS {
	return c.pts
}

// HasDuration returns true if there is a duration
func (c *spliceInsert) HasDuration() bool {
	return c.hasDuration
}

// Duration returns the PTS duration of the command
func (c *spliceInsert) Duration() gots.PTS {
	return c.duration
}

// Components returns the components of the splice command
func (c *spliceInsert) Components() []Component {
	return c.components
}

// IsAutoReturn returns the boolean value of the auto return field
func (c *spliceInsert) IsAutoReturn() bool {
	return c.autoReturn
}

// UniqueProgramId returns the unique_program_id field
func (c *spliceInsert) UniqueProgramId() uint16 {
	return c.uniqueProgramId
}

// AvailNum returns the avail_num field, index of this avail or zero if unused
func (c *spliceInsert) AvailNum() uint8 {
	return c.availNum
}

// AvailsExpected returns avails_expected field, number of avails for program
func (c *spliceInsert) AvailsExpected() uint8 {
	return c.availsExpected
}

// IsProgramSplice returns if the program_splice_flag is set
func (c *spliceInsert) IsProgramSplice() bool {
	return c.isProgramSplice
}

// SpliceImmediate returns if the splice_immediate_flag is set
func (c *spliceInsert) SpliceImmediate() bool {
	return c.spliceImmediate
}
```
  </details>

 `Splice Insert, Splice Null, Time Signal ,Bandwidth Reservation,and Private Splice Commands as defined by cuei. `

  <details><summary> Total lines 133</summary>


```go

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


//Decode returns a Command by type
func (cmd *SpliceCommand) Decode(cmdtype uint8, gob *gobs.Gob) {
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
func (cmd *SpliceCommand) bandwidthReservation(gob *gobs.Gob) {
	cmd.Name = "Bandwidth Reservation"
	gob.Forward(0)
}

// private Command
func (cmd *SpliceCommand) private(gob *gobs.Gob) {
	cmd.Name = "Private Command"
	cmd.Identifier = gob.UInt32(32)
	cmd.Bites = gob.Bytes(24)
}

// splice Null
func (cmd *SpliceCommand) spliceNull(gob *gobs.Gob) {
	cmd.Name = "Splice Null"
	gob.Forward(0)
}

// splice Insert
func (cmd *SpliceCommand) spliceInsert(gob *gobs.Gob) {
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

// time Signal
func (cmd *SpliceCommand) timeSignal(gob *gobs.Gob) {
	cmd.Name = "Time Signal"
	cmd.spliceTime(gob)
}
```
    
  </details>
