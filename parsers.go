package cuei

// Cue Parser is wrapper for Cue
type CueParser struct {
	Cue
}

// Parse decodes bites into SCTE-35 values
func (cuep *CueParser) Parse(bites []byte) {
	cuep.Decode(bites)
}

// Show prints the Cue values as JSON
func (cuep *CueParser) Show() {
	cuep.Show()
}

// initialize and return a *CueParser
func NewCueParser() *CueParser {
	cuep := &CueParser{}
	return cuep
}

// StreamParser parses a []byte of mpegts for SCTE-35
type StreamParser struct {
	Stream
}

// Parse parses  mpegts bytes and returns any SCTE-35 Cues found
func (streamp *StreamParser) Parse(bites []byte) []*Cue {
	cues := streamp.decodeBytes(bites)
	return cues
}

// ParseFile parses a mpegts file and returns any SCTE-35 Cues found
func (streamp *StreamParser) ParseFile(filename string) []*Cue {
	cues := streamp.decode(filename)
	return cues
}

// initialize and return a *StreamParser
func NewStreamParser() *StreamParser {
	sp := &StreamParser{}
	sp.Pids = &Pids{}
	sp.mkMaps()
	return sp
}
