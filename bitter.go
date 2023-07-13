package cuei

import (
	"fmt"
	"math/big"
)

// Decoder converts bytes to a list of bits.
type BitDecoder struct {
	idx  uint
	bits string
}

// Load raw bytes and convert to bits
func (bd *BitDecoder) Load(bites []byte) {
	i := new(big.Int)
	i.SetBytes(bites)
	bd.bits = fmt.Sprintf("%b", i)
	bd.idx = 0
}

// Chunk slices bitcount of bits and returns it as a uint64
func (bd *BitDecoder) Chunk(bitcount uint) *big.Int {
	j := new(big.Int)
	d := bd.idx + bitcount
	j.SetString(bd.bits[bd.idx:d], 2)
	bd.idx = d
	return j
}

// uint8 trims uint64 to 8 bits
func (bd *BitDecoder) UInt8(bitcount uint) uint8 {
	j := bd.UInt64(bitcount)
	return uint8(j)

}

// uint16 trims uint64 to 16 bits
func (bd *BitDecoder) UInt16(bitcount uint) uint16 {
	j := bd.UInt64(bitcount)
	return uint16(j)

}

// uint32 trims uint64 to 32 bits
func (bd *BitDecoder) UInt32(bitcount uint) uint32 {
	j := bd.UInt64(bitcount)
	return uint32(j)

}

// uint64 is a wrapper for Chunk
func (bd *BitDecoder) UInt64(bitcount uint) uint64 {
	j := bd.Chunk(bitcount)
	return j.Uint64()

}

// Flag slices 1 bit and returns true for 1 , false for 0
func (bd *BitDecoder) Flag() bool {
	var bitcount uint
	bitcount = 1
	j := bd.UInt64(bitcount)
	return j == 1
}

func (bd *BitDecoder) Bool() bool {
	return bd.Flag()
}

// Float slices bitcount of bits and returns  float64
func (bd *BitDecoder) Float(bitcount uint) float64 {
	j := bd.UInt64(bitcount)
	return float64(j)
}

// As90k is Float / 90000.00 rounded to six decimal places.
func (bd *BitDecoder) As90k(bitcount uint) float64 {
	as90k := bd.Float(bitcount) / 90000.00
	return float64(uint64(as90k*1000000)) / 1000000
}

// Hex slices bitcount of bits and returns as hex string
func (bd *BitDecoder) Hex(bitcount uint) string {
	j := bd.UInt64(bitcount)
	ashex := fmt.Sprintf("%#x", j)
	return ashex
}

// Bytes slices bitcount of bits and returns as []bytes
func (bd *BitDecoder) Bytes(bitcount uint) []byte {
	j := bd.Chunk(bitcount)
	return j.Bytes()
}

// Ascii returns the Ascii chars of Bytes
func (bd *BitDecoder) Ascii(bitcount uint) string {
	return string(bd.Bytes(bitcount))
}

// Forward advances g.idx by bitcount
func (bd *BitDecoder) Forward(bitcount uint) {
	bd.idx += bitcount
}

// Encoder packs  data as bits for encoding.
type BitEncoder struct {
	Bites big.Int
}

// Append a []byte as bits
func (be *BitEncoder) AddBytes(bites []byte, nbits uint) {
	t := new(big.Int)
	t.SetBytes(bites)
	o := be.Bites.Lsh(&be.Bites, nbits)
	be.Bites = *be.Bites.Add(o, t)
}

/*
Add left shifts Encoder.Bites by nbits and add val interface{} as bits.
Supports val as bool, float64, int, uint8, uint16, uint32,or  uint64.
*/
func (be *BitEncoder) Add(val interface{}, nbits uint) {
	t := new(big.Int)
	t.SetUint64(u64(val))
	o := be.Bites.Lsh(&be.Bites, nbits)
	be.Bites = *be.Bites.Add(o, t)
}

// AddHex64 append a hex string as uint64 in bits
func (be *BitEncoder) AddHex64(val string, nbits uint) {
	u := new(big.Int)
	_, err := fmt.Sscan(val, u)
	if err != nil {
		fmt.Println("error scanning value:", err)
	} else {
		be.Add(u.Uint64(), nbits)
	}
}

// Reserve left shifts Encoder.Bites by num and adds num bits  set to 1
func (be *BitEncoder) Reserve(num int) {
	for i := 0; i < num; i++ {
		be.Add(1, 1)
	}
}

/*
	 u64 takes a bool, float64, int ,uint8, ,uint16 uint32, or uint64
		and returns a uint64
*/
func u64(i interface{}) uint64 {
	switch i.(type) {
	case int:
		return uint64(i.(int))
	case uint8:
		return uint64(i.(uint8))
	case uint16:
		return uint64(i.(uint16))
	case uint32:
		return uint64(i.(uint32))
	case float64:
		return uint64(i.(float64) * float64(90000.0))
	case bool:
		if i == true {
			return uint64(1)
		}
		return uint64(0)
	default:
		return uint64(i.(uint64))

	}
}
