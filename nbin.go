package cuei

import (
	"fmt"
	"math/big"
)

// Nbin packs  SCTE-35 data as bits for encoding and stores them in a big.Int called Bites
type Nbin struct {
	Bites big.Int
}

// Append a string  as bits to NBin.Bites for encoding
func (nb *Nbin) AddBytes(bites []byte, nbits uint) {
	t := new(big.Int)
	t.SetBytes(bites)
	o := nb.Bites.Lsh(&nb.Bites, nbits)
	nb.Bites = *nb.Bites.Add(o, t)
}

// Add64  append  a uint64 as bits to NBin.Bites
func (nb *Nbin) Add64(val uint64, nbits uint) {
	t := new(big.Int)
	t.SetUint64(val)
	o := nb.Bites.Lsh(&nb.Bites, nbits)
	nb.Bites = *nb.Bites.Add(o, t)
}

// Add32  append a  unit32 as bits to NBin.Bites
func (nb *Nbin) Add32(val uint32, nbits uint) {
	u := uint64(val)
	nb.Add64(u, nbits)
}

// Add16  append a  unit16 as bits to NBin.Bites
func (nb *Nbin) Add16(val uint16, nbits uint) {
	u := uint64(val)
	nb.Add64(u, nbits)
}

// Add8  append a  unit8 as bits to NBin.Bites
func (nb *Nbin) Add8(val uint8, nbits uint) {
	u := uint64(val)
	nb.Add64(u, nbits)
}

// AddFlag append a bool as a bit to NBin.Bytes
func (nb *Nbin) AddFlag(val bool) {
	if val == true {
		nb.Add64(1, 1)
	} else {
		nb.Add64(0, 1)
	}
}

// Add90k append a 90k clock value, as ticks, in bits, to NBin.Bytes
func (nb *Nbin) Add90k(val float64, nbits uint) {
	u := uint64(val * float64(90000.0))
	nb.Add64(u, nbits)
}

// AddHex64 append a hex string as uint64 in bits to NBin.Bites
func (nb *Nbin) AddHex64(val string, nbits uint) {
	u := new(big.Int)
	_, err := fmt.Sscan(val, u)
	if err != nil {
		fmt.Println("error scanning value:", err)
	} else {
		fmt.Println(u.Uint64())
		nb.Add64(u.Uint64(), nbits)
	}
}

// Reserve num bits by setting them to 1
func (nb *Nbin) Reserve(num int) {

	for i := 0; i < num; i++ {
		nb.Add64(1, 1)
	}
}
