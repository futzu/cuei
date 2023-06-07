package cuei

import (
	"fmt"
	"math/big"
)

type Nbin struct{
    Bites big.Int
    
}

func (nb *Nbin) AddBytes(str string,nbits uint) {
        t := new(big.Int)
        t.SetBytes([]byte(str))
        o := nb.Bites.Lsh(&nb.Bites, nbits)
        nb.Bites = *nb.Bites.Add(o,t)
}

func (nb *Nbin) Add64(val uint64,nbits uint) {
        t := new(big.Int)
        t.SetUint64(val)
        o := nb.Bites.Lsh(&nb.Bites, nbits)
        nb.Bites = *nb.Bites.Add(o,t)
}

func (nb *Nbin) Add32(val uint32,nbits uint) {
        u := uint64(val)
        nb.Add64(u,nbits)
}

func (nb *Nbin) Add16(val uint16,nbits uint) {
       u := uint64(val)
        nb.Add64(u,nbits)
}

func (nb *Nbin) Add8(val uint8,nbits uint) {
      u := uint64(val)
      nb.Add64(u,nbits)
}

func (nb *Nbin) AddBool(val bool) {
    if val == true {
        nb.Add64(1,1)
    } else {
        nb.Add64(0,1)
    }
}

