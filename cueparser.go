package cuei

// Cue Parser is wrapper for Cue
type CueParser struct {
	Cue
}

// Parse decodes bites into SCTE-35 values
func (cuep *CueParser) Parse(bites []byte)  {
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
