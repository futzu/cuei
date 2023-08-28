package cuei

const zero = 0x00
const one = 0x01
const eight = 0x08
const twentyFour = 0x18
const twoFiftyFive = 0xFF
const twoFiftySix = 0x100
const mask = 0x80000000      // Crc32 mask
const initValue = 0xFFFFFFFF // initial Crc32 value
const initPoly = 0x104C11DB7 // Polynomial value for cRC32 table

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
func mkTable() [twoFiftySix]int {
	var tbl [twoFiftySix]int
	newPoly := initPoly & initValue
	for idx := range tbl {
		tbl[idx] = bytecrc((idx << twentyFour), newPoly)
	}
	return tbl
}

// MkCrc32 generate a 32 bit Crc
func MkCrc32(data []byte) uint32 {
	crc := initValue
	tbl := mkTable()
	for _, bite := range data {
		crc = tbl[int(bite)^((crc>>twentyFour)&twoFiftyFive)] ^ ((crc << eight) & (initValue - twoFiftyFive))
	}
	return uint32(crc)
}
