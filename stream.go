package cuei

import (
	"bytes"
	"os"
)

type PacketData struct {
	PacketNumber int     `json:",omitempty"`
	Pid          uint16  `json:",omitempty"`
	Program      uint16  `json:",omitempty"`
	Pcr          float64 `json:",omitempty"`
	Pts          float64 `json:",omitempty"`
}

// pktSz is the size of an MPEG-TS packet in bytes.
const pktSz = 188

// bufSz is the size of a read when parsing files.
const bufSz = 16384 * pktSz

// Stream for parsing MPEGTS for SCTE-35
type Stream struct {
	Cues []*Cue
	Pids
	Pid2Prgm map[uint16]uint16 // pid to program map
	Pid2Type map[uint16]uint8  // pid to stream type map
	Programs []uint16
	Prgm2Pcr map[uint16]uint64 // program to pcr map
	Prgm2Pts map[uint16]uint64 // program to pts map
	last     map[uint16][]byte // last compares current packet payload to last packet payload by pid
	partial  map[uint16][]byte // partial manages tables spread across multiple packets by pid
	pktNum   int               // packet count.

}

func (stream *Stream) mkMaps() {
	stream.Pid2Prgm = make(map[uint16]uint16)
	stream.Pid2Type = make(map[uint16]uint8)
	stream.Prgm2Pcr = make(map[uint16]uint64)
	stream.Prgm2Pts = make(map[uint16]uint64)
	stream.last = make(map[uint16][]byte)
	stream.partial = make(map[uint16][]byte)
}

// Decode fname (a file name) for SCTE-35
func (stream *Stream) Decode(fname string) []*Cue {
	stream.mkMaps()
	stream.pktNum = 0
	file, err := os.Open(fname)
	chk(err)
	defer file.Close()
	buffer := make([]byte, bufSz)
	for {
		bytesread, err := file.Read(buffer)
		if err != nil {
			break
		}
		for i := 1; i <= (bytesread / pktSz); i++ {
			end := i * pktSz
			start := end - pktSz
			p := buffer[start:end]
			pkt := &p
			stream.pktNum++
			stream.parse(*pkt)
		}
	}
	return stream.Cues
}

// MkPcr takes a program and returns the current PCR
func (stream *Stream) MkPcr(prgm uint16) float64 {
	pcrb := stream.Prgm2Pcr[prgm]
	return mk90k(pcrb)
}

// MkPts takes a program and returns the current PTS
func (stream *Stream) MkPts(prgm uint16) float64 {
	pts := stream.Prgm2Pts[prgm]
	return mk90k(pts)
}

// parsePusi returns true if PUSI flag is set
func (stream *Stream) parsePusi(pkt []byte) bool {
	return (pkt[1]&0x40 == 0x40)

}

// afcFlag returns true if AFC flag is set
func (stream *Stream) afcFlag(pkt []byte) bool {
	return (pkt[3]&0x20 == 0x20)
}

// pcrFlag returns true if PCR flag is set
func (stream *Stream) pcrFlag(pkt []byte) bool {
	return (pkt[5]&0x10 == 0x10)
}

// ptsFlag returns true if PTS flag is set
func (stream *Stream) ptsFlag(pay []byte) bool {
	return (pay[7]&0x80 == 0x80)
}

// parsePts parses a packet for PTS
func (stream *Stream) parsePts(pay []byte, pid uint16) {
	if stream.ptsFlag(pay) {
		prgm, ok := stream.Pid2Prgm[pid]
		if ok {
			pts := uint64(pay[9]&14) << 29
			pts |= uint64(pay[10]) << 22
			pts |= (uint64(pay[11]) >> 1) << 15
			pts |= uint64(pay[12]) << 7
			pts |= uint64(pay[13]) >> 1
			stream.Prgm2Pts[prgm] = pts
		}
	}
}

// parsePcr parses a packet for PCR
func (stream *Stream) parsePcr(pkt []byte, pid uint16) {
	if stream.afcFlag(pkt) {
		if stream.pcrFlag(pkt) {
			pcr := (uint64(pkt[6]) << 25)
			pcr |= (uint64(pkt[7]) << 17)
			pcr |= (uint64(pkt[8]) << 9)
			pcr |= (uint64(pkt[9]) << 1)
			pcr |= uint64(pkt[10]) >> 7
			prgm := stream.Pid2Prgm[pid]
			stream.Prgm2Pcr[prgm] = pcr
		}
	}
}

// parsePay packet payload starts after header and afc (if present)
func (stream *Stream) parsePayload(pkt []byte) []byte {
	head := 4
	if stream.afcFlag(pkt) {
		afl := int(pkt[4])
		head += afl + 1
	}
	if head > pktSz {
		head = pktSz
	}
	return pkt[head:]
}

// chkPartial appends the current packet payload to partial table by pid.
func (stream *Stream) chkPartial(pay []byte, pid uint16, sep []byte) []byte {
	val, ok := stream.partial[pid]
	if ok {
		pay = append(val, pay...)
	}
	return splitByIdx(pay, sep)
}

// sameAsLast compares the current packet to the last packet by pid.
func (stream *Stream) sameAsLast(pay []byte, pid uint16) bool {
	val, ok := stream.last[pid]
	if ok {
		if bytes.Compare(pay, val) == 0 {
			return true
		}
	}
	stream.last[pid] = pay
	return false
}

// sectionDone aggregates partial tables by pid until the section is complete.
func (stream *Stream) sectionDone(pay []byte, pid uint16, seclen uint16) bool {
	if seclen+3 > uint16(len(pay)) {
		stream.partial[pid] = pay
		return false
	}
	delete(stream.partial, pid)
	return true
}

// parse is the parser method for Stream
func (stream *Stream) parse(pkt []byte) {
	p := parsePid(pkt[1], pkt[2])
	pid := &p
	pl := stream.parsePayload(pkt)
	pay := &pl

	if *pid == 0 {
		stream.parsePat(*pay, *pid)
	}
	if stream.isPmtPid(*pid) {
		stream.parsePmt(*pay, *pid)
	}
	if stream.isPcrPid(*pid) {
		stream.parsePcr(pkt, *pid)
	}
	if stream.parsePusi(pkt) {
		stream.parsePts(*pay, *pid)
	}
	if stream.isScte35Pid(*pid) {
		stream.parseScte35(*pay, *pid)
	}
}

// parsePat parses PAT payload
func (stream *Stream) parsePat(pay []byte, pid uint16) {
	if stream.sameAsLast(pay, pid) {
		return
	}
	pay = stream.chkPartial(pay, pid, []byte("\x00\x00"))
	if len(pay) < 1 {
		return
	}
	seclen := parseLen(pay[2], pay[3])
	if stream.sectionDone(pay, pid, seclen) {
		seclen -= 5 // pay bytes 4,5,6,7,8
		idx := uint16(9)
		end := idx + seclen - 4 //  4 bytes for crc
		chunksize := uint16(4)
		for idx < end {
			prgm := parsePrgm(pay[idx], pay[idx+1])
			if prgm > 0 {
				if !IsIn(stream.Programs, prgm) {
					stream.Programs = append(stream.Programs, prgm)
				}
				pmtpid := parsePid(pay[idx+2], pay[idx+3])
				stream.addPmtPid(pmtpid)
			}
			idx += chunksize
		}
	}
}

// parsePmt parses PMT payload
func (stream *Stream) parsePmt(pay []byte, pid uint16) {
	if stream.sameAsLast(pay, pid) {
		return
	}
	pay = stream.chkPartial(pay, pid, []byte("\x02"))
	if len(pay) < 1 {
		return
	}
	secinfolen := parseLen(pay[1], pay[2])
	if stream.sectionDone(pay, pid, secinfolen) {
		prgm := parsePrgm(pay[3], pay[4])
		pcrpid := parsePid(pay[8], pay[9])
		stream.addPcrPid(pcrpid)
		proginfolen := parseLen(pay[10], pay[11])
		idx := uint16(12)
		idx += proginfolen
		silen := secinfolen - 9
		silen -= proginfolen
		stream.parseStreams(silen, pay, idx, prgm)
	}
}

// parseStreams parses program stream information
func (stream *Stream) parseStreams(silen uint16, pay []byte, idx uint16, prgm uint16) {
	chunksize := uint16(5)
	endidx := (idx + silen) - chunksize
	for idx < endidx {
		streamtype := pay[idx]
		elpid := parsePid(pay[idx+1], pay[idx+2])
		eilen := parseLen(pay[idx+3], pay[idx+4])
		idx += chunksize
		idx += eilen
		stream.Pid2Prgm[elpid] = prgm
		stream.Pid2Type[elpid] = streamtype
		stream.vrfyStreamType(elpid, streamtype)
	}
}

// vrfyStreamType checks for stream types 6 and 134 and adds them to Stream.Pids.Scte35Pids
func (stream *Stream) vrfyStreamType(pid uint16, streamtype uint8) {
	if streamtype == 6 || streamtype == 134 {
		stream.addScte35Pid(pid)
	}
}

// parseSCTE35 parses SCTE35 packets
func (stream *Stream) parseScte35(pay []byte, pid uint16) {
	pay = stream.chkPartial(pay, pid, []byte("\xfc0"))
	if len(pay) == 0 {
		stream.Pids.delScte35Pid(pid)
		return
	}
	seclen := parseLen(pay[1], pay[2])
	if stream.sectionDone(pay, pid, seclen) {
		cue := stream.mkCue(pid)
		if cue.Decode(pay) {
			stream.Cues = append(stream.Cues, cue)
			cue.Show()
		} else {
			stream.Pids.delScte35Pid(pid)
		}
	}
}

// mkCue adds PID,PCR, PTS and PacketNumber to a Cue
func (stream *Stream) mkCue(pid uint16) *Cue {
	cue := &Cue{}
	cue.Packet = &PacketData{}
	cue.Packet.Pid = pid
	p := stream.Pid2Prgm[pid]
	prgm := &p
	cue.Packet.Program = *prgm
	cue.Packet.Pcr = stream.MkPcr(*prgm)
	cue.Packet.Pts = stream.MkPts(*prgm)
	cue.Packet.PacketNumber = stream.pktNum
	return cue
}
