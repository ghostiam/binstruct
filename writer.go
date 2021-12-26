package binstruct

import (
	"bytes"
	"encoding/binary"
	"io"
)

type Writer interface {
	io.Writer

	// WriteBool write one byte boolean.
	WriteBool(v bool) error

	// WriteUint8 write one byte uint8.
	WriteUint8(v uint8) error
	// WriteUint16 write two bytes uint16.
	WriteUint16(v uint16) error
	// WriteUint32 write four bytes uint32.
	WriteUint32(v uint32) error
	// WriteUint64 write eight bytes uint64.
	WriteUint64(v uint64) error

	// WriteInt8 write one byte int8.
	WriteInt8(v int8) error
	// WriteInt16 write two bytes int16.
	WriteInt16(v int16) error
	// WriteInt32 write four bytes int32.
	WriteInt32(v int32) error
	// WriteInt64 write eight bytes int64.
	WriteInt64(v int64) error

	// WriteFloat32 write four bytes float32.
	WriteFloat32(v float32) error
	// WriteFloat64 write eight bytes float64.
	WriteFloat64(v float64) error

	// Marshal returns the bytes encoding of v.
	Marshal(v interface{}) ([]byte, error)

	// WithOrder changes the byte order for the new Writer.
	WithOrder(order binary.ByteOrder) Writer
}

// NewWriter returns a new writer that write to r with byte order.
// If debug set true, all write bytes and offsets will be displayed.
func NewWriter(w io.Writer, order binary.ByteOrder, debug bool) Writer {
	return &writer{
		w:     w,
		order: order,
		debug: debug,
	}
}

// NewWriterToBuf returns a new writer that writes to buffer with byte order.
// If debug set true, all write bytes and offsets will be displayed.
func NewWriterToBuf(buf []byte, order binary.ByteOrder, debug bool) Writer {
	return NewWriter(bytes.NewBuffer(buf), order, debug)
}

type writer struct {
	w     io.Writer
	order binary.ByteOrder

	debug bool
}

func (w *writer) Write(p []byte) (n int, err error) {
	return w.w.Write(p)
}

func (w *writer) WriteBool(b bool) error {
	panic("implement me")
}

func (w *writer) WriteUint8(v uint8) error {
	panic("implement me")
}

func (w *writer) WriteUint16(v uint16) error {
	panic("implement me")
}

func (w *writer) WriteUint32(v uint32) error {
	panic("implement me")
}

func (w *writer) WriteUint64(v uint64) error {
	panic("implement me")
}

func (w *writer) WriteInt8(v int8) error {
	panic("implement me")
}

func (w *writer) WriteInt16(v int16) error {
	panic("implement me")
}

func (w *writer) WriteInt32(v int32) error {
	panic("implement me")
}

func (w *writer) WriteInt64(v int64) error {
	panic("implement me")
}

func (w *writer) WriteFloat32(v float32) error {
	panic("implement me")
}

func (w *writer) WriteFloat64(v float64) error {
	panic("implement me")
}

func (w *writer) Marshal(v interface{}) ([]byte, error) {
	panic("implement me")
}

func (w *writer) WithOrder(order binary.ByteOrder) Writer {
	panic("implement me")
}
