package cuei

// StreamParser parses a []byte of mpegts for SCTE-35
type StreamParser struct {
	Stream
}

// Parse parses  mpegts bytes and returns any SCTE-35 Cues found
func (streamp *StreamParser) Parse(bites []byte) []*Cue {
	cues := streamp.DecodeBytes(bites)
	return cues
}

// initialize and return a *StreamParser
func NewStreamParser() *StreamParser {
	sp := &StreamParser{}
	sp.Pids = &Pids{}
	sp.mkMaps()
	return sp
}
