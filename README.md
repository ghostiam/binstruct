# binstruct
Golang binary decoder to structure

**Warning: This project is under development, backward compatibility is not guaranteed.**

# Examples

[ZIP decoder](examples/zip) \
[PNG decoder](examples/png)

# Use

## For struct

### From file:
```go
file, err := os.Open("sample.png")
if err != nil {
    log.Fatal(err)
}

var png PNG
decoder := binstruct.NewDecoder(file, binary.BigEndian)
decoder.SetDebug(true) // you can enable the output of bytes read for debugging
err = decoder.Decode(&png)
if err != nil {
    log.Fatal(err)
}

spew.Dump(png)
```

### From bytes

```go
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
err := UnmarshalBE(data, &actual) // UnmarshalLE() or Unmarshal()
if err != nil {
    log.Fatal(err)
}
```

## or just use reader without mapping data into the structure

You can not use the functionality for mapping data into the structure, you can use the interface to get data from the stream (io.ReadSeeker)

[reader.go](reader.go)
```go
type Reader interface {
	io.ReadSeeker

	Peek(n int) ([]byte, error)

	ReadBytes(n int) (an int, b []byte, err error)
	ReadAll() ([]byte, error)

	ReadByte() (byte, error)
	ReadBool() (bool, error)

	ReadUint8() (uint8, error)
	ReadUint16() (uint16, error)
	ReadUint32() (uint32, error)
	ReadUint64() (uint64, error)

	ReadInt8() (int8, error)
	ReadInt16() (int16, error)
	ReadInt32() (int32, error)
	ReadInt64() (int64, error)

	ReadFloat32() (float32, error)
	ReadFloat64() (float64, error)

	Unmarshaler
}
```

# Struct tags

```go
type test struct {
	IgnoredField []byte `bin:"-"` // ignore field
	CallMethod   []byte `bin:"MethodName"` // Call method "MethodName"
	ReadLength   []byte `bin:"len:42"` // read 42 bytes
	
	// Offsets test binstruct_test.go:9
	Offset      []byte `bin:"offset:42"` // move to 42 bytes from current position
	OffsetStart []byte `bin:"offsetStart:42"` // move to 42 bytes from start position
	OffsetEnd   []byte `bin:"offsetEnd:-42"` // move to -42 bytes from end position
	OffsetStart []byte `bin:"offsetStart:42, offset:10"` // also worked and equally `offsetStart:52`

	CalcTagValue []byte `bin:"len:10+5+2+3"` // equally len:20, supported +,-,/,* (calculations are performed from left to right that is 2+2*2=8 not 6!!!)
	
	// You can refer to another field to get the value.
	DataLength              int // actual value
	ValueFromOtherField     int // `bin:"len:DataLength"`
	// also work calculations
	CalcValueFromOtherField int // `bin:"len:DataLength+10"`
} 

// Method can be:
func (*test) MethodName(r Reader) (FieldType, error) {
    ...
}
// or
func (*test) MethodName(r Reader) (error) {
    ...
}
```

See the tests and sample files for more information.

# License

MIT License