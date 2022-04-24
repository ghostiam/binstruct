package binstruct

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math"
)

var (
	// ErrNegativeCount is returned when an attempt is made to read a negative number of bytes
	ErrNegativeCount = errors.New("binstruct: negative count")
)

// Reader is the interface that wraps the binstruct reader methods.
type Reader interface {
	io.ReadSeeker

	// Peek returns the next n bytes without advancing the reader.
	Peek(n int) ([]byte, error)

	// ReadBytes reads up to n bytes. It returns the number of bytes
	// read, bytes and any error encountered.
	ReadBytes(n int) (an int, b []byte, err error)
	// ReadAll reads until an error or EOF and returns the data it read.
	ReadAll() ([]byte, error)

	// ReadByte read and return one byte
	ReadByte() (byte, error)
	// ReadBool read one byte and return boolean value
	ReadBool() (bool, error)

	// ReadUint8 read one byte and return uint8 value
	ReadUint8() (uint8, error)
	// ReadUint16 read two bytes and return uint16 value
	ReadUint16() (uint16, error)
	// ReadUint32 read four bytes and return uint32 value
	ReadUint32() (uint32, error)
	// ReadUint64 read eight bytes and return uint64 value
	ReadUint64() (uint64, error)
	// ReadUintX read X bytes and return uint64 value
	ReadUintX(x int) (uint64, error)

	// ReadInt8 read one byte and return int8 value
	ReadInt8() (int8, error)
	// ReadInt16 read two bytes and return int16 value
	ReadInt16() (int16, error)
	// ReadInt32 read four bytes and return int32 value
	ReadInt32() (int32, error)
	// ReadInt64 read eight bytes and return int64 value
	ReadInt64() (int64, error)
	// ReadIntX read X bytes and return int64 value
	ReadIntX(x int) (int64, error)

	// ReadFloat32 read four bytes and return float32 value
	ReadFloat32() (float32, error)
	// ReadFloat64 read eight bytes and return float64 value
	ReadFloat64() (float64, error)

	// Unmarshal parses the binary data and stores the result
	// in the value pointed to by v.
	Unmarshal(v interface{}) error

	// WithOrder changes the byte order for the new Reader
	WithOrder(order binary.ByteOrder) Reader
}

// NewReader returns a new reader that reads from r with byte order.
// If debug set true, all read bytes and offsets will be displayed.
func NewReader(r io.ReadSeeker, order binary.ByteOrder, debug bool) Reader {
	return &reader{
		r:     r,
		order: order,
		debug: debug,
	}
}

// NewReaderFromBytes returns a new reader that reads from data with byte order.
// If debug set true, all read bytes and offsets will be displayed.
func NewReaderFromBytes(data []byte, order binary.ByteOrder, debug bool) Reader {
	return NewReader(bytes.NewReader(data), order, debug)
}

type reader struct {
	r     io.ReadSeeker
	order binary.ByteOrder

	debug bool
}

func (r *reader) ReadAll() ([]byte, error) {
	b, err := ioutil.ReadAll(r)

	if r.debug {
		fmt.Printf("ReadAll(): %s", hex.Dump(b))
	}

	return b, err
}

// If an EOF happens after reading some but not all the bytes, ReadBytes returns io.ErrUnexpectedEOF.
func (r *reader) ReadBytes(n int) (an int, b []byte, err error) {
	if n < 0 {
		return 0, nil, ErrNegativeCount
	}

	if n == 0 {
		return 0, []byte{}, nil
	}

	b = make([]byte, n)
	an, err = io.ReadFull(r, b)

	if r.debug {
		fmt.Printf("Read(want: %d|actual: %d): %s", n, an, hex.Dump(b))
	}

	if err != nil {
		return an, b, err
	}

	return an, b, nil
}

func (r *reader) ReadByte() (byte, error) {
	return r.ReadUint8()
}

func (r *reader) ReadBool() (bool, error) {
	b, err := r.ReadByte()
	return b != 0, err
}

func (r *reader) ReadUint8() (uint8, error) {
	_, b, err := r.ReadBytes(1)
	if err != nil {
		return 0, err
	}

	return b[0], nil
}

func (r *reader) ReadUint16() (uint16, error) {
	_, b, err := r.ReadBytes(2)
	if err != nil {
		return 0, err
	}

	return r.order.Uint16(b), nil
}

func (r *reader) ReadUint32() (uint32, error) {
	_, b, err := r.ReadBytes(4)
	if err != nil {
		return 0, err
	}

	return r.order.Uint32(b), nil
}

func (r *reader) ReadUint64() (uint64, error) {
	_, b, err := r.ReadBytes(8)
	if err != nil {
		return 0, err
	}

	return r.order.Uint64(b), nil
}

func (r *reader) ReadUintX(x int) (uint64, error) {
	var i uint64

	if x > 8 {
		return 0, errors.New("cannot read more than 8 bytes for custom length (u)int")
	}

	_, b, err := r.ReadBytes(x)
	if err != nil {
		return 0, err
	}

	switch r.order {
	case binary.BigEndian:
		for j := 0; j < x; j++ {
			i |= uint64(b[x-j-1]) << (8 * j)
		}

	case binary.LittleEndian:
		for j := 0; j < x; j++ {
			i |= uint64(b[j]) << (8 * j)
		}

	default:
		err = errors.New("cannot determine endianness for custom (u)int length read")
	}

	return i, err
}

func (r *reader) ReadInt8() (int8, error) {
	i, err := r.ReadUint8()
	return int8(i), err
}

func (r *reader) ReadInt16() (int16, error) {
	i, err := r.ReadUint16()
	return int16(i), err
}

func (r *reader) ReadInt32() (int32, error) {
	i, err := r.ReadUint32()
	return int32(i), err
}

func (r *reader) ReadInt64() (int64, error) {
	i, err := r.ReadUint64()
	return int64(i), err
}

func (r *reader) ReadIntX(x int) (int64, error) {
	u, err := r.ReadUintX(x)
	if err != nil {
		return 0, err
	}

	// Properly handle negatives by shifting fully left and then right
	u = u << (64 - 8*x)
	i := int64(u) >> (64 - 8*x)

	return i, nil
}

func (r *reader) ReadFloat32() (float32, error) {
	b, err := r.ReadUint32()
	if err != nil {
		return 0, err
	}

	float := math.Float32frombits(b)
	return float, nil
}

func (r *reader) ReadFloat64() (float64, error) {
	b, err := r.ReadUint64()
	if err != nil {
		return 0, err
	}

	float := math.Float64frombits(b)
	return float, nil
}

// io.Reader
func (r *reader) Read(p []byte) (n int, err error) {
	return r.r.Read(p)
}

// io.Seeker
func (r *reader) Seek(offset int64, whence int) (int64, error) {
	i, err := r.r.Seek(offset, whence)

	if r.debug {
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

func (r *reader) Peek(n int) ([]byte, error) {
	rn, b, err := r.ReadBytes(n)
	if err != nil {
		return nil, err
	}

	_, err = r.Seek(int64(-rn), io.SeekCurrent) // set offset back
	return b, err
}

func (r *reader) Unmarshal(v interface{}) error {
	u := &unmarshal{r}
	return u.Unmarshal(v)
}

func (r *reader) WithOrder(order binary.ByteOrder) Reader {
	return NewReader(r, order, r.debug)
}
