package cuei

/*
Scte35Parser is for incorporating with another MPEGTS parser.

	Usage:

	     import(
	         "fmt"
	         "github.com/futzu/cuei"
	     )

	     scte35parser := cuei.NewScte35Parser()

	 // Each time your parser/demuxer finds
	 // a SCTE-35 packet do something like

	     cue := scte35parser.Parse(aScte35Packet)
	     if cue != nil {
	            // process the Cue
	            fmt.Printf("%#v",cue.Command)
	     }
*/
type Scte35Parser struct {
	Stream
}

/*
		Parse accepts a MPEGTS SCTE-35 packet as input.


	        If the MPEGTS SCTE-35 packet contains a complete cue message

	            The cue message is decoded into a Cue and returned.


		    If the MPEGTS SCTE-35 packet is a partial cue message

	            It will be stored and aggregated with the next MPEGTS SCTE-35 packet until complete.

	            Completed cue messages are decoded into a Cue and returned.
*/
func (scte35p *Scte35Parser) Parse(pkt []byte) (cue *Cue) {
	pid := parsePid(pkt[1], pkt[2])
	pl := scte35p.parsePayload(pkt)
	pay := scte35p.chkPartial(pl, pid, []byte("\xfc0"))
	if len(pay) != 0 {
		seclen := parseLen(pay[1], pay[2])
		if scte35p.sectionDone(pay, pid, seclen) {
			cue = &Cue{}
			if cue.Decode(pay) {
				cue.Show()
				scte35p.Cues = append(scte35p.Cues, cue)
				return cue
			}
		}
	}
	return
}

// initialize and return a *Scte35parser
func NewScte35Parser() *Scte35Parser {
	sp := &Scte35Parser{}
	sp.Pids = &Pids{}
	sp.mkMaps()
	return sp
}
