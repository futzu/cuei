package cuei_test

import (
	"fmt"
	"github.com/futzu/cuei"
	"testing"
)

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
	cue.Decode(data)
	cue.Show()
}

func ExampleCue_Decode() {
	data := "/DAWAAAAAAAAAP/wBQb+AKmKxwAACzuu2Q=="
	cue := cuei.NewCue()
	cue.Decode(data)

	// Cue values can be accessed via dot notiation

	fmt.Println("Cue values can be accessed via dot notiation,")
	fmt.Println("Cue.Command.Show()")
	cue.Command.Show()
	fmt.Println("\n\nSet cue.Command.PTS")
	cue.Command.PTS = 987.654321
	fmt.Printf("cue.Command.PTS = %v\n\n", cue.Command.PTS)
	fmt.Println("Cue.Command.Show()")
	cue.Command.Show()

}

func ExampleCue_Encode() {
	data := "/DAWAAAAAAAAAP/wBQb+AKmKxwAACzuu2Q=="
	cue := cuei.NewCue()
	cue.Decode(data)
	fmt.Println(cue.Encode())
}

func ExampleCue_Encode2B64() {
	data := "/DAWAAAAAAAAAP/wBQb+AKmKxwAACzuu2Q=="
	cue := cuei.NewCue()
	cue.Decode(data)
	fmt.Println(cue.Encode2B64())
}

func ExampleCue_Encode2Hex() {
	data := "/DAWAAAAAAAAAP/wBQb+AKmKxwAACzuu2Q=="
	cue := cuei.NewCue()
	cue.Decode(data)
	fmt.Println(cue.Encode2Hex())
	cue.Decode(cue.Encode2Hex())
	cue.Show()
	cue.Command.Show()
}

func ExampleCue_AdjustPts() {
	data := "/DAWAAAAAAAAAP/wBQb+AKmKxwAACzuu2Q=="
	cue := cuei.NewCue()
	cue.Decode(data)
	cue.Show()
	// Change cue.InfoSection.PtsAdjustment and re-encode cue to bytes
	cue.AdjustPts(33.333)
	cue.Show()
	fmt.Println("Was", data)
	fmt.Println("Is", cue.Encode2B64())
}

func Test(t *testing.T) {

	t.Run("Json2Cue", func(t *testing.T) {
		ExampleJson2Cue()
	})
	t.Run("Cue.Decode()", func(t *testing.T) {
		ExampleCue_Decode()
	})
	t.Run("Cue.AdjustPts()", func(t *testing.T) {
		ExampleCue_AdjustPts()
	})
	t.Run("NewCue", func(t *testing.T) {
		ExampleNewCue()
	})
	t.Run("Cue_Encode", func(t *testing.T) {
		ExampleCue_Encode()
	})
	t.Run("Cue_Encode2B64", func(t *testing.T) {
		ExampleCue_Encode2B64()
	})
	t.Run("Cue_Encode2Hex", func(t *testing.T) {
		ExampleCue_Encode2Hex()
	})

}
