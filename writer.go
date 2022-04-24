package binstruct

import (
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"math"
)

type Writer interface {
	io.WriteSeeker

	io.ByteWriter
	io.StringWriter

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
	// WriteUintX write X bytes uint64.
	WriteUintX(v uint64, x int) error

	// WriteInt8 write one byte int8.
	WriteInt8(v int8) error
	// WriteInt16 write two bytes int16.
	WriteInt16(v int16) error
	// WriteInt32 write four bytes int32.
	WriteInt32(v int32) error
	// WriteInt64 write eight bytes int64.
	WriteInt64(v int64) error
	// WriteIntX write X bytes int64.
	WriteIntX(v int64, x int) error

	// WriteFloat32 write four bytes float32.
	WriteFloat32(v float32) error
	// WriteFloat64 write eight bytes float64.
	WriteFloat64(v float64) error

	// Len returns the number of bytes written.
	Len() int

	// Marshal returns the bytes encoding of v.
	Marshal(v interface{}) ([]byte, error)

	// WithOrder changes the byte order for the new Writer.
	WithOrder(order binary.ByteOrder) Writer
}

// NewWriter returns a new writer that write to w with byte order.
// If debug set true, all write bytes and offsets will be displayed.
func NewWriter(w io.Writer, order binary.ByteOrder, debug bool) Writer {
	return &writer{
		w:     w,
		order: order,
		debug: debug,
	}
}

type writer struct {
	w     io.Writer
	order binary.ByteOrder

	debug   bool
	written int
}

func (w *writer) Write(p []byte) (n int, err error) {
	n, err = w.w.Write(p)
	w.written += n

	if w.debug {
		fmt.Printf("Write(want: %d|actual: %d): %s", len(p), n, hex.Dump(p))
	}

	return n, err
}

func (w *writer) WriteByte(c byte) (err error) {
	if sw, ok := w.w.(io.ByteWriter); ok {
		err = sw.WriteByte(c)
	} else {
		_, err = w.w.Write([]byte{c})
	}
	w.written += 1 // always 1 byte

	if w.debug {
		fmt.Printf("WriteByte: %s", hex.Dump([]byte{c}))
	}

	return err
}

func (w *writer) WriteString(s string) (n int, err error) {
	if sw, ok := w.w.(io.StringWriter); ok {
		n, err = sw.WriteString(s)
	} else {
		n, err = w.w.Write([]byte(s))
	}
	w.written += n

	if w.debug {
		fmt.Printf("WriteString(want: %d|actual: %d): %s\n", len(s), n, s)
	}

	return n, err
}

func (w *writer) WriteBool(v bool) error {
	b := byte(0)
	if v {
		b = 1
	}

	_, err := w.Write([]byte{b})
	return err
}

func (w *writer) WriteUint8(v uint8) error {
	_, err := w.Write([]byte{v})
	return err
}

func (w *writer) WriteUint16(v uint16) error {
	b := make([]byte, 2)
	w.order.PutUint16(b, v)
	_, err := w.Write(b)
	return err
}

func (w *writer) WriteUint32(v uint32) error {
	b := make([]byte, 4)
	w.order.PutUint32(b, v)
	_, err := w.Write(b)
	return err
}

func (w *writer) WriteUint64(v uint64) error {
	b := make([]byte, 8)
	w.order.PutUint64(b, v)
	_, err := w.Write(b)
	return err
}

func (w *writer) WriteUintX(v uint64, x int) error {
	if x <= 0 {
		return errors.New("cannot write less than 1 bytes for custom length (u)int")
	}

	if x > 8 {
		return errors.New("cannot write more than 8 bytes for custom length (u)int")
	}

	b := make([]byte, x)
	for i := x; i > 0; i-- {
		shift := (i - 1) * 8
		b[x-i] = byte(v >> shift)
	}

	_, err := w.Write(b)
	return err
}

func (w *writer) WriteInt8(v int8) error {
	return w.WriteUint8(uint8(v))
}

func (w *writer) WriteInt16(v int16) error {
	return w.WriteUint16(uint16(v))
}

func (w *writer) WriteInt32(v int32) error {
	return w.WriteUint32(uint32(v))
}

func (w *writer) WriteInt64(v int64) error {
	return w.WriteUint64(uint64(v))
}

func (w *writer) WriteIntX(v int64, x int) error {
	return w.WriteUintX(uint64(v), x)
}

func (w *writer) WriteFloat32(v float32) error {
	return w.WriteUint32(math.Float32bits(v))
}

func (w *writer) WriteFloat64(v float64) error {
	return w.WriteUint64(math.Float64bits(v))
}

func (w *writer) Len() int {
	return w.written
}

func (w *writer) Marshal(v interface{}) ([]byte, error) {
	// TODO implement me
	panic("implement me")
}

func (w *writer) WithOrder(order binary.ByteOrder) Writer {
	return NewWriter(w, order, w.debug)
}

// io.Seeker
func (w *writer) Seek(offset int64, whence int) (int64, error) {
	ws, ok := w.w.(io.Seeker)
	if !ok {
		return 0, errors.New("writer not implemented io.Seeker")
	}

	i, err := ws.Seek(offset, whence)

	if w.debug {
		whenceStr := "invalid"
		switch whence {
		case io.SeekStart:
			whenceStr = "SeekStart"
		case io.SeekCurrent:
			whenceStr = "SeekCurrent"
		case io.SeekEnd:
			whenceStr = "SeekEnd"
		}

		fmt.Printf("Seek(%d, %s) CurPos:%d\n", offset, whenceStr, i)
	}

	return i, err
}
