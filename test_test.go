// +build test

package binstruct

import (
	"bytes"
	"encoding/binary"
	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/require"
	"testing"
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

func (d *data) StringFunc(r ReadSeekPeeker) (string, error) {
	_, _, err := r.ReadBytes(1)
	if err != nil {
		return "", err
	}

	lenStr, err := r.ReadUint16()
	_, str, err := r.ReadBytes(int(lenStr))
	if err != nil {
		return "", err
	}

	return string(str), nil
}

func (d *data) MapFunc(r ReadSeekPeeker) error {
	s := make(map[int]string)

	for i := 0; i < 2; i++ {
		_, _, err := r.ReadBytes(1)
		if err != nil {
			return err
		}

		lenStr, err := r.ReadUint16()
		_, str, err := r.ReadBytes(int(lenStr))
		if err != nil {
			return err
		}

		s[i] = string(str)
	}

	d.Map = s
	return nil
}

func Test_Decoder(t *testing.T) {
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

	var want = data{
		StrLen: 5,
		Str:    "hello",
		Int:    10,
		ArrLen: 2,
		ISlice: []int{17, 34},
		IArr:   [2]int32{51, 68},
		SSlice: []string{"hi", "yay"},
		Map:    map[int]string{0: "hi", 1: "yay"},
		Custom: custom{
			ID:      255,
			TypeLen: 4,
			Type:    "test",
			B:       []byte{'h', 'i', '!'},
		},
	}

	var actual data

	decoder := NewDecoder(bytes.NewReader(b), binary.BigEndian)
	err := decoder.Decode(&actual)
	if err != nil {
		panic(err)
	}

	require.Equal(t, want, actual, spew.Sdump(actual))
}

type data2 struct {
	ID      int32
	Type    string `bin:"NullTerminatedString"`
	OtherID int32
}

func (*data2) NullTerminatedString(r ReadSeekPeeker) (string, error) {
	var b []byte

	for {
		readByte, err := r.ReadByte()
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

func Test_Decoder2(t *testing.T) {
	b := []byte{
		// ID
		0x00, 0x00, 0x00, 0x05,
		// Type as null-terminated string
		't', 'e', 's', 't', 0x00,
		// OtherID
		0xff, 0xff, 0xff, 0xf0,
	}
	var actual data2
	want := data2{
		ID:      5,
		Type:    "test",
		OtherID: -16,
	}

	decoder := NewDecoder(bytes.NewReader(b), binary.BigEndian)
	err := decoder.Decode(&actual)
	if err != nil {
		panic(err)
	}

	require.Equal(t, want, actual, spew.Sdump(actual))
}
