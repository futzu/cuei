package cuei

import (
	"fmt"
	"math"
	"math/big"
)

// Decoder converts bytes to a list of bits.
type bitDecoder struct {
	idx  uint
	bits string
	last uint
}

// Load raw bytes and convert to bits
func (bd *bitDecoder) load(bites []byte) {
	i := new(big.Int)
	i.SetBytes(bites)
	bd.bits = fmt.Sprintf("%b", i)
	bd.last = uint(len(bd.bits))
	bd.idx = 0
}

// chunk slices bitcount of bits and returns it as a uint64
func (bd *bitDecoder) chunk(bitcount uint) *big.Int {
	j := new(big.Int)
	if (bd.idx + bitcount) <= bd.last-32 {
		d := bd.idx + bitcount
		j.SetString(bd.bits[bd.idx:d], 2)
		bd.idx = d
	}
	return j
}

// crc
func (bd *bitDecoder) crc() string {
	j := new(big.Int)
	j.SetString(bd.bits[bd.last-32:], 2)
	ashex := fmt.Sprintf("%#x", j)
	return ashex
}

// uInt8 trims uint64 to 8 bits
func (bd *bitDecoder) uInt8(bitcount uint) uint8 {
	j := bd.uInt64(bitcount)
	return uint8(j)

}

// uInt16 trims uint64 to 16 bits
func (bd *bitDecoder) uInt16(bitcount uint) uint16 {
	j := bd.uInt64(bitcount)
	return uint16(j)

}

// uInt32 trims uint64 to 32 bits
func (bd *bitDecoder) uInt32(bitcount uint) uint32 {
	j := bd.uInt64(bitcount)
	return uint32(j)

}

// uInt64 is a wrapper for chunk
func (bd *bitDecoder) uInt64(bitcount uint) uint64 {
	j := bd.chunk(bitcount)
	return j.Uint64()

}

// asFlag slices 1 bit and returns true for 1 , false for 0
func (bd *bitDecoder) asFlag() bool {
	var bitcount uint
	bitcount = 1
	j := bd.uInt64(bitcount)
	return j == 1
}

// asFloat slices bitcount of bits and returns  float64
func (bd *bitDecoder) asFloat(bitcount uint) float64 {
	j := bd.uInt64(bitcount)
	return float64(j)
}

// as90k is Float / 90000.00 rounded to six decimal places.
func (bd *bitDecoder) as90k(bitcount uint) float64 {
	as90k := bd.asFloat(bitcount) / 90000.00
	return float64(uint64(as90k*1000000)) / 1000000
}

// asHex slices bitcount of bits and returns as hex string
func (bd *bitDecoder) asHex(bitcount uint) string {
	j := bd.uInt64(bitcount)
	ashex := fmt.Sprintf("%#x", j)
	return ashex
}

// asBytes slices bitcount of bits and returns as []bytes
func (bd *bitDecoder) asBytes(bitcount uint) []byte {
	j := bd.chunk(bitcount)
	return j.Bytes()
}

// asAscii returns the ascii chars of Bytes
func (bd *bitDecoder) asAscii(bitcount uint) string {
	return string(bd.asBytes(bitcount))
}

// goForward advances g.idx by bitcount
func (bd *bitDecoder) goForward(bitcount uint) {
	bd.idx += bitcount
}

// Encoder packs  data as bits for encoding.
type bitEncoder struct {
	Bites big.Int
}

// Append a []byte as bits
func (be *bitEncoder) AddBytes(bites []byte, nbits uint) {
	t := new(big.Int)
	t.SetBytes(bites)
	o := be.Bites.Lsh(&be.Bites, nbits)
	be.Bites = *be.Bites.Add(o, t)
}

/*
Add left shifts Encoder.Bites by nbits and add val interface{} as bits.
Supports val as bool, float64, int, uint8, uint16, uint32,or  uint64.
*/
func (be *bitEncoder) Add(val interface{}, nbits uint) {
	t := new(big.Int)
	t.SetUint64(u64(val))
	o := be.Bites.Lsh(&be.Bites, nbits)
	be.Bites = *be.Bites.Add(o, t)
}

// AddHex64 append a hex string as uint64 in bits
func (be *bitEncoder) AddHex64(val string, nbits uint) {
	u := new(big.Int)
	_, err := fmt.Sscan(val, u)
	if err != nil {
		fmt.Println("error scanning value:", err)
	} else {
		be.Add(u.Uint64(), nbits)
	}
}

// AddHex32 append a hex string as uint32 in bits
func (be *bitEncoder) AddHex32(val string, nbits uint) {
	u := new(big.Int)
	_, err := fmt.Sscan(val, u)
	if err != nil {
		fmt.Println("error scanning value:", err)
	} else {
		be.Add(uint32(u.Uint64()), nbits)
	}
}

// Reserve left shifts Encoder.Bites by num and adds num bits  set to 1
func (be *bitEncoder) Reserve(num int) {
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
		return uint64(math.Round(i.(float64) * float64(90000.0)))
	case bool:
		if i == true {
			return uint64(1)
		}
		return uint64(0)
	default:
		return uint64(i.(uint64))

	}
}
