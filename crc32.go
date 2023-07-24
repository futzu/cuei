package cuei

// Polynomial value for cRC32 table
const initPoly = 0x104C11DB7
// initial Crc32 value
const initValue = 0xFFFFFFFF
// Crc32 mask
const mask = 0x80000000
// zero 
const zero = 0x0
// one
const one = 0x1
// eight
const eight = 0x8
// twenty-four
const twentyFour = 0x18
// two hundred and fifty-five
const twoFiftyFive = 0xFF
// two hundred and fifty-six
const twoFiftySix = 0x100

// bytecrc creates the values used to populate the table
func bytecrc(crc int, aPoly int) int {
	for i := 0; i < eight; i++ {
		if crc&mask != zero {
			crc = crc<<one ^ aPoly
		} else {
			crc = crc << one
		}
	}
	return int(crc & initValue)
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
func cRC32(data []byte) uint32 {
	crc := initValue
	tbl := mkTable()
	for _, bite := range data {
		crc = tbl[int(bite)^((crc>>twentyFour)&twoFiftyFive)] ^ ((crc << eight) & (initValue - twoFiftyFive))
	}
	return uint32(crc)
}
