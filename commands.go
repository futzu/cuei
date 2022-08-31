package cuei

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
    
}

// CommandDecoder returns a Command by cmdtype
func (cmd *SpliceCommand) Decoder(cmdtype uint8,gob *gob.Gob) {
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
func (cmd *SpliceCommand) BandwidthReservation(gob *gob.Gob) {
	cmd.Name = "Bandwidth Reservation"
	gob.Forward(0)
}
// Private Command
func (cmd *SpliceCommand) Private(gob *gob.Gob) {
	cmd.Name = "Private Command"
	cmd.Identifier = gob.UInt32(32)
	cmd.Bites = gob.AsBytes(24)
}

// Splice Null
func (cmd *SpliceCommand) SpliceNull(gob *gob.Gob) {
	cmd.Name = "Splice Null"
	gob.Forward(0)
}
