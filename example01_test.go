package binstruct_test

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/davecgh/go-spew/spew"
	"github.com/ghostiam/binstruct"
)

type dataWithNullTerminatedString struct {
	ID      int32
	Type    string `bin:"NullTerminatedString"`
	OtherID int32
}

func (*dataWithNullTerminatedString) NullTerminatedString(r binstruct.Reader) (string, error) {
	var b []byte

	for {
		readByte, err := r.ReadByte()
		if binstruct.IsEOF(err) {
			break
		}
		if err != nil {
			return "", err
		}

		if readByte == 0x00 {
			break
		}

		b = append(b, readByte)
	}

	return string(b), nil
}

func Example_decoderDataWithNullTerminatedString() {
	b := []byte{
		// ID
		0x00, 0x00, 0x00, 0x05,
		// Type as null-terminated string
		't', 'e', 's', 't', 0x00,
		// OtherID
		0xff, 0xff, 0xff, 0xf0,
	}

	var actual dataWithNullTerminatedString

	decoder := binstruct.NewDecoder(bytes.NewReader(b), binary.BigEndian)
	err := decoder.Decode(&actual)
	if err != nil {
		panic(err)
	}

	fmt.Print(spew.Sdump(actual))

	// Output:
	// (binstruct_test.dataWithNullTerminatedString) {
	//  ID: (int32) 5,
	//  Type: (string) (len=4) "test",
	//  OtherID: (int32) -16
	// }
}
