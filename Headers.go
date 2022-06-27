package gowoz

import (
	"encoding/binary"
	"errors"
	"fmt"
	"os"
)

///////////////////////////////////////////
//             File Header               //
///////////////////////////////////////////

func (W *WOZHeader) read(f *os.File) {
	var tmp [1]byte
	var tmp4 [4]byte

	f.Read(tmp4[:])
	W.Format = fmt.Sprintf("%s", tmp4)
	f.Read(tmp[:])
	W.HighBits = tmp[0]
	f.Read(W.LFCRLF[:])
	f.Read(W.CRC[:])
}

func (W *WOZHeader) dump() {
	fmt.Printf("== Header\n")
	fmt.Printf("\tFormat: %s\n", W.Format)
	fmt.Printf("\tHigBits: %02X - CRLF: %03X\n", W.HighBits, W.LFCRLF)
	fmt.Printf("\tCRC: %04X\n\n", W.CRC)
}

///////////////////////////////////////////
//             Chunk Header              //
///////////////////////////////////////////

func (W *WOZChunkHeader) read(f *os.File) (int, error) {
	var tmp [4]byte
	n, err := f.Read(tmp[:])
	W.ID = fmt.Sprintf("%s", tmp)
	if n == 0 {
		return 0, nil
	}
	if err != nil {
		return n, err
	}
	n, err = f.Read(tmp[:])
	if n == 0 {
		return -1, errors.New("Malformed file")
	}
	if err != nil {
		return n, err
	}
	W.Size = binary.LittleEndian.Uint32(tmp[:])
	return n, nil
}
