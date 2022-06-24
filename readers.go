package gowoz

import (
	"encoding/binary"
	"errors"
	"fmt"
	"os"

	"github.com/tunabay/go-bitarray"
)

///////////////////////////////////////////
//             File Header               //
///////////////////////////////////////////

func (W *WOZHeader) read(f *os.File) {
	var tmp [1]byte
	f.Read(W.Format[:])
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

///////////////////////////////////////////
//                  INFO                 //
///////////////////////////////////////////

func (W *WOZInfoChunk) read(f *os.File, header WOZChunkHeader) {
	var tmp1 [1]byte
	var tmp2 [2]byte

	W.Header = header

	f.Read(tmp1[:])
	W.Version = tmp1[0]
	f.Read(tmp1[:])
	W.DiskType = tmp1[0]
	f.Read(tmp1[:])
	W.WriteProtected = tmp1[0]
	f.Read(tmp1[:])
	W.Synchronized = tmp1[0]
	f.Read(tmp1[:])
	W.Cleaned = tmp1[0]
	f.Read(W.Creator[:])

	if W.Version > 1 {
		f.Read(tmp1[:])
		W.DiskSides = tmp1[0]
		f.Read(tmp1[:])
		W.BootSectorFormat = tmp1[0]
		f.Read(tmp1[:])
		W.OptimalBitTiming = tmp1[0]
		f.Read(tmp2[:])
		W.CompatibleHardware = binary.LittleEndian.Uint16(tmp2[:])
		f.Read(tmp2[:])
		W.RequiredRAM = binary.LittleEndian.Uint16(tmp2[:])
		f.Read(tmp2[:])
		W.LargestTrack = binary.LittleEndian.Uint16(tmp2[:])
	}

	if W.Version > 2 {
		f.Read(tmp2[:])
		W.FLUXBlock = binary.LittleEndian.Uint16(tmp2[:])
		f.Read(tmp2[:])
		W.LargestFluxTrack = binary.LittleEndian.Uint16(tmp2[:])
	}

	f.Seek(80, 0)
}

func (W *WOZInfoChunk) dump() {
	fmt.Printf("== Infos\n")
	fmt.Printf("\tVersion: %d\n", W.Version)
	fmt.Printf("\tDiskType: %s\n", DiskType[W.DiskType])
	fmt.Printf("\tWriteProtected: %d\n", W.WriteProtected)
	fmt.Printf("\tSynchronized: %d\n", W.Synchronized)
	fmt.Printf("\tCleaned: %d\n", W.Cleaned)
	fmt.Printf("\tCreator: %s\n", W.Creator)
	if W.Version >= 2 {
		fmt.Printf("\tDiskSides: %d\n", W.DiskSides)
		fmt.Printf("\tBootSectorFormat: %s\n", BootSector[W.BootSectorFormat])
		fmt.Printf("\tOptimalBitTiming: %d\n", W.OptimalBitTiming)
		fmt.Printf("\tCompatibleHardware: %d\n", W.CompatibleHardware)
		fmt.Printf("\tRequiredRAM: %dK\n", W.RequiredRAM)
		fmt.Printf("\tLargestTrack: %d blocks (%d bytes)\n", W.LargestTrack, W.LargestTrack*512)
	}
	if W.Version >= 3 {
		fmt.Printf("\tFLUXBlock: %d\n", W.FLUXBlock)
		fmt.Printf("\tLargestFluxTrack: %d\n", W.LargestFluxTrack)
	}
}

///////////////////////////////////////////
//                  TMAP                 //
///////////////////////////////////////////

func (W *WOZTMapChunk) read(f *os.File, header WOZChunkHeader) {
	var tmp []byte
	var cpt float32
	W.Header = header

	tmp = make([]byte, 160)
	f.Read(tmp)

	W.Map = make(map[float32]byte)
	cpt = 0
	for _, val := range tmp {
		W.Map[cpt] = val
		cpt += 0.25
	}
}

func (W *WOZTMapChunk) dump() {
	var cpt float32

	fmt.Printf("== TMap\n")
	for cpt = 0; cpt <= 40; cpt += 0.25 {
		val, ok := W.Map[cpt]
		if ok {
			if val != 0xFF {
				fmt.Printf("Physical Track %0.2f: %d\n", cpt, val)
			}
		}
	}
}

///////////////////////////////////////////
//                  META                 //
///////////////////////////////////////////

func (W *WOZChunkMeta) read(f *os.File, header WOZChunkHeader) {
	var tmp []byte
	W.Header = header

	tmp = make([]byte, header.Size)
	f.Read(tmp)
	W.Metadata = fmt.Sprintf("%s", tmp)
}

func (W *WOZChunkMeta) dump() {
	fmt.Printf("== Meta\n")
	fmt.Printf("\t%s\n", W.Metadata)
}

///////////////////////////////////////////
//                  TRKS                 //
///////////////////////////////////////////

func (W *WOZTrackDesc) read(f *os.File) {
	var tmp2 [2]byte
	var tmp4 [4]byte

	f.Read(tmp2[:])
	W.StartBlock = binary.LittleEndian.Uint16(tmp2[:])
	f.Read(tmp2[:])
	W.BlockCount = binary.LittleEndian.Uint16(tmp2[:])
	f.Read(tmp4[:])
	W.BitCount = binary.LittleEndian.Uint32(tmp4[:])
}

func (W *WOZTRKSChunk) read(f *os.File, header WOZChunkHeader) {
	var dataStart uint32
	var blockBuff []byte
	// var countBit uint32
	// var mask byte
	// var bitLoaded bool

	W.Header = header

	// Read tracks infos
	for t := 0; t < 160; t++ {
		W.Tracks[t].read(f)
	}

	// Read tracks data
	for t := 0; t < 160; t++ {
		// if t > 0 {
		// 	panic(1)
		// }
		if W.Tracks[t].BlockCount == 0 {
			continue
		}
		dataStart = uint32(W.Tracks[t].StartBlock) << 9
		f.Seek(int64(dataStart), 0)
		blockBuff = make([]byte, W.Tracks[t].BlockCount<<9)
		f.Read(blockBuff)

		W.Data[t] = bitarray.NewBufferFromByteSlice(blockBuff)
		// fmt.Printf("blocks: %d - Bits: %d\n", W.Tracks[t].BlockCount, W.Tracks[t].BitCount)
		// countBit = 0
		// bitLoaded = false
		// for _, pack := range blockBuff {
		// 	// fmt.Printf("%08b", pack)
		// 	for i := 0; i < 8; i++ {
		// 		mask = 0b10000000 >> i
		// 		if pack&mask == mask {
		// 			W.Data[t].
		// 		}
		// 		// if pack&mask == mask {
		// 		// 	fmt.Printf("1")
		// 		// } else {
		// 		// 	fmt.Printf("0")
		// 		// }
		// 		countBit++
		// 		if countBit == W.Tracks[t].BitCount {
		// 			bitLoaded = true
		// 			break
		// 		}
		// 	}
		// 	if bitLoaded {
		// 		break
		// 	}
		// }
	}

	f.Seek(int64(256+header.Size), 0)
}

func (W *WOZTRKSChunk) dump() {
	fmt.Printf("== TRKS\n")
	for t := 0; t < 160; t++ {
		if W.Tracks[t].BlockCount == 0 {
			continue
		}
		fmt.Printf("Track %02d : %d blocks (%d bits / %d bytes) starts at block %d (byte %d)- len: %d\n", t, W.Tracks[t].BlockCount, W.Tracks[t].BitCount, W.Tracks[t].BitCount/8, W.Tracks[t].StartBlock, uint32(W.Tracks[t].StartBlock)<<9, W.Data[t].Len())
	}
}
