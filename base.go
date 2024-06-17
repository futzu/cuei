package cuei

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/big"
)

// chk generic catchall error checking
func chk(e error) {
	if e != nil {
		fmt.Println(e)

	}
}

// decB64 decodes base64 strings.
func decB64(b64 string) []byte {
	deb64, err := base64.StdEncoding.DecodeString(b64)
	chk(err)
	return deb64
}

// encB64 encodes  bytes to a Base64 string
func encB64(data []byte) string {
	b64 := base64.StdEncoding.EncodeToString(data)
	return b64
}

// hex2Int Hexadecimal string to uint64
func hex2Int(str string) uint64 {
	i := new(big.Int)
	_, err := fmt.Sscan(str, i)
	chk(err)
	return i.Uint64()
}

// Hexed converts bytes to hex string
func Hexed(somebites []byte) string {
	b := new(big.Int)
	b.SetBytes(somebites)
	return fmt.Sprintf("0x%v", b.Text(16))
}

// IsIn is a test for slice membership
func IsIn(slice []uint16, val uint16) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}

func mk90k(raw uint64) float64 {
	nk := float64(raw) / 90000.0
	return float64(uint64(nk*1000000)) / 1000000
}

// Mk90k converts ticks to seconds
func Mk90k(raw uint64) float64 {
	return mk90k(raw)
}

// mkJson structs to JSON
func mkJson(i interface{}) string {
	jason, err := json.MarshalIndent(&i, "", "    ")
	chk(err)
	return string(jason)
}

// Take a JSON string and return a *Cue
func Json2Cue(s string) *Cue {
	b := []byte(s)
	cue := NewCue()
	json.Unmarshal(b, cue)
	cue.Encode()
	return cue
}

func parseLen(byte1, byte2 byte) uint16 {
	return uint16(byte1&0xf)<<8 | uint16(byte2)
}

func parsePid(byte1, byte2 byte) uint16 {
	return uint16(byte1&0x1f)<<8 | uint16(byte2)
}

func parsePrgm(byte1, byte2 byte) uint16 {
	return uint16(byte1)<<8 | uint16(byte2)
}

func splitByIdx(payload, sep []byte) []byte {
	idx := bytes.Index(payload, sep)
	if idx == -1 {
		return []byte("")
	}
	return payload[idx:]
}
