package cuei

// Polynomial value for CRC32 table
const initPoly = 0x104C11DB7

// initial Crc32 value
const initValue= 0xFFFFFFFF

// Crc32 mask
const mask = 0x80000000
const zero = 0x0
const one = 0x1
const eight = 0x8
const twentyFour = 0x18
const twoFiftyFive= 0xFF
const twoFiftySix = 0x100

// bytecrc creates the values used to populate the table
func bytecrc(crc int, aPoly int) int {
	for i := 0; i < eight; i++ {
		if crc& mask != zero {
			crc = crc<<one ^ aPoly
		} else {
			crc = crc << one
		}
	}
	return int(crc &initValue)
}

// mkTable makes the Crc32 table
func mkTable() [256]int {
	var tbl [twoFiftySix]int
	newPoly := initPoly & initValue
	for idx, _ := range tbl {
		tbl[idx] = (bytecrc((idx << twentyFour), newPoly))
	}
	return tbl
}

// generate a 32 bit Crc
func CRC32(data []byte) uint32 {
	crc := initValue
	tbl := mkTable()
	for _, bite := range data {
		crc = tbl[int(bite)^((crc>>twentyFour)& twoFiftyFive)] ^ ((crc << eight) & (initValue- twoFiftyFive))
	}
	return uint32(crc)
}
