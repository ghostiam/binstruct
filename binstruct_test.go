package binstruct

import (
	"encoding/binary"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Offsets(t *testing.T) {
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

func Test_OffsetsMany(t *testing.T) {
	data := []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F}

	type dataStruct struct {
		ManyOffset  byte `bin:"offsetStart:2, offset:4, offset:-2"`
		CheckOffset byte `bin:"offsetStart:4"`
	}

	want := dataStruct{
		ManyOffset:  0x05,
		CheckOffset: 0x05,
	}

	var actual dataStruct
	err := UnmarshalBE(data, &actual)
	require.NoError(t, err)
	require.Equal(t, want, actual)
}

func Test_IntLE(t *testing.T) {
	data := []byte{
		0x01,
		0x02, 0x00,
		0x03, 0x00, 0x00, 0x00,
		0x04, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
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
	err := UnmarshalLE(data, &actual)
	require.NoError(t, err)
	require.Equal(t, want, actual)
}

func Test_IntXLE(t *testing.T) {
	// Test the weird sizes, i.e. not powers of 2
	data := []byte{
		0x03, 0x00, 0xf0, // -1048573
		0x0a, 0xfb, 0xc2, 0x10, 0xf0, // -68438263030
		0x0a, 0xfb, 0xc2, 0x10, 0xf0, 0x0c, // 14225212898058
		0x0a, 0xfb, 0xc2, 0x10, 0xf0, 0x0c, 0x7d, // 35198597301730058
	}

	type dataStruct struct {
		I3 int32 `bin:"len:3"`
		I5 int64 `bin:"len:5"`
		I6 int64 `bin:"len:6"`
		I7 int64 `bin:"len:7"`
	}

	want := dataStruct{
		I3: -1048573,
		I5: -68438263030,
		I6: 14225212898058,
		I7: 35198597301730058,
	}

	var actual dataStruct
	err := UnmarshalLE(data, &actual)
	require.NoError(t, err)
	require.Equal(t, want, actual)
}

func Test_IntBE(t *testing.T) {
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

func Test_IntXBE(t *testing.T) {
	// Test the weird sizes, i.e. not powers of 2
	data := []byte{
		0xf0, 0x00, 0x03, // -1048573
		0xf0, 0x10, 0xc2, 0xfb, 0x0a, // -68438263030
		0x0c, 0xf0, 0x10, 0xc2, 0xfb, 0x0a, // 14225212898058
		0x7d, 0x0c, 0xf0, 0x10, 0xc2, 0xfb, 0x0a, // 35198597301730058
	}

	type dataStruct struct {
		I3 int32 `bin:"len:3"`
		I5 int64 `bin:"len:5"`
		I6 int64 `bin:"len:6"`
		I7 int64 `bin:"len:7"`
	}

	want := dataStruct{
		I3: -1048573,
		I5: -68438263030,
		I6: 14225212898058,
		I7: 35198597301730058,
	}

	var actual dataStruct
	err := UnmarshalBE(data, &actual)
	require.NoError(t, err)
	require.Equal(t, want, actual)
}

func Test_IntBETag(t *testing.T) {
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

func Test_IntBEWithoutLenTag(t *testing.T) {
	data := []byte{
		0x01,
	}

	type dataStruct struct {
		I8 int
	}

	var actual dataStruct
	err := UnmarshalBE(data, &actual)
	require.EqualError(t, err, `failed set value to field "I8": need set tag with len or use int8/int16/int32/int64`)
	require.Equal(t, dataStruct{}, actual)
}

func Test_UintBE(t *testing.T) {
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

func Test_UintBETag(t *testing.T) {
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

func Test_UintBEWithoutLenTag(t *testing.T) {
	data := []byte{
		0x01,
	}

	type dataStruct struct {
		I8 uint
	}

	var actual dataStruct
	err := UnmarshalBE(data, &actual)
	require.EqualError(t, err, `failed set value to field "I8": need set tag with len or use uint8/uint16/uint32/uint64`)
	require.Equal(t, dataStruct{}, actual)
}

func Test_FloatBE(t *testing.T) {
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

func Test_Bool(t *testing.T) {
	data := []byte{
		0x00,
		0x01,
		0xFF,
	}

	type dataStruct struct {
		B1 bool
		B2 bool
		B3 bool
	}

	want := dataStruct{
		B1: false,
		B2: true,
		B3: true,
	}

	var actual dataStruct
	err := UnmarshalBE(data, &actual)
	require.NoError(t, err)
	require.Equal(t, want, actual)
}

func Test_Slice(t *testing.T) {
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

func Test_SliceWithoutLenTag(t *testing.T) {
	data := []byte{
		0x00, 0x01,
		0x00, 0x02,
		0x00, 0x03,
		0x00, 0x04,
	}

	type dataStruct struct {
		Arr []int16
	}

	var actual dataStruct
	err := UnmarshalBE(data, &actual)
	require.EqualError(t, err, `failed set value to field "Arr": need set tag with len for slice`)
	require.Equal(t, dataStruct{}, actual)
}

func Test_SliceOfSlice(t *testing.T) {
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

func Test_SliceOfSliceOfSlice(t *testing.T) {
	data := []byte{
		0x00, 0x01,
		0x00, 0x02,
		0x00, 0x03,
		0x00, 0x04,
		0x00, 0x05,
		0x00, 0x06,
		0x00, 0x07,
		0x00, 0x08,
	}

	type dataStruct struct {
		Arr [][][]int16 `bin:"len:2,[len:2,[len:2]]"`
	}

	want := dataStruct{
		Arr: [][][]int16{{{1, 2}, {3, 4}}, {{5, 6}, {7, 8}}},
	}

	var actual dataStruct
	err := UnmarshalBE(data, &actual)
	require.NoError(t, err)
	require.Equal(t, want, actual)
}

func Test_Array(t *testing.T) {
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

func Test_ArrayOfArray(t *testing.T) {
	data := []byte{
		0x00, 0x01,
		0x00, 0x02,
		0x00, 0x03,
		0x00, 0x04,
	}

	type dataStruct struct {
		Arr [2][2]int16
	}

	want := dataStruct{
		Arr: [2][2]int16{{1, 2}, {3, 4}},
	}

	var actual dataStruct
	err := UnmarshalBE(data, &actual)
	require.NoError(t, err)
	require.Equal(t, want, actual)
}

func Test_ArrayOfArrayOfArray(t *testing.T) {
	data := []byte{
		0x00, 0x01,
		0x00, 0x02,
		0x00, 0x03,
		0x00, 0x04,
		0x00, 0x05,
		0x00, 0x06,
		0x00, 0x07,
		0x00, 0x08,
	}

	type dataStruct struct {
		Arr [2][2][2]int16
	}

	want := dataStruct{
		Arr: [2][2][2]int16{{{1, 2}, {3, 4}}, {{5, 6}, {7, 8}}},
	}

	var actual dataStruct
	err := UnmarshalBE(data, &actual)
	require.NoError(t, err)
	require.Equal(t, want, actual)
}

func Test_ByteArray(t *testing.T) {
	data := []byte{0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F}

	type dataStruct struct {
		B [4]byte
	}

	want := dataStruct{
		B: [4]byte{0x0A, 0x0B, 0x0C, 0x0D},
	}

	var actual dataStruct
	err := UnmarshalBE(data, &actual)
	require.NoError(t, err)
	require.Equal(t, want, actual)
}

func Test_ByteSlice(t *testing.T) {
	data := []byte{0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F}

	type dataStruct struct {
		B []byte `bin:"len:4"`
	}

	want := dataStruct{
		B: []byte{0x0A, 0x0B, 0x0C, 0x0D},
	}

	var actual dataStruct
	err := UnmarshalBE(data, &actual)
	require.NoError(t, err)
	require.Equal(t, want, actual)
}

func Test_StringEmpty(t *testing.T) {
	data := []byte{}

	type dataStruct struct {
		Str string `bin:"len:0"`
	}

	want := dataStruct{
		Str: "",
	}

	var actual dataStruct
	err := UnmarshalBE(data, &actual)
	require.NoError(t, err)
	require.Equal(t, want, actual)
}

func Test_String(t *testing.T) {
	data := []byte{'h', 'e', 'l', 'l', 'o'}

	type dataStruct struct {
		Str string `bin:"len:5"`
	}

	want := dataStruct{
		Str: "hello",
	}

	var actual dataStruct
	err := UnmarshalBE(data, &actual)
	require.NoError(t, err)
	require.Equal(t, want, actual)
}

func Test_StringWithoutLenTag(t *testing.T) {
	data := []byte{'h', 'e', 'l', 'l', 'o'}

	type dataStruct struct {
		Str string
	}

	var actual dataStruct
	err := UnmarshalBE(data, &actual)
	require.EqualError(t, err, `failed set value to field "Str": need set tag with len for string`)
	require.Equal(t, dataStruct{}, actual)
}

func Test_StringWithLenFromField(t *testing.T) {
	data := []byte{0x00, 0x05, 'h', 'e', 'l', 'l', 'o'}

	type dataStruct struct {
		StrLen int16
		Str    string `bin:"len:StrLen"`
	}

	want := dataStruct{
		StrLen: 5,
		Str:    "hello",
	}

	var actual dataStruct
	err := UnmarshalBE(data, &actual)
	require.NoError(t, err)
	require.Equal(t, want, actual)
}

type dataCustomMethod1Struct struct {
	Custom map[string]string `bin:"CustomMap"`
}

func (d *dataCustomMethod1Struct) CustomMap(r Reader) error {
	m := make(map[string]string)

	lenMap, err := r.ReadInt8()
	if err != nil {
		return err
	}

	for i := 0; i < int(lenMap); i++ {
		_, name, err := r.ReadBytes(1)
		if err != nil {
			return err
		}

		_, value, err := r.ReadBytes(1)
		if err != nil {
			return err
		}

		m[string(name)] = string(value)
	}

	d.Custom = m
	return nil
}

func Test_CustomMethod1(t *testing.T) {
	data := []byte{0x03, 'q', 'w', 'e', 'r', 't', 'y'}

	want := dataCustomMethod1Struct{
		Custom: map[string]string{"q": "w", "e": "r", "t": "y"},
	}

	var actual dataCustomMethod1Struct
	err := UnmarshalBE(data, &actual)
	require.NoError(t, err)
	require.Equal(t, want, actual)
}

type dataCustomMethod2Struct struct {
	Custom [2]map[string]string `bin:"len:2,[CustomMap]"`
}

func (*dataCustomMethod2Struct) CustomMap(r Reader) (map[string]string, error) {
	m := make(map[string]string)

	lenMap, err := r.ReadInt8()
	if err != nil {
		return nil, err
	}

	for i := 0; i < int(lenMap); i++ {
		_, name, err := r.ReadBytes(1)
		if err != nil {
			return nil, err
		}

		_, value, err := r.ReadBytes(1)
		if err != nil {
			return nil, err
		}

		m[string(name)] = string(value)
	}

	return m, nil
}

func Test_CustomMethod2(t *testing.T) {
	data := []byte{
		0x03, 'q', 'w', 'e', 'r', 't', 'y',
		0x03, 'a', 's', 'd', 'f', 'g', 'h',
	}

	want := dataCustomMethod2Struct{
		Custom: [2]map[string]string{
			{"q": "w", "e": "r", "t": "y"},
			{"a": "s", "d": "f", "g": "h"},
		},
	}

	var actual dataCustomMethod2Struct
	err := UnmarshalBE(data, &actual)
	require.NoError(t, err)
	require.Equal(t, want, actual)
}

func Test_CustomMethodNotExist(t *testing.T) {
	data := []byte{}

	type dataCustomMethod3Struct struct {
		Custom string `bin:"CustomMethodNotExist"`
	}

	var actual dataCustomMethod3Struct
	err := UnmarshalBE(data, &actual)
	require.EqualError(t, err, `failed set value to field "Custom": 
failed call method, expected methods:
	func (*dataCustomMethod3Struct) CustomMethodNotExist(r binstruct.Reader) error {} 
or
	func (*dataCustomMethod3Struct) CustomMethodNotExist(r binstruct.Reader) (string, error) {}
`)
	require.Equal(t, dataCustomMethod3Struct{}, actual)
}

func Test_InvalidType(t *testing.T) {
	data := []byte{}

	type dataStruct struct {
		Invalid interface{}
	}

	var actual dataStruct
	err := UnmarshalBE(data, &actual)
	require.EqualError(t, err, `failed set value to field "Invalid": type "interface" not supported`)
	require.Equal(t, dataStruct{}, actual)
}

type CustomMethodFromParent struct {
	Pin struct {
		Checksum uint16 `bin:"CustomMethodFromParent,len:2"`

		Pin struct {
			Checksum uint16 `bin:"CustomMethodFromParent,len:2"`

			Pin struct {
				Checksum uint16 `bin:"CustomMethodFromParent,len:2"`
			}
		}
	}
}

func (*CustomMethodFromParent) CustomMethodFromParent(r Reader) (uint16, error) {
	var out uint16
	if err := binary.Read(r, binary.LittleEndian, &out); err != nil {
		return 0, err
	}
	return out, nil
}

func Test_CustomMethodFromParent_Issue4(t *testing.T) {
	data := []byte{0x01, 0x00, 0x02, 0x00, 0x03, 0x00}

	want := CustomMethodFromParent{
		Pin: struct {
			Checksum uint16 `bin:"CustomMethodFromParent,len:2"`
			Pin      struct {
				Checksum uint16 `bin:"CustomMethodFromParent,len:2"`
				Pin      struct {
					Checksum uint16 `bin:"CustomMethodFromParent,len:2"`
				}
			}
		}{
			Checksum: 1,
			Pin: struct {
				Checksum uint16 `bin:"CustomMethodFromParent,len:2"`
				Pin      struct {
					Checksum uint16 `bin:"CustomMethodFromParent,len:2"`
				}
			}{
				Checksum: 2,
				Pin: struct {
					Checksum uint16 `bin:"CustomMethodFromParent,len:2"`
				}{
					Checksum: 3,
				},
			},
		},
	}

	var actual CustomMethodFromParent
	err := UnmarshalLE(data, &actual)
	require.NoError(t, err)
	require.Equal(t, want, actual)
}

type LeAndBeInOneStruct struct {
	UInt16             uint16
	UInt16LE           uint16 `bin:"le"`
	UInt16BE           uint16 `bin:"be"`
	UInt16WithLEReader uint16 `bin:"ParseUInt16WithLEReader,le"`
	UInt16WithBEReader uint16 `bin:"be,ParseUInt16WithBEReader"`
	UInt16Check        uint16
}

func (*LeAndBeInOneStruct) ParseUInt16WithLEReader(r Reader) (uint16, error) {
	return r.ReadUint16()
}

func (*LeAndBeInOneStruct) ParseUInt16WithBEReader(r Reader) (uint16, error) {
	return r.ReadUint16()
}

func Test_LeAndBeInOneStruct(t *testing.T) {
	data := []byte{0x01, 0x00, 0x02, 0x00, 0x03, 0x00, 0x04, 0x00, 0x05, 0x00, 0x06, 0x00}

	wantLE := LeAndBeInOneStruct{
		UInt16:             1,
		UInt16LE:           2,
		UInt16BE:           768,
		UInt16WithLEReader: 4,
		UInt16WithBEReader: 1280,
		UInt16Check:        6,
	}

	wantBE := LeAndBeInOneStruct{
		UInt16:             256,
		UInt16LE:           2,
		UInt16BE:           768,
		UInt16WithLEReader: 4,
		UInt16WithBEReader: 1280,
		UInt16Check:        1536,
	}

	var actual LeAndBeInOneStruct
	err := UnmarshalLE(data, &actual)
	require.NoError(t, err)
	require.Equal(t, wantLE, actual)

	err = UnmarshalBE(data, &actual)
	require.NoError(t, err)
	require.Equal(t, wantBE, actual)
}
