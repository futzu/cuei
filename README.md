
 [Install](#install-cuei)  | [Go Docs](https://pkg.go.dev/github.com/futzu/cuei)  | [Examples](https://pkg.go.dev/github.com/futzu/cuei) | [cuei wins the SCTE-35 Parser Shoot Out](https://github.com/futzu/SCTE35-threefive/blob/master/speedtest.md)




### cuei is a SCTE-35 parser lib written in Go.<br>
*️⃣  CUEI is the FourCC / identifier for SCTE-35, that's where I got the name.
___

<br>


- [x] Parses SCTE-35 Cues from MPEGTS or Bytes or Base64 or Hex or Int or Octal or even Base 36.
- [x] Parses SCTE-35 Cues spread over multiple MPEGTS packets  
- [x] Supports multi-packet PAT and PMT tables  
- [x] Supports multiple MPEGTS Programs and multiple SCTE-35 streams 
- [x] Encodes Time Signals and Splice Inserts with Descriptors and Upids. 
 
---
#### 

### Want to parse an MPEGTS video and print the SCTE-35?  🛰️

### Do it in ten lines.
```go
package main                        

import (                              
        "os"                            
        "github.com/futzu/cuei"       
)                                    

func main(){                         
        arg := os.Args[1]             
        stream := cuei.NewStream()    
        stream.Decode(arg)           
}                                    
```
---
### Latest version is `One`  `Two`   `nine`
* Cyclomatic complexity score for v1.2.9 is __1.95__

# Documentation

* [Install](#install-cuei)  


* [Examples](https://pkg.go.dev/github.com/futzu/cuei) 

	* [Parse SCTE-35 from MPEGTS](#quick-demo) 
	
	* [Parse Base64 encoded SCTE-35](#parse-base64-encoded-scte-35) 
	      
	* [Use Dot Notation to access SCTE-35 Cue values](#use-dot-notation-to-access-scte-35-cue-values)
		
	* [Shadow a Cue Struct Method ( override ) ](#shadow-a-cue-struct-method)
	
	* [Shadow a Cue Method and call the Shadowed Method ( like super in python )](#call-a-shadowed-method)

  	* [Load a SCTE-35 Cue from JSON and Encode it](#load-json-and-encode)

### `Install cuei` 

```go
go get github.com/futzu/cuei@latest

```



### `Quick Demo` 

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

        stream := cuei.NewStream()
        cues := stream.Decode(arg)
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
```go
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
*  Use cuei.Stream.DecodeBytes for more fine-grained control of MPEGTS stream parsing. 
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
		stream := cuei.NewStream() // New StreamParser for each file
		stream.Quiet = true // suppress printing SCTE-35 messages 
		
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
			cues = stream.DecodeBytes(buffer)   // StreamParser.Parse returns a [] *cuei.Cue 
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

### `Parse base64 encoded SCTE-35`
```go
package main

import (
	"fmt"
	"github.com/futzu/cuei"
)

func main(){

	cue := cuei.NewCue()
	data := "/DA7AAAAAAAAAP/wFAUAAAABf+/+AItfZn4AKTLgAAEAAAAWAhRDVUVJAAAAAX//AAApMuABACIBAIoXZrM="
        cue.Decode(data) 
        fmt.Println("Cue as Json")
        cue.Show()
}
```

### `Shadow a Cue struct method`
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
	data := "/DA7AAAAAAAAAP/wFAUAAAABf+/+AItfZn4AKTLgAAEAAAAWAhRDVUVJAAAAAX//AAApMuABACIBAIoXZrM="
        cue2.Decode(data) 
        cue2.Show()
	
}

```

### `Call a shadowed method`
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
	data := "/DA7AAAAAAAAAP/wFAUAAAABf+/+AItfZn4AKTLgAAEAAAAWAhRDVUVJAAAAAX//AAApMuABACIBAIoXZrM="
        cue2.Decode(data) 
        cue2.Show()
	
}


```
### `Use Dot notation to access SCTE-35 Cue values`
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

	arg := os.Args[1]
	stream := cuei.NewStream()
	cues :=	stream.Decode(arg)
	for _,c := range cues {
		fmt.Printf("PTS: %v, Splice Command: %v\n",c.PacketData.Pts, c.Command.Name )
	}
}

```

### `Load JSON and Encode`
* cuei can accept SCTE-35 data as JSON and encode it to Base64, Bytes, or Hex string.
```go
package main

import (
	"fmt"
	"github.com/futzu/cuei"
)

func main() {

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
    "Crc32": "0xd7165c79"
}
`
	cue :=  cuei.Json2Cue(js)
	
	cue.AdjustPts(28.0)   	 // Apply pts adjustment
	
	fmt.Println("\nBytes:\n\t", cue.Encode())	// Bytes
	
	fmt.Println("\nBase64:\n\t",cue.Encode2B64())  	// Base64
	
	fmt.Println("\nHex:\n\t",cue.Encode2Hex()) 	// Hex

}


```
* Output
```smalltalk
Bytes:
	[252 48 42 0 0 0 38 115 192 255 255 240 15 5 0 0 22 58 127 207 254 127 12 79 115
	0 0 0 0 0 10 0 8 67 85 69 73 0 0 0 0 236 139 53 78]

Base64:
	 /DAqAAAAJnPA///wDwUAABY6f8/+fwxPcwAAAAAACgAIQ1VFSQAAAADsizVO

Hex:
	 0xfc302a0000002673c0fffff00f050000163a7fcffe7f0c4f7300000000000a00084355454900000000ec8b354e

```
