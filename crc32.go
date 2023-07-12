package cuei

// Polynomial value for CRC32 table
const POLY = 0x104C11DB7
// initial Crc32 value
const INIT_VALUE = 0xFFFFFFFF
// Crc32 mask
const MASK  = 0x80000000
const ZERO = 0x0
const ONE = 0x1
const EIGHT = 0x8
const TWENTY_FOUR = 0x18
const TWO_FIFTY_FIVE = 0xFF
const TWO_FIFTY_SIX = 0x100


// bytecrc creates the values used to populate the table
func bytecrc(crc int, poly int) int {
	for i := 0; i < EIGHT; i++ {
		if crc& MASK != ZERO {
			crc = crc<<ONE ^ poly
		} else {
			crc = crc << ONE
		}
	}
	return int(crc & INIT_VALUE)
}

// mkTable makes the Crc32 table
func mkTable() [256]int {
	var tbl [TWO_FIFTY_SIX]int
	poly := POLY & INIT_VALUE
	for idx, _ := range tbl {
		tbl[idx] = (bytecrc((idx << TWENTY_FOUR), poly))
	}
	return tbl
}

// generate a 32 bit Crc
func CRC32(data []byte) uint32 {
	crc := INIT_VALUE
	tbl := mkTable()
	for _, bite := range data {
		crc = tbl[int(bite)^((crc>>TWENTY_FOUR)&TWO_FIFTY_FIVE)] ^ ((crc << EIGHT) & (INIT_VALUE - TWO_FIFTY_FIVE))
	}
	return uint32(crc)
}
