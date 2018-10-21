package binstruct

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_UnmarshalOffsets(t *testing.T) {
	data := []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F}

	type dataStruct struct {
		First  byte
		Second byte
		Last   byte `bin:"offsetEnd:-1"`

		OffsetFromStart5    byte `bin:"offsetStart:5"`
		OffsetFromStart10   byte `bin:"offsetStart:10"`
		OffsetFromEnd8      byte `bin:"offsetEnd:-8"`
		AfterOffsetFromEnd8 byte

		FirstAgain            byte `bin:"offsetStart:0"`
		SecondAgain           byte
		Skip1AfterSecondAgain byte `bin:"offset:1"`
	}

	want := dataStruct{
		First:                 0x01,
		Second:                0x02,
		Last:                  0x0f,
		OffsetFromStart5:      0x06,
		OffsetFromStart10:     0x0B,
		OffsetFromEnd8:        0x08,
		AfterOffsetFromEnd8:   0x09,
		FirstAgain:            0x01,
		SecondAgain:           0x02,
		Skip1AfterSecondAgain: 0x04,
	}

	var actual dataStruct
	err := UnmarshalBE(data, &actual)
	require.NoError(t, err)
	require.Equal(t, want, actual)
}

func Test_UnmarshalIntBE(t *testing.T) {
	data := []byte{
		0x01,
		0x00, 0x02,
		0x00, 0x00, 0x00, 0x03,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x04,
	}

	type dataStruct struct {
		I8  int8
		I16 int16
		I32 int32
		I64 int64
	}

	want := dataStruct{
		I8: 1, I16: 2, I32: 3, I64: 4,
	}

	var actual dataStruct
	err := UnmarshalBE(data, &actual)
	require.NoError(t, err)
	require.Equal(t, want, actual)
}

func Test_UnmarshalIntBETag(t *testing.T) {
	data := []byte{
		0x01,
		0x00, 0x02,
		0x00, 0x00, 0x00, 0x03,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x04,
	}

	type dataStruct struct {
		I8  int `bin:"len:1"`
		I16 int `bin:"len:2"`
		I32 int `bin:"len:4"`
		I64 int `bin:"len:8"`
	}

	want := dataStruct{
		I8: 1, I16: 2, I32: 3, I64: 4,
	}

	var actual dataStruct
	err := UnmarshalBE(data, &actual)
	require.NoError(t, err)
	require.Equal(t, want, actual)
}

func Test_UnmarshalUintBE(t *testing.T) {
	data := []byte{
		0x01,
		0x00, 0x02,
		0x00, 0x00, 0x00, 0x03,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x04,
	}

	type dataStruct struct {
		I8  uint8
		I16 uint16
		I32 uint32
		I64 uint64
	}

	want := dataStruct{
		I8: 1, I16: 2, I32: 3, I64: 4,
	}

	var actual dataStruct
	err := UnmarshalBE(data, &actual)
	require.NoError(t, err)
	require.Equal(t, want, actual)
}

func Test_UnmarshalUintBETag(t *testing.T) {
	data := []byte{
		0x01,
		0x00, 0x02,
		0x00, 0x00, 0x00, 0x03,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x04,
	}

	type dataStruct struct {
		I8  uint `bin:"len:1"`
		I16 uint `bin:"len:2"`
		I32 uint `bin:"len:4"`
		I64 uint `bin:"len:8"`
	}

	want := dataStruct{
		I8: 1, I16: 2, I32: 3, I64: 4,
	}

	var actual dataStruct
	err := UnmarshalBE(data, &actual)
	require.NoError(t, err)
	require.Equal(t, want, actual)
}

func Test_UnmarshalFloatBE(t *testing.T) {
	data := []byte{
		0x40, 0x49, 0x0f, 0xdb,
		0x40, 0x09, 0x21, 0xfb, 0x54, 0x44, 0x2d, 0x18,
	}

	type dataStruct struct {
		F32 float32
		F64 float64
	}

	want := dataStruct{
		F32: 3.1415927,
		F64: 3.141592653589793,
	}

	var actual dataStruct
	err := UnmarshalBE(data, &actual)
	require.NoError(t, err)
	require.Equal(t, want, actual)
}

func Test_UnmarshalSlice(t *testing.T) {
	data := []byte{
		0x00, 0x01,
		0x00, 0x02,
		0x00, 0x03,
		0x00, 0x04,
	}

	type dataStruct struct {
		Arr []int16 `bin:"len:4"`
	}

	want := dataStruct{
		Arr: []int16{1, 2, 3, 4},
	}

	var actual dataStruct
	err := UnmarshalBE(data, &actual)
	require.NoError(t, err)
	require.Equal(t, want, actual)
}

func Test_UnmarshalSliceOfSlice(t *testing.T) {
	data := []byte{
		0x00, 0x01,
		0x00, 0x02,
		0x00, 0x03,
		0x00, 0x04,
	}

	type dataStruct struct {
		Arr [][]int16 `bin:"len:2,[len:2]"`
	}

	want := dataStruct{
		Arr: [][]int16{{1, 2}, {3, 4}},
	}

	var actual dataStruct
	err := UnmarshalBE(data, &actual)
	require.NoError(t, err)
	require.Equal(t, want, actual)
}

func Test_UnmarshalArray(t *testing.T) {
	data := []byte{
		0x00, 0x01,
		0x00, 0x02,
		0x00, 0x03,
		0x00, 0x04,
	}

	type dataStruct struct {
		Arr [4]int16
	}

	want := dataStruct{
		Arr: [4]int16{1, 2, 3, 4},
	}

	var actual dataStruct
	err := UnmarshalBE(data, &actual)
	require.NoError(t, err)
	require.Equal(t, want, actual)
}
