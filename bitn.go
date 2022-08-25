package cuei

import (
	//"encoding/hex"
	"fmt"
	"math/big"
)

// Bitn converts bytes to a list of bits.
type Bitn struct {
	idx  uint
	bits string
}

// Load raw bytes and convert to bits
func (b *Bitn) Load(bites []byte) {
	i := new(big.Int)
	i.SetBytes(bites)
	b.bits = fmt.Sprintf("%b", i)
	b.idx = 0
}

// Chunk slices bitcount of bits and returns it as a uint64
func (b *Bitn) Chunk(bitcount uint) *big.Int {
	j := new(big.Int)
	d := b.idx + bitcount
	j.SetString(b.bits[b.idx:d], 2)
	b.idx = d
	return j
}

// AsUInt8 trims AsUInt64 to 8 bits for smaller numbers
func (b *Bitn) AsUInt8(bitcount uint) uint8 {
	j := b.AsUInt64(bitcount)
	return uint8(j)

}

// AsUInt64 is a wrapper for Chunk
func (b *Bitn) AsUInt64(bitcount uint) uint64 {
	j := b.Chunk(bitcount)
	return j.Uint64()

}

// AsBool slices 1 bit and returns true for 1 , false for 0
func (b *Bitn) AsBool() bool {
	var bitcount uint
	bitcount = 1
	j := b.AsUInt64(bitcount)
	return j == 1
}

// AsFloat slices bitcount of bits and returns as float64
func (b *Bitn) AsFloat(bitcount uint) float64 {
	j := b.AsUInt64(bitcount)
	return float64(j)
}

// As90k is AsFloat / 90000.00 rounded to six decimal places.
func (b *Bitn) As90k(bitcount uint) float64 {
	as90k := b.AsFloat(bitcount) / 90000.00
	return float64(uint64(as90k*1000000)) / 1000000
}

// AsHex slices bitcount of bits and returns as hex string
func (b *Bitn) AsHex(bitcount uint) string {
	j := b.AsUInt64(bitcount)
	ashex := fmt.Sprintf("%#x", j)
	return ashex
}

// AsBytes slices bitcount of bits and returns as []bytes
func (b *Bitn) AsBytes(bitcount uint) []byte {
	j := b.Chunk(bitcount)
	return j.Bytes()
}

// AsAscii returns the Ascii chars of AsBytes
func (b *Bitn) AsAscii(bitcount uint) string {
	return string(b.AsBytes(bitcount))
}

// Forward advances b.idx by bitcount
func (b *Bitn) Forward(bitcount uint) {
	b.idx += bitcount
}
