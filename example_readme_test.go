package binstruct_test

import (
	"encoding/binary"
	"fmt"
	"log"
	"os"

	"github.com/ghostiam/binstruct"
)

func Example_readmeFromIOReadSeeker() {
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

func Example_readmeFromIOReadSeekerWithDebuging() {
	file, err := os.Open("testdata/file.bin")
	if err != nil {
		log.Fatal(err)
	}

	type dataStruct struct {
		Arr []int16 `bin:"len:4"`
	}

	var actual dataStruct
	decoder := binstruct.NewDecoder(file, binary.BigEndian)
	decoder.SetDebug(true) // you can enable the output of bytes read for debugging
	err = decoder.Decode(&actual)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%+v", actual)

	// Output:
	// Read(want: 2|actual: 2): ([]uint8) (len=2 cap=2) {
	//  00000000  00 01                                             |..|
	// }
	// Read(want: 2|actual: 2): ([]uint8) (len=2 cap=2) {
	//  00000000  00 02                                             |..|
	// }
	// Read(want: 2|actual: 2): ([]uint8) (len=2 cap=2) {
	//  00000000  00 03                                             |..|
	// }
	// Read(want: 2|actual: 2): ([]uint8) (len=2 cap=2) {
	//  00000000  00 04                                             |..|
	// }
	// {Arr:[1 2 3 4]}
}

func Example_readmeFromBytes() {
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

	// Output:
	// {Arr:[1 2 3 4]}
}

func Example_readmeReader() {
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
