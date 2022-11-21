package cuei

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
)

// chk generic catchall error checking
func chk(e error) {
	if e != nil {
		fmt.Println(e)

	}
}

// DeB64 decodes base64 strings.
func DeB64(b64 string) []byte {
	deb64, err := base64.StdEncoding.DecodeString(b64)
	chk(err)
	return deb64
}

// IsIn is a test for slice membership
func IsIn[T comparable](slice []T, val T) bool {
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

// MkJson structs to JSON
func MkJson(i interface{}) string {
	jason, err := json.MarshalIndent(&i, "", "    ")
	chk(err)
	return string(jason)
}

func parseLen(byte1 byte, byte2 byte) uint16 {
	return uint16(byte1&0xf)<<8 | uint16(byte2)
}

func parsePid(byte1 byte, byte2 byte) uint16 {
	return uint16(byte1&0x1f)<<8 | uint16(byte2)
}

func parsePrgm(byte1 byte, byte2 byte) uint16 {
	return uint16(byte1)<<8 | uint16(byte2)
}

func splitByIdx(payload, sep []byte) []byte {
	idx := bytes.Index(payload, sep)
	if idx == -1 {
		return []byte("")
	}
	return payload[idx:]
}
