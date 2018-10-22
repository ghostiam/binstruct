package binstruct

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math"

	"github.com/davecgh/go-spew/spew"
)

var (
	ErrNegativeCount = errors.New("binstruct: negative count")
)

type Reader interface {
	ReadBytes(n int) (an int, b []byte, err error)
	ReadAll() ([]byte, error)

	ReadByte() (byte, error)
	ReadBool() (bool, error)

	ReadUint8() (uint8, error)
	ReadUint16() (uint16, error)
	ReadUint32() (uint32, error)
	ReadUint64() (uint64, error)

	ReadInt8() (int8, error)
	ReadInt16() (int16, error)
	ReadInt32() (int32, error)
	ReadInt64() (int64, error)

	ReadFloat32() (float32, error)
	ReadFloat64() (float64, error)
}

type Seeker interface {
	io.Seeker
}

type Peeker interface {
	Peek(n int) ([]byte, error)
}

type ReadSeekPeeker interface {
	Reader
	Seeker
	Peeker
	Unmarshaler
}

func NewReader(r io.ReadSeeker, order binary.ByteOrder, debug bool) ReadSeekPeeker {
	return &readSeekPeeker{
		r:     r,
		order: order,
		debug: debug,
	}
}

type readSeekPeeker struct {
	r     io.ReadSeeker
	order binary.ByteOrder

	debug bool
}

func (r *readSeekPeeker) Unmarshal(v interface{}) error {
	u := &unmarshal{r}
	return u.Unmarshal(v)
}

func (r *readSeekPeeker) ReadAll() ([]byte, error) {
	b, err := ioutil.ReadAll(r.r)

	if r.debug {
		fmt.Printf("ReadAll(): %s", spew.Sdump(b))
	}

	return b, err
}

// If an EOF happens after reading some but not all the bytes, ReadBytes returns io.ErrUnexpectedEOF.
func (r *readSeekPeeker) ReadBytes(n int) (an int, b []byte, err error) {
	if n < 0 {
		return 0, nil, ErrNegativeCount
	}

	if n == 0 {
		return 0, []byte{}, nil
	}

	b = make([]byte, n)
	an, err = io.ReadFull(r.r, b)

	if r.debug {
		fmt.Printf("Read(want: %d|actual: %d): %s", n, an, spew.Sdump(b))
	}

	if err != nil {
		return an, b, err
	}

	return an, b, nil
}

func (r *readSeekPeeker) ReadByte() (byte, error) {
	return r.ReadUint8()
}

func (r *readSeekPeeker) ReadBool() (bool, error) {
	b, err := r.ReadByte()
	return b != 0, err
}

func (r *readSeekPeeker) ReadUint8() (uint8, error) {
	_, b, err := r.ReadBytes(1)
	if err != nil {
		return 0, err
	}

	return b[0], nil
}

func (r *readSeekPeeker) ReadUint16() (uint16, error) {
	_, b, err := r.ReadBytes(2)
	if err != nil {
		return 0, err
	}

	return r.order.Uint16(b), nil
}

func (r *readSeekPeeker) ReadUint32() (uint32, error) {
	_, b, err := r.ReadBytes(4)
	if err != nil {
		return 0, err
	}

	return r.order.Uint32(b), nil
}

func (r *readSeekPeeker) ReadUint64() (uint64, error) {
	_, b, err := r.ReadBytes(8)
	if err != nil {
		return 0, err
	}

	return r.order.Uint64(b), nil
}

func (r *readSeekPeeker) ReadInt8() (int8, error) {
	i, err := r.ReadUint8()
	return int8(i), err
}

func (r *readSeekPeeker) ReadInt16() (int16, error) {
	i, err := r.ReadUint16()
	return int16(i), err
}

func (r *readSeekPeeker) ReadInt32() (int32, error) {
	i, err := r.ReadUint32()
	return int32(i), err
}

func (r *readSeekPeeker) ReadInt64() (int64, error) {
	i, err := r.ReadUint64()
	return int64(i), err
}

func (r *readSeekPeeker) ReadFloat32() (float32, error) {
	b, err := r.ReadUint32()
	if err != nil {
		return 0, err
	}

	float := math.Float32frombits(b)
	return float, nil
}

func (r *readSeekPeeker) ReadFloat64() (float64, error) {
	b, err := r.ReadUint64()
	if err != nil {
		return 0, err
	}

	float := math.Float64frombits(b)
	return float, nil
}

func (r *readSeekPeeker) Seek(offset int64, whence int) (int64, error) {
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

func (r *readSeekPeeker) Peek(n int) ([]byte, error) {
	rn, b, err := r.ReadBytes(n)
	r.Seek(int64(-rn), io.SeekCurrent) // set offset back
	return b, err
}
