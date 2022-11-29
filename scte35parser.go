package cuei

/*
	*
	   Scte35Parser is for incorporating with another MPEGTS parser.
	       Usage:

	       import(
	           "github.com/futzu/cuei"
	           )
	       scte35parser := cuei.Scte35Parser{}

	           // Each time your parser/demuxer finds a SCTE-35 packet (stream type 0x86)

	           // do something like

	           cue := scte35parser.Parse(aScte35Packet)
	           if cue != nil {
	                   // do something with the cue
	           }

*
*/
type Scte35Parser struct {
	Stream
}

/*
*

	Parse accepts a pkt as input.

	    If the packet is a partial Cue, it will be stored and aggregated
	    with the next packet until complete.

	    completed packet(s) with be decoded into a Cue and returned.

*
*/
func (scte35p *Scte35Parser) Parse(pkt []byte) (cue *Cue) {
	cue = scte35p.Scte35Parse(pkt)
	if cue != nil {
		cue.Show()
		return cue
	}

	return
}
