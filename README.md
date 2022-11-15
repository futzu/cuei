# cuei

* [Install](#install-cuei)
* Examples
	* [Parse Base64 encoded SCTE-35](#parse-base64-encoded-scte-35)      
	* [Parse SCTE-35 from MPEGTS](#parse-mpegts-video-for-scte35)
	* [Shadow a Cue Struct Method ( override ) ](#shadow-a-cue-struct-method)
	* [Shadow a Cue Method and call the Shadowed Method ( like super in python )](#call-a-shadowed-method)




#### `install cuei`

```go
go install github.com/futzu/cuei
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
        fmt.Printf("\nCue.Command\n\n%+v\n",cue.Command)          
        fmt.Printf("\nCue.Descriptors[0]\n\n%+v\n",cue.Descriptors[0])         
}
```
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

