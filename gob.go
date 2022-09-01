// Group of bits

package cuei

import (
	"fmt"
	"math/big"
)

// Gob converts bytes to a list of bits.
type Gob struct {
	idx  uint
	bits string
}

// Load raw bytes and convert to bits
func (g *Gob) Load(bites []byte) {
	i := new(big.Int)
	i.SetBytes(bites)
	g.bits = fmt.Sprintf("%b", i)
	g.idx = 0
}

// Chunk slices bitcount of bits and returns it as a uint64
func (g *Gob) Chunk(bitcount uint) *big.Int {
	j := new(big.Int)
	d := g.idx + bitcount
	j.SetString(g.bits[g.idx:d], 2)
	g.idx = d
	return j
}

// uint8 trims uint64 to 8 bits
func (g *Gob) UInt8(bitcount uint) uint8 {
	j := g.UInt64(bitcount)
	return uint8(j)

}

// uint16 trims uint64 to 16 bits
func (g *Gob) UInt16(bitcount uint) uint16 {
	j := g.UInt64(bitcount)
	return uint16(j)

}

// uint32 trims uint64 to 32 bits
func (g *Gob) UInt32(bitcount uint) uint32 {
	j := g.UInt64(bitcount)
	return uint32(j)

}

// uint64 is a wrapper for Chunk
func (g *Gob) UInt64(bitcount uint) uint64 {
	j := g.Chunk(bitcount)
	return j.Uint64()

}

// Flag slices 1 bit and returns true for 1 , false for 0
func (g *Gob) Flag() bool {
	var bitcount uint
	bitcount = 1
	j := g.UInt64(bitcount)
	return j == 1
}

func (g *Gob) Bool() bool {
	return g.Flag()
}

// Float slices bitcount of bits and returns  float64
func (g *Gob) Float(bitcount uint) float64 {
	j := g.UInt64(bitcount)
	return float64(j)
}

// As90k is Float / 90000.00 rounded to six decimal places.
func (g *Gob) As90k(bitcount uint) float64 {
	as90k := g.Float(bitcount) / 90000.00
	return float64(uint64(as90k*1000000)) / 1000000
}

// Hex slices bitcount of bits and returns as hex string
func (g *Gob) Hex(bitcount uint) string {
	j := g.UInt64(bitcount)
	ashex := fmt.Sprintf("%#x", j)
	return ashex
}

// Bytes slices bitcount of bits and returns as []bytes
func (g *Gob) Bytes(bitcount uint) []byte {
	j := g.Chunk(bitcount)
	return j.Bytes()
}

// Ascii returns the Ascii chars of Bytes
func (g *Gob) Ascii(bitcount uint) string {
	return string(g.Bytes(bitcount))
}

// Forward advances g.idx by bitcount
func (g *Gob) Forward(bitcount uint) {
	g.idx += bitcount
}
