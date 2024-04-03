package binstruct_test

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"

	"github.com/davecgh/go-spew/spew"

	"github.com/ghostiam/binstruct"
)

type custom struct {
	ID      int16
	_       [1]byte
	TypeLen int16
	Type    string `bin:"len:TypeLen"`
	B       []byte `bin:"len:3"`
}

type data struct {
	StrLen int    `bin:"len:2,offset:1"`
	Str    string `bin:"len:StrLen"`
	Int    int32  `bin:"len:2"`
	ArrLen uint16
	ISlice []int `bin:"len:ArrLen,[len:4]"`
	IArr   [2]int32
	SSlice []string       `bin:"len:ArrLen,[StringFunc]"`
	Map    map[int]string `bin:"MapFunc"`

	Skip []byte `bin:"-"`

	Custom custom
}

func (d *data) StringFunc(r binstruct.Reader) (string, error) {
	_, _, err := r.ReadBytes(1)
	if err != nil {
		return "", err
	}

	lenStr, err := r.ReadUint16()
	if err != nil {
		return "", err
	}

	_, str, err := r.ReadBytes(int(lenStr))
	if err != nil {
		return "", err
	}

	return string(str), nil
}

func (d *data) MapFunc(r binstruct.Reader) error {
	s := make(map[int]string)

	for i := 0; i < 2; i++ {
		_, _, err := r.ReadBytes(1)
		if err != nil {
			return err
		}

		lenStr, err := r.ReadUint16()
		if err != nil {
			return err
		}

		_, str, err := r.ReadBytes(int(lenStr))
		if err != nil {
			return err
		}

		s[i] = string(str)
	}

	d.Map = s
	return nil
}

func Example_decodeCustom() {
	var b = []byte{
		's', 0x00, 0x05, 'h', 'e', 'l', 'l', 'o', // string
		// Int
		0x00, 0x0A,
		// ArrLen
		0x00, 0x02,
		// ISlice
		0x00, 0x00, 0x00, 0x11, // [0]int
		0x00, 0x00, 0x00, 0x22, // [1]int
		// IArr
		0x00, 0x00, 0x00, 0x33, // [0]int
		0x00, 0x00, 0x00, 0x44, // [1]int
		// SSlice
		's', 0x00, 0x02, 'h', 'i', // [0]string
		's', 0x00, 0x03, 'y', 'a', 'y', // [1]string
		// Map
		's', 0x00, 0x02, 'h', 'i', // [0]string
		's', 0x00, 0x03, 'y', 'a', 'y', // [1]string
		// Custom
		0x00, 0xff, // int
		0xff,       // skip
		0x00, 0x04, // str len
		't', 'e', 's', 't', // string
		'h', 'i', '!', // bytes
	}

	var actual data

	decoder := binstruct.NewDecoder(bytes.NewReader(b), binary.BigEndian)
	err := decoder.Decode(&actual)
	if err != nil {
		log.Fatal(err)
	}

	spewCfg := spew.NewDefaultConfig()
	spewCfg.SortKeys = true
	fmt.Print(spewCfg.Sdump(actual))

	// Output: (binstruct_test.data) {
	//  StrLen: (int) 5,
	//  Str: (string) (len=5) "hello",
	//  Int: (int32) 10,
	//  ArrLen: (uint16) 2,
	//  ISlice: ([]int) (len=2 cap=2) {
	//   (int) 17,
	//   (int) 34
	//  },
	//  IArr: ([2]int32) (len=2 cap=2) {
	//   (int32) 51,
	//   (int32) 68
	//  },
	//  SSlice: ([]string) (len=2 cap=2) {
	//   (string) (len=2) "hi",
	//   (string) (len=3) "yay"
	//  },
	//  Map: (map[int]string) (len=2) {
	//   (int) 0: (string) (len=2) "hi",
	//   (int) 1: (string) (len=3) "yay"
	//  },
	//  Skip: ([]uint8) <nil>,
	//  Custom: (binstruct_test.custom) {
	//   ID: (int16) 255,
	//   _: ([1]uint8) (len=1 cap=1) {
	//    00000000  00                                                |.|
	//   },
	//   TypeLen: (int16) 4,
	//   Type: (string) (len=4) "test",
	//   B: ([]uint8) (len=3 cap=3) {
	//    00000000  68 69 21                                          |hi!|
	//   }
	//  }
	// }
}
