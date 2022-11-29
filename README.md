# cuei is a SCTE35 parser library, in Go.

- [x] Parses SCTE-35 Cues from MPEGTS or Bytes or Base64
- [x] Parses SCTE-35 Cues spread over multiple MPEGTS packets
- [x] Supports multi-packet PAT and PMT tables
- [x] Supports multiple MPEGTS Programs and multiple SCTE-35 streams

* [Install](#install-cuei)
* [Go Docs](https://pkg.go.dev/github.com/futzu/cuei)
* [Examples](#parse-base64-encoded-scte-35) 
	* [Parse Base64 encoded SCTE-35](#parse-base64-encoded-scte-35)      
	* [Parse SCTE-35 from MPEGTS](#parse-mpegts-video-for-scte35)
	* [Use cuei with another MPEGTS stream parser / demuxer](#use-cuei-with-another-mpegts-stream-parser--demuxer)
	* [Shadow a Cue Struct Method ( override ) ](#shadow-a-cue-struct-method)
	* [Shadow a Cue Method and call the Shadowed Method ( like super in python )](#call-a-shadowed-method)
	* [Use Dot Notation to access SCTE-35 Cue values](#use-dot-notation-to-access-scte-35-cue-values)


#### `install cuei`

```go
go install github.com/futzu/cuei@latest

```
#### `fetch cueidemo.go`
```rebol
curl http://iodisco.com/cueidemo.go -o cueidemo.go
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
```rebol
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


* cueidemo.go

```go 
package main

import (
	"os"
	"fmt"
	"github.com/futzu/cuei"
)

func main(){

	args := os.Args[1:]
	for i := range args{
		fmt.Printf( "\nNext File: %s\n\n",args[i] )
		var stream   cuei.Stream
		stream.Decode(args[i])
	}
} 


```
#### `parse base64 encoded SCTE-35`
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
#### `Use cuei with another MPEGTS stream parser / demuxer`
* Scte35Parser is for incorporating with another MPEGTS parser.
* Example
```go
	       import(
	           "github.com/futzu/cuei"
	           )
	       scte35parser := cuei.Scte35Parser{}

	       // Each time your parser/demuxer finds a SCTE-35 packet (stream type 0x86)

	       // do something like

	       cue := scte35parser.Parse(aScte35Packet)
	       if cue != nil {
	           // do something with the cue
	       }

```

*  If the packet is a partial Cue
	* The packet will be stored and aggregated with the next packet until complete.

* Single packet Cues and completed multiple packet Cues 
	* Are decoded into a Cue and returned.

#### `Shadow a Cue struct method`
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

#### Call a shadowed method
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
#### Use Dot notation to access SCTE-35 Cue values
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
