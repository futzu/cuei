package cuei_test

import (
	"fmt"
	"github.com/futzu/cuei"
	"testing"
)

func ExampleDecB64() {
	b64 := "/DARAAAAAAAAAP/wAAAAAHpPv/8="
	dcoded := cuei.DecB64(b64)
	fmt.Println(dcoded)
}

func ExampleEncB64() {
	somebytes := []byte{252, 48, 17, 0, 0, 0, 0, 0, 0, 0, 255, 240, 0, 0, 0, 0, 122, 79, 191, 255}
	ncoded := cuei.EncB64(somebytes)
	fmt.Println(ncoded)
}

func ExampleJson2Cue() {

	js := `{
    "InfoSection": {
        "Name": "Splice Info Section",
        "TableID": "0xfc",
        "SectionSyntaxIndicator": false,
        "Private": false,
        "Reserved": "0x3",
        "SectionLength": 42,
        "ProtocolVersion": 0,
        "EncryptedPacket": false,
        "EncryptionAlgorithm": 0,
        "PtsAdjustment": 0,
        "CwIndex": "0xff",
        "Tier": "0xfff",
        "CommandLength": 15,
        "CommandType": 5
    },
    "Command": {
        "Name": "Splice Insert",
        "CommandType": 5,
        "SpliceEventID": 5690,
        "OutOfNetworkIndicator": true,
        "ProgramSpliceFlag": true,
        "TimeSpecifiedFlag": true,
        "PTS": 23683.480033
    },
    "DescriptorLoopLength": 10,
    "Descriptors": [
        {
            "Length": 8,
            "Identifier": "CUEI",
            "Name": "Avail Descriptor"
        }
    ],
    "Crc32": 3608566905
}
`
	cue := cuei.Json2Cue(js)
	cue.Encode()
	cue.Show()
}

func ExampleNewCue() {

	data := "/DCtAAAAAAAAAP/wBQb+Tq9DwQCXAixDVUVJCUvhcH+fAR1QQ1IxXzEyMTYyMTE0MDBXQUJDUkFDSEFFTFJBWSEBAQIsQ1VFSQlL4W9/nwEdUENSMV8xMjE2MjExNDAwV0FCQ1JBQ0hBRUxSQVkRAQECGUNVRUkJTBwVf58BClRLUlIxNjA4NEEQAQECHkNVRUkJTBwWf98AA3clYAEKVEtSUjE2MDg0QSABAdHBXYA="

	cue := cuei.NewCue()
	cue.Decode(cuei.DecB64(data))
	cue.Show()
}

func Test(t *testing.T) {

	t.Run("DecB64", func(t *testing.T) {
		ExampleDecB64()
	})
	t.Run("EncB64", func(t *testing.T) {
		ExampleEncB64()
	})
	t.Run("Json2Cue", func(t *testing.T) {
		ExampleJson2Cue()
	})
	t.Run("NewCue", func(t *testing.T) {
		ExampleNewCue()
	})

}
