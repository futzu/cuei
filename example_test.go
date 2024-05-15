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
    ]

}
`
	cue := cuei.Json2Cue(js)
	cue.Encode()
	cue.Show()
	cue.Command.PTS = 12345.12345
	cue.Encode()
	cue.Show()
	cue2 := cuei.NewCue()
	cue2.Decode(cue.Encode())
}

func ExampleNewCue() {
	data := "/DCtAAAAAAAAAP/wBQb+Tq9DwQCXAixDVUVJCUvhcH+fAR1QQ1IxXzEyMTYyMTE0MDBXQUJDUkFDSEFFTFJBWSEBAQIsQ1VFSQlL4W9/nwEdUENSMV8xMjE2MjExNDAwV0FCQ1JBQ0hBRUxSQVkRAQECGUNVRUkJTBwVf58BClRLUlIxNjA4NEEQAQECHkNVRUkJTBwWf98AA3clYAEKVEtSUjE2MDg0QSABAdHBXYA="
	cue := cuei.NewCue()
	cue.Command = &cuei.Command{}
	cue.Decode(data)
	cue.Show()

}

func ExampleCue_Decode() {
	data := "/DAWAAAAAAAAAP/wBQb+AKmKxwAACzuu2Q=="
	cue := cuei.NewCue()
	cue.Decode(data)
	fmt.Println("Cue.Decode() parses data and populate the fields in the Cue.")
	cue.Show()
	fmt.Println("\n\nCue values can be accessed via dot notiation,")
	cue.Command.PTS = 987.654321
	fmt.Printf("cue.Command.PTS = %v\n", cue.Command.PTS)
	cue.Show()

}

func ExampleCue_Encode() {
	data := "/DAWAAAAAAAAAP/wBQb+AKmKxwAACzuu2Q=="
	cue := cuei.NewCue()
	cue.Decode(data)
	// encode to bytes
	fmt.Println(cue.Encode())
	sp := cuei.SpliceInsert{
		SpliceEventID:              1,
		SpliceEventCancelIndicator: false,
		OutOfNetworkIndicator:      true,
		ProgramSpliceFlag:          true,
		DurationFlag:               true,
		BreakDuration:              60,
		BreakAutoReturn:            true,
		SpliceImmediateFlag:        true,
		EventIDComplianceFlag:      true,
		UniqueProgramID:            1,
		AvailNum:                   0,
		AvailExpected:              0,
	}
	cue.Command = &cuei.Command{}
	cue.Command.CommandType = 5
	cue.Command.Name = "Splice Insert"
	cue.Command.SpliceInsert = sp
	cue.Command.TimeSpecifiedFlag = true
	cue.Command.PTS = 123456.9
	cue.Encode()
	cue.Show()

}

func ExampleCue_Encode2B64() {
	data := "/DAWAAAAAAAAAP/wBQb+AKmKxwAACzuu2Q=="
	cue := cuei.NewCue()
	cue.Decode(data)
	// encode to base64
	fmt.Println(cue.Encode2B64())
}

func ExampleCue_Encode2Hex() {
	data := "/DAWAAAAAAAAAP/wBQb+AKmKxwAACzuu2Q=="
	cue := cuei.NewCue()
	// decode base64 data into cue
	cue.Decode(data)
	// Encode the cue as hex
	hexed := cue.Encode2Hex()
	fmt.Println(hexed)
	// decode the hex back into a Cue
	cue.Decode(hexed)
	cue.Show()
}

func ExampleCue_AdjustPts() {
	data := "/DAWAAAAAAAAAP/wBQb+AKmKxwAACzuu2Q=="
	cue := cuei.NewCue()
	cue.Decode(data)
	fmt.Println("Before calling Cue.AdjustPts")
	fmt.Println(data)
	cue.InfoSection.Show()
	fmt.Println()
	// Change cue.InfoSection.PtsAdjustment and re-encode cue to bytes
	cue.AdjustPts(33.333)
	fmt.Println("After calling Cue.AdjustPts")
	fmt.Println(cue.Encode2B64())
	cue.InfoSection.Show()

}

func ExampleCue_Show() {
	dee := "/DA0AAAAAAAAAAAABQb/4zZ7tQAeAhxDVUVJAA6Gjz/TAAESy7EICAAAAAAA0/cuIgAAjFLk9Q=="
	cue := cuei.NewCue()
	cue.Decode(dee)
	cue.Show()
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
	t.Run("Cue_Show", func(t *testing.T) {
		ExampleCue_Show()
	})
}
