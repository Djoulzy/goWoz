package gowoz

import (
	"encoding/binary"
	"fmt"
	"os"

	"github.com/tunabay/go-bitarray"
)

///////////////////////////////////////////
//                  TRKS                 //
///////////////////////////////////////////

func (W *WOZTrackDesc) read(version int, f *os.File) {
	var tmp [1]byte
	var tmp2 [2]byte
	var tmp4 [4]byte

	W.Version = version

	if version >= 2 {
		f.Read(tmp2[:])
		W.StartBlock = binary.LittleEndian.Uint16(tmp2[:])
		f.Read(tmp2[:])
		W.BlockCount = binary.LittleEndian.Uint16(tmp2[:])
		f.Read(tmp4[:])
		W.BitCount = binary.LittleEndian.Uint32(tmp4[:])
	} else {
		f.Read(tmp2[:])
		W.BytesUsed = binary.LittleEndian.Uint16(tmp2[:])
		f.Read(tmp2[:])
		W.BitCount = uint32(binary.LittleEndian.Uint16(tmp2[:]))
		f.Read(tmp4[:])
		W.SplicePoint = binary.LittleEndian.Uint16(tmp2[:])
		f.Read(tmp[:])
		W.SpliceNibble = tmp[0]
		f.Read(tmp[:])
		W.SpliceBitCount = tmp[0]
		f.Read(tmp2[:])
		W.Reserved = binary.LittleEndian.Uint16(tmp2[:])
	}
}

func (W *WOZTRKSChunk) read(MAP map[float32]byte, version int, f *os.File, header WOZChunkHeader) {
	var dataStart uint32
	var blockBuff []byte

	W.Header = header
	W.Version = version

	if version >= 2 {
		// Read tracks infos v2
		for t := 0; t < 160; t++ {
			W.Tracks[t].read(version, f)
		}
		// Read tracks data
		for _, track := range MAP {
			if track == 0xFF {
				continue
			} else {
				dataStart = uint32(W.Tracks[track].StartBlock) << 9
				f.Seek(int64(dataStart), 0)
				blockBuff = make([]byte, W.Tracks[track].BlockCount<<9)
				f.Read(blockBuff)

				W.Data[track] = bitarray.NewBufferFromByteSlice(blockBuff)
			}
		}

	} else {
		// Read tracks data
		for _, track := range MAP {
			if track == 0xFF {
				continue
			} else {
				dataStart = 256 + (uint32(track) * 6656)
				f.Seek(int64(dataStart), 0)
				blockBuff = make([]byte, 6646)
				f.Read(blockBuff)

				W.Data[track] = bitarray.NewBufferFromByteSlice(blockBuff)
				W.Tracks[track].read(version, f)
			}
		}
	}
	f.Seek(int64(256+header.Size), 0)
}

func (W *WOZTRKSChunk) dump(MAP map[float32]byte) {
	var cpt float32

	fmt.Printf("== TRKS\n")
	for cpt = 0; cpt <= 40; cpt += 0.25 {
		val, ok := MAP[cpt]
		if ok {
			if val == 0xFF {
				continue
			}
			if W.Version >= 2 {
				fmt.Printf("Physical Track %0.2f = Track %02d : %d blocks (%d bits / %d bytes) starts at block %d (byte %d) - len: %d\n", cpt, val, W.Tracks[val].BlockCount, W.Tracks[val].BitCount, W.Tracks[val].BitCount/8, W.Tracks[val].StartBlock, uint32(W.Tracks[val].StartBlock)<<9, W.Data[val].Len())
			} else {
				fmt.Printf("Physical Track %0.2f = Track %02d : %d bits / %d bytes (used: %d) - len: %d\n", cpt, val, W.Tracks[val].BitCount, W.Tracks[val].BitCount/8, W.Tracks[val].BytesUsed, W.Data[val].Len())
			}
		}
	}
}
