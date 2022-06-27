package gowoz

import (
	"encoding/binary"
	"fmt"
	"os"
)

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
