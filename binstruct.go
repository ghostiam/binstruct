package binstruct

import (
	"bytes"
	"encoding/binary"
	"io"
)

func UnmarshalLE(data []byte, v interface{}) error {
	return Unmarshal(data, binary.LittleEndian, v)
}

func UnmarshalBE(data []byte, v interface{}) error {
	return Unmarshal(data, binary.BigEndian, v)
}

func Unmarshal(data []byte, order binary.ByteOrder, v interface{}) error {
	return NewDecoder(bytes.NewReader(data), order).Decode(v)
}

type Decoder struct {
	r     io.ReadSeeker
	order binary.ByteOrder
	debug bool
}

func NewDecoder(r io.ReadSeeker, order binary.ByteOrder) *Decoder {
	return &Decoder{
		r:     r,
		order: order,
		debug: false,
	}
}

func (dec *Decoder) SetDebug(debug bool) {
	dec.debug = debug
}

func (dec *Decoder) Decode(v interface{}) error {
	return NewReader(dec.r, dec.order, dec.debug).Unmarshal(v)
}
