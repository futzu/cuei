
# cuei is a SCTE35 parser library in Go.

- [x] Parses SCTE-35 Cues from MPEGTS or Bytes or Base64
- [x] Parses SCTE-35 Cues spread over multiple MPEGTS packets
- [x] Supports multi-packet PAT and PMT tables
- [x] Supports multiple MPEGTS Programs and multiple SCTE-35 streams


###  
<details> <summary>Heads Up!</summary>

 
### I am about to make some backend changes, I need to reorganize things. 
### Most of you won't even notice.

#### Changes:
- [ ] Combine Gob and NBin into an external module
- [x]  Remove SCTE35Parser in favor of StreamParser
- [x]  Add CueParser for consistency
- [x]  Rename SpliceCommand to Command
- [x]  Rename SpliceDescriptor to Descriptor
- [ ]  .....

> A lot of this is to make the godocs easier to follow. It has to be done. 

	
</details>
	

* [Install](#install-cuei)

* [In a nutshell](#nutshell)

* [Go Docs](https://pkg.go.dev/github.com/futzu/cuei)

* [Examples](#parse-base64-encoded-scte-35) 

	* [Parse SCTE-35 from MPEGTS](#quick-demo)
	
	* [Parse Base64 encoded SCTE-35](#parse-base64-encoded-scte-35)
	      
	* [Use Dot Notation to access SCTE-35 Cue values](#use-dot-notation-to-access-scte-35-cue-values)
		
	* [Shadow a Cue Struct Method ( override ) ](#shadow-a-cue-struct-method)
	
	* [Shadow a Cue Method and call the Shadowed Method ( like super in python )](#call-a-shadowed-method)


# `Install cuei`

```go
go install github.com/futzu/cuei@latest

```
# `Nutshell`
| Use this        |   To do this                                                  |
|-----------------|---------------------------------------------------------------|
|[cuei.CueParser](https://github.com/futzu/cuei/blob/eac3a19eeb26/parsers.go)         | Parse SCTE-35 from a Base64 or Byte string.                   |
|[cuei.StreamParser](https://github.com/futzu/cuei/blob/main/parsers.go)      |        Parse SCTE35 Cues from a MPEGTS file.    |
|[cuei.StreamParser](https://github.com/futzu/cuei/blob/main/parsers.go) | Parse MPEGTS packets as an array of bytes, like from a network stream. | 


# `Quick Demo`

* cueidemo.go

```go 
package main

import (
        "os"
        "fmt"
        "github.com/futzu/cuei"
)

func main(){

        arg := os.Args[1]

        streamp := cuei.NewStreamParser()
        cues := streamp.ParseFile(arg)
        for _,cue := range cues {
        fmt.Printf("Command is a %v\n", cue.Command.Name)
        }

}



```

#### `build cueidemo`
```go
go build cueidemo.go
```
#### `parse mpegts video for scte35` 
```rebol
./cueidemo a_video_with_scte35.ts
```

#### `output`
```json
Next File: mpegts/out.ts

{
    "Name": "Splice Info Section",
    "TableID": "0xfc",
    "SectionSyntaxIndicator": false,
    "Private": false,
    "Reserved": "0x3",
    "SectionLength": 49,
    "ProtocolVersion": 0,
    "EncryptedPacket": false,
    "EncryptionAlgorithm": 0,
    "PtsAdjustment": 0,
    "CwIndex": "0x0",
    "Tier": "0xfff",
    "SpliceCommandLength": 20,
    "SpliceCommandType": 5,
    "DescriptorLoopLength": 12,
    "Command": {
        "Name": "Splice Insert",
        "CommandType": 5,
        "SpliceEventID": "0x5d",
        "OutOfNetworkIndicator": true,
        "ProgramSpliceFlag": true,
        "DurationFlag": true,
        "BreakDuration": 90.023266,
        "TimeSpecifiedFlag": true,
        "PTS": 38113.135577
    },
    "Descriptors": [
        {
            "Tag": 1,
            "Length": 10,
            "Identifier": "CUEI",
            "Name": "DTMF Descriptor",
            "PreRoll": 177,
            "DTMFCount": 4,
            "DTMFChars": 4186542473
        }
    ],
    "Packet": {
        "PacketNumber": 73885,
        "Pid": 515,
        "Program": 51,
        "Pcr": 38104.526277,
        "Pts": 38105.268588
    }
}


```
*  Use cuei.StreamParser for more fine-grained control of MPEGTS stream parsing. 
```go
package main

import (
	"fmt"
	"github.com/futzu/cuei"
	"os"
)

func main() {

	args := os.Args[1:] // Take multiple command line args
	for _, arg := range args {
		var cues []*cuei.Cue
		streamp := cuei.NewStreamParser() // New StreamParser for each file
		streamp.Quiet = true // suppress printing SCTE-35 messages 
		
		// you don't have to use a file
		// StreamParser.Parse takes a []byte as input
		
		file, err := os.Open(arg)
		if err != nil {
			break
		}
		defer file.Close()

		buffer := make([]byte, cuei.BufSz) // Parse in chunks
		for {
			_, err := file.Read(buffer)
			if err != nil {
				break
			}
			cues = streamp.Parse(buffer)   // StreamParser.Parse returns a [] *cuei.Cue 
			for _,cue := range cues {
			// do stuff with the cues like:
				cue.Show()
			// or
			fmt.Printf("Command is a %v\n", cue.Command.Name)
			}
			
		}
	}
}
```

# `Parse base64 encoded SCTE-35`
```go
package main

import (
	"fmt"
	"github.com/futzu/cuei"
)

func main(){

	var cue cuei.Cue
	data := cuei.DeB64("/DA7AAAAAAAAAP/wFAUAAAABf+/+AItfZn4AKTLgAAEAAAAWAhRDVUVJAAAAAX//AAApMuABACIBAIoXZrM=")
        cue.Decode(data) 
        fmt.Println("Cue as Json")
        cue.Show()
}
```
---
# `Shadow a Cue struct method`
```go
package main

import (
	"fmt"
	"github.com/futzu/cuei"
)

type Cue2 struct {
    cuei.Cue               		// Embed cuei.Cue
}
func (cue2 *Cue2) Show() {        	// Override Show
	fmt.Printf("%+v",cue2.Command)
}

func main(){

	var cue2 Cue2
	data := cuei.DeB64("/DA7AAAAAAAAAP/wFAUAAAABf+/+AItfZn4AKTLgAAEAAAAWAhRDVUVJAAAAAX//AAApMuABACIBAIoXZrM=")
        cue2.Decode(data) 
        cue2.Show()
	
}

```

# `Call a shadowed method`
```go
package main

import (
	"fmt"
	"github.com/futzu/cuei"
)

type Cue2 struct {
    cuei.Cue               		// Embed cuei.Cue
}
func (cue2 *Cue2) Show() {        	// Override Show

	fmt.Println("Cue2.Show()")
	fmt.Printf("%+v",cue2.Command) 
	
	fmt.Println("\n\ncuei.Cue.Show() from cue2.Show()")
	
	cue2.Cue.Show()			// Call the Show method from embedded cuei.Cue
}

func main(){

	var cue2 Cue2
	data := cuei.DeB64("/DA7AAAAAAAAAP/wFAUAAAABf+/+AItfZn4AKTLgAAEAAAAWAhRDVUVJAAAAAX//AAApMuABACIBAIoXZrM=")
        cue2.Decode(data) 
        cue2.Show()
	
}


```
# `Use Dot notation to access SCTE-35 Cue values`
```go

/**
Show  the packet PTS time and Splice Command Name of SCTE-35 Cues
in a MPEGTS stream.
**/


package main

import (
	"os"
	"fmt"
	"github.com/futzu/cuei"
)

func main() {

	args := os.Args[1:]
	for _,arg := range args {
		fmt.Printf("\nNext File: %s\n\n", arg)
		var stream cuei.Stream
		stream.Decode(arg)
		for _,c:= range stream.Cues {
			fmt.Printf("PTS: %v, Splice Command: %v\n",c.PacketData.Pts, c.Command.Name )
		}
	}
}
```
