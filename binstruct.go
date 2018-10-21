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
	rs    io.ReadSeeker
	order binary.ByteOrder
	debug bool
}

func NewDecoder(rs io.ReadSeeker, order binary.ByteOrder) *Decoder {
	return &Decoder{
		rs:    rs,
		order: order,
		debug: false,
	}
}

func (dec *Decoder) SetDebug(debug bool) {
	dec.debug = debug
}

func (dec *Decoder) Decode(v interface{}) error {
	u := unmarshal{
		r: NewReader(dec.rs, dec.order, dec.debug),
	}

	return u.Unmarshal(v)
}
