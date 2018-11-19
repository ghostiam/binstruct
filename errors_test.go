package binstruct

import "testing"

func TestIsEOF(t *testing.T) {
	var v struct {
		_ [2]byte
		I int64
	}
	err := UnmarshalBE([]byte{0x00, 0x01}, &v)
	if !IsEOF(err) {
		t.Error(err)
	}
}

func TestIsUnexpectedEOF(t *testing.T) {
	var v struct {
		I int64
	}
	err := UnmarshalBE([]byte{0x00, 0x01, 0x02}, &v)
	if !IsUnexpectedEOF(err) {
		t.Error(err)
	}
}
