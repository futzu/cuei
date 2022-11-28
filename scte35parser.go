package cuei

// Scte35Parser is for incorporating with another MPEGTS parser.
type Scte35Parser struct {
	Stream
}

/*
	*
	Parse accepts a pkt as input.

	If the packet is a partial Cue,
	it will be stored and aggregated
	with the next packet until complete.

	Completed packet(s)
	with be decoded into a Cue and returned.

	Usage:

	 import(
	   "github.com/futzu/cuei"
	 )
	 scte35parser := Scte35Parser{}
	 // parse mpegts for scte35 packets with your MPEGTS parser.
	 ...


	 cue := scte35parser.Parse(aScte35Packet)
	 if cue != nil {
	     // do something with the cue
	 }

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
