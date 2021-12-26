package main

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/davecgh/go-spew/spew"

	"github.com/ghostiam/binstruct"
)

// .ZIP File Format Specification: https://pkware.cachefly.net/webdocs/casestudies/APPNOTE.TXT
func main() {
	file, err := os.Open("sample.zip")
	if err != nil {
		log.Fatal(err)
	}

	var zip ZIP
	decoder := binstruct.NewDecoder(file, binary.LittleEndian)
	// decoder.SetDebug(true)
	err = decoder.Decode(&zip)
	if err != nil {
		log.Fatal(err)
	}

	spew.Dump(zip)
}

type ZIP struct {
	_                       byte                        `bin:"ParseZIPSections"` // Call helper method for scan
	LocalFileSections       []ZIPLocalFileSection       `bin:"-"`
	CentralDirEntrySections []ZIPCentralDirEntrySection `bin:"-"`
	EndOfCentralDirSection  ZIPEndOfCentralDirSection   `bin:"-"`
}

func (zip *ZIP) ParseZIPSections(r binstruct.Reader) error {
	for {
		// Find magic PK (0x50 0x4B)
		var magicPrevByte byte
		for {
			b, err := r.ReadByte()
			if errors.Is(err, io.EOF) {
				return nil
			}
			if err != nil {
				return fmt.Errorf("failed read magic PK: %w", err)
			}

			if magicPrevByte == 'P' && b == 'K' {
				magicPrevByte = 0x00
				break // exit from loop
			}

			magicPrevByte = b
		}

		// read section type
		_, sectionType, err := r.ReadBytes(2)
		if err != nil {
			return fmt.Errorf("failed read section type: %w", err)
		}

		switch {
		case bytes.Equal(sectionType, []byte{0x03, 0x04}):
			// parse ZIPLocalFileSection
			var localFileSection ZIPLocalFileSection
			err = r.Unmarshal(&localFileSection)
			if err != nil {
				return fmt.Errorf("failed Unmarshal ZIPLocalFileSection: %w", err)
			}

			zip.LocalFileSections = append(zip.LocalFileSections, localFileSection)

		case bytes.Equal(sectionType, []byte{0x01, 0x02}):
			// parse CentralDirEntry
			var centralDirEntrySections ZIPCentralDirEntrySection
			err = r.Unmarshal(&centralDirEntrySections)
			if err != nil {
				return fmt.Errorf("failed Unmarshal ZIPCentralDirEntrySection: %w", err)
			}

			zip.CentralDirEntrySections = append(zip.CentralDirEntrySections, centralDirEntrySections)

		case bytes.Equal(sectionType, []byte{0x05, 0x06}):
			// parse EndOfCentralDir
			var endOfCentralDirSections ZIPEndOfCentralDirSection
			err = r.Unmarshal(&endOfCentralDirSections)
			if err != nil {
				return fmt.Errorf("failed Unmarshal ZIPEndOfCentralDir: %w", err)
			}

			zip.EndOfCentralDirSection = endOfCentralDirSections

		default:
			log.Printf("unknown section type: %#x", sectionType)
		}
	}
}

type ZIPLocalFileSection struct {
	LocalFileHeader
	Body []byte `bin:"len:CompressedSize"`
}

type LocalFileHeader struct {
	Version           uint16
	Flags             [2]byte
	CompressionMethod uint16
	FileModTime       uint16
	FileModDate       uint16
	Crc32             [4]byte
	CompressedSize    uint32
	UncompressedSize  uint32
	FileNameLen       uint16
	ExtraLen          uint16
	FileName          string `bin:"len:FileNameLen"`
	Extra             []byte `bin:"len:ExtraLen"`
}

type ZIPCentralDirEntrySection struct {
	VersionMadeBy          int16
	VersionNeededToExtract int16
	Flags                  [2]byte
	CompressionMethod      int16
	LastModFileTime        int16
	LastModFileDate        int16
	Crc32                  [4]byte
	CompressedSize         int32
	UncompressedSize       int32
	FileNameLen            int16
	ExtraLen               int16
	CommentLen             int16
	DiskNumberStart        int16
	IntFileAttr            int16
	ExtFileAttr            int32
	LocalHeaderOffset      int32
	FileName               string `bin:"len:FileNameLen"`
	Extra                  []byte `bin:"len:ExtraLen"`
	Comment                string `bin:"len:CommentLen"`
}

type ZIPEndOfCentralDirSection struct {
	DiskOfEndOfCentralDir      int16
	DiskOfCentralDir           int16
	QtyCentralDirEntriesOnDisk int16
	QtyCentralDirEntriesTotal  int16
	CentralDirSize             int32
	CentralDirOffset           int32
	CommentLen                 int16
	Comment                    string `bin:"len:CommentLen"`
}
