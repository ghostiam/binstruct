[![Go Report Card](https://goreportcard.com/badge/github.com/ghostiam/binstruct)](https://goreportcard.com/report/github.com/ghostiam/binstruct) [![CodeCov](https://codecov.io/gh/ghostiam/binstruct/branch/master/graph/badge.svg)](https://codecov.io/gh/ghostiam/binstruct) [![GoDoc](https://godoc.org/github.com/ghostiam/binstruct?status.svg)](https://godoc.org/github.com/ghostiam/binstruct) [![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://github.com/ghostiam/binstruct/blob/master/LICENSE)

# binstruct
Golang binary decoder to structure

# Install
```go get -u github.com/ghostiam/binstruct```

# Examples

[ZIP decoder](examples/zip) \
[PNG decoder](examples/png)

# Use

## For struct

### From file or other io.ReadSeeker:
```go
package main

import (
	"encoding/binary"
	"fmt"
	"log"
	"os"

	"github.com/ghostiam/binstruct"
)

func main() {
	file, err := os.Open("testdata/file.bin")
	if err != nil {
		log.Fatal(err)
	}

	type dataStruct struct {
		Arr []int16 `bin:"len:4"`
	}

	var actual dataStruct
	decoder := binstruct.NewDecoder(file, binary.BigEndian)
	// decoder.SetDebug(true) // you can enable the output of bytes read for debugging
	err = decoder.Decode(&actual)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%+v", actual)

	// Output:
	// {Arr:[1 2 3 4]}
}
```

### From bytes

```go
package main

import (
	"fmt"
	"log"
	
	"github.com/ghostiam/binstruct"
)

func main() {
	data := []byte{
		0x00, 0x01,
		0x00, 0x02,
		0x00, 0x03,
		0x00, 0x04,
	}

	type dataStruct struct {
		Arr []int16 `bin:"len:4"`
	}

	var actual dataStruct
	err := binstruct.UnmarshalBE(data, &actual) // UnmarshalLE() or Unmarshal()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%+v", actual)

	// Output: {Arr:[1 2 3 4]}
}
```

## or just use reader without mapping data into the structure

You can not use the functionality for mapping data into the structure, you can use the interface to get data from the stream (io.ReadSeeker)

[reader.go](reader.go)
```go
type Reader interface {
	io.ReadSeeker

	// Peek returns the next n bytes without advancing the reader.
	Peek(n int) ([]byte, error)

	// ReadBytes reads up to n bytes. It returns the number of bytes
	// read, bytes and any error encountered.
	ReadBytes(n int) (an int, b []byte, err error)
	// ReadAll reads until an error or EOF and returns the data it read.
	ReadAll() ([]byte, error)

	// ReadByte read and return one byte
	ReadByte() (byte, error)
	// ReadBool read one byte and return boolean value
	ReadBool() (bool, error)

	// ReadUint8 read one byte and return uint8 value
	ReadUint8() (uint8, error)
	// ReadUint16 read two bytes and return uint16 value
	ReadUint16() (uint16, error)
	// ReadUint32 read four bytes and return uint32 value
	ReadUint32() (uint32, error)
	// ReadUint64 read eight bytes and return uint64 value
	ReadUint64() (uint64, error)
	// ReadUintX read X bytes and return uint64 value
	ReadUintX(x int) (uint64, error)

	// ReadInt8 read one byte and return int8 value
	ReadInt8() (int8, error)
	// ReadInt16 read two bytes and return int16 value
	ReadInt16() (int16, error)
	// ReadInt32 read four bytes and return int32 value
	ReadInt32() (int32, error)
	// ReadInt64 read eight bytes and return int64 value
	ReadInt64() (int64, error)
	// ReadIntX read X bytes and return int64 value
	ReadIntX(x int) (int64, error)

	// ReadFloat32 read four bytes and return float32 value
	ReadFloat32() (float32, error)
	// ReadFloat64 read eight bytes and return float64 value
	ReadFloat64() (float64, error)

	// Unmarshal parses the binary data and stores the result
	// in the value pointed to by v.
	Unmarshal(v interface{}) error

	// WithOrder changes the byte order for the new Reader
	WithOrder(order binary.ByteOrder) Reader
}
```

Example:
```go
package main

import (
	"encoding/binary"
	"fmt"
	"log"
	
	"github.com/ghostiam/binstruct"
)

func main() {
	data := []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F}

	reader := binstruct.NewReaderFromBytes(data, binary.BigEndian, false)

	i16, err := reader.ReadInt16()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(i16)

	i32, err := reader.ReadInt32()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(i32)

	b, err := reader.Peek(4)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Peek bytes: %#v\n", b)

	an, b, err := reader.ReadBytes(4)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Read %d bytes: %#v\n", an, b)

	other, err := reader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Read all: %#v\n", other)

	// Output:
	// 258
	// 50595078
	// Peek bytes: []byte{0x7, 0x8, 0x9, 0xa}
	// Read 4 bytes: []byte{0x7, 0x8, 0x9, 0xa}
	// Read all: []byte{0xb, 0xc, 0xd, 0xe, 0xf}
}
```

# Decode to fields

```go
type test struct {
	// Read 1 byte
	Field bool
	Field byte
	Field [1]byte
	Field int8
	Field uint8

	// Read 2 bytes
	Field int16
	Field uint16
	Field [2]byte

	// Read 4 bytes
	Field int32
	Field uint32
	Field [4]byte

	// Read 8 bytes
	Field int64
	Field uint64
	Field [8]byte

	// You can override length
	Field int64 `bin:"len:2"`
	// Or even use very weird byte lengths for int
	Field int64 `bin:"len:3"`
	Field int64 `bin:"len:5"`
	Field int64 `bin:"len:7"`

	// Fields of type int, uint and string are not read automatically 
	// because the size is not known, you need to set it manually
	Field int    `bin:"len:2"`
	Field uint   `bin:"len:4"`
	Field string `bin:"len:42"`
	
	// Can read arrays and slices
	Array [2]int32              // read 8 bytes (4+4byte for 2 int32)
	Slice []int32 `bin:"len:2"` // read 8 bytes (4+4byte for 2 int32)
	
	// Also two-dimensional slices work (binstruct_test.go:307 Test_SliceOfSlice)
	Slice2D [][]int32 `bin:"len:2,[len:2]"`
	// and even three-dimensional slices (binstruct_test.go:329 Test_SliceOfSliceOfSlice)
	Slice3D [][][]int32 `bin:"len:2,[len:2,[len:2]]"`
	
	// Structures and embedding are also supported.
	Struct struct {
		...
	}
	OtherStruct Other
	Other // embedding
}

type Other struct {
	...
}
```

# Tags

```go
type test struct {
	IgnoredField []byte `bin:"-"`          // ignore field
	CallMethod   []byte `bin:"MethodName"` // Call method "MethodName"
	ReadLength   []byte `bin:"len:42"`     // read 42 bytes

	// Offsets test binstruct_test.go:9
	Offset      byte `bin:"offset:42"`      // move to 42 bytes from current position and read byte
	OffsetStart byte `bin:"offsetStart:42"` // move to 42 bytes from start position and read byte
	OffsetEnd   byte `bin:"offsetEnd:-42"`  // move to -42 bytes from end position and read byte
	OffsetStart byte `bin:"offsetStart:42, offset:10"` // also worked and equally `offsetStart:52`

	// Calculations supported +,-,/,* and are performed from left to right that is 2+2*2=8 not 6!!!
	CalcTagValue []byte `bin:"len:10+5+2+3"` // equally len:20

	// You can refer to another field to get the value.
	DataLength              int    // actual length
	ValueFromOtherField     string `bin:"len:DataLength"`
	CalcValueFromOtherField string `bin:"len:DataLength+10"` // also work calculations

	// You can change the byte order directly from the tag
	UInt16LE uint16 `bin:"le"`
	UInt16BE uint16 `bin:"be"`
	// Or when you call the method, it will contain the Reader with the byte order you need
	CallMethodWithLEReader uint16 `bin:"MethodNameWithLEReader,le"`
	CallMethodWithBEReader uint16 `bin:"be,MethodNameWithBEReader"`
} 

// Method can be:
func (*test) MethodName(r binstruct.Reader) (error) {}
// or
func (*test) MethodName(r binstruct.Reader) (FieldType, error) {}
```

See the tests and examples for more information.

# License

MIT License
