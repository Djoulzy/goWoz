package gowoz

import (
	"encoding/binary"
	"fmt"
	"os"
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
		W.ByteCount = (uint32(W.BlockCount)) << 9
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
				W.Data[track] = make([]byte, int((W.Tracks[track].BitCount)>>3))
				f.Read(W.Data[track])
			}
		}
		// for _, tmp := range W.Data[0] {
		// 	fmt.Printf("%02X ", tmp)
		// }
	} else {
		// Read tracks data
		for _, track := range MAP {
			if track == 0xFF {
				continue
			} else {
				dataStart = 256 + (uint32(track) * 6656)
				f.Seek(int64(dataStart), 0)
				W.Data[track] = make([]byte, 6646)
				f.Read(W.Data[track])

				W.Tracks[track].read(version, f)
			}
		}
	}
	f.Seek(int64(256+header.Size), 0)
}

// func (W *WOZTRKSChunk) dump(MAP map[float32]byte) {
// 	var cpt float32

// 	fmt.Printf("== TRKS\n")
// 	fmt.Printf(" Ph.Trk | Dat.Trk | Blks |  Bits | Start |  Len \n")
// 	fmt.Printf("--------+---------+------+-------+-------+------\n")
// 	for cpt = 0; cpt <= 40; cpt += 0.25 {
// 		val, ok := MAP[cpt]
// 		if ok {
// 			if val == 0xFF {
// 				// fmt.Printf("Physical Track %0.2f = %02X\n", cpt, MAP[cpt])
// 				continue
// 			}
// 			if W.Version >= 2 {
// 				fmt.Printf("  %05.02f |    %02d   |  %02d  | %05d |  %03d  | %04d\n", cpt, val, W.Tracks[val].BlockCount, W.Tracks[val].BitCount, W.Tracks[val].StartBlock, len(W.Data[val]))
// 			} else {
// 				fmt.Printf("Physical Track %0.2f = Track %02d : %d bits / %d bytes (used: %d) - len: %d\n", cpt, val, W.Tracks[val].BitCount, W.Tracks[val].BitCount/8, W.Tracks[val].BytesUsed, len(W.Data[val]))
// 			}
// 		}
// 	}
// }

func (W *WOZTRKSChunk) reverseByte(val byte) string {
	return fmt.Sprintf("%d%d%d%d%d%d%d%d",
		(val&0b10000000)>>7,
		(val&0b01000000)>>6,
		(val&0b00100000)>>5,
		(val&0b00010000)>>4,
		(val&0b00001000)>>3,
		(val&0b00000100)>>2,
		(val&0b00000010)>>1,
		val&0b00000001,
	)
}

func (W *WOZTRKSChunk) dump(MAP map[float32]byte) {
	fmt.Printf("== TRKS\n")
	for index, trk := range W.Tracks {
		if trk.BlockCount > 0 {
			fmt.Printf("TRK index %02X: start byte in BITS %08x; ends %08x; %08x bytes; %08x bits (", index, 1536+(trk.StartBlock)<<9, 1536+(trk.StartBlock)<<9+trk.BlockCount<<9, trk.BlockCount<<9, trk.BitCount)
			for i := 0; i < 8; i++ {
				fmt.Printf("%s", W.reverseByte(W.Data[index][i]))
			}
			fmt.Printf("...)\n")
		}
	}
}

// #define BYTE_TO_BINARY_PATTERN "%c%c%c%c%c%c%c%c"
// #define BYTE_TO_BINARY(byte)  \
//   (byte & 0b10000000 ? '1' : '0'), \
//   (byte & 0b01000000 ? '1' : '0'), \
//   (byte & 0b00100000 ? '1' : '0'), \
//   (byte & 0b00010000 ? '1' : '0'), \
//   (byte & 0b00001000 ? '1' : '0'), \
//   (byte & 0b00000100 ? '1' : '0'), \
//   (byte & 0b00000010 ? '1' : '0'), \
//   (byte & 0b00000001 ? '1' : '0')
