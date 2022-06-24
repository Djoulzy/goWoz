package gowoz

import (
	"fmt"
	"os"
)

func InitWozFile(fileName string) (*WOZFileFormat, error) {
	file, err := os.Open(fileName)
	defer file.Close()

	if err != nil {
		return nil, err
	}

	tmp := WOZFileFormat{}
	tmp.init(file)

	return &tmp, err
}

func (W *WOZFileFormat) init(f *os.File) {
	var chunkHeader WOZChunkHeader
	var n int
	var err error

	W.fdesc = f

	W.Header.read(f)

	for {
		n, err = chunkHeader.read(f)
		if err != nil {
			fmt.Printf("Lecture header: %d\n", n)
			panic(err)
		}
		if n == 0 {
			fmt.Printf("End of file\n")
			break
		}

		switch chunkHeader.ID {
		case "INFO":
			W.INFO.read(f, chunkHeader)
		case "TMAP":
			W.TMAP.read(f, chunkHeader)
		case "TRKS":
			W.TRKS.read(f, chunkHeader)
		case "META":
			W.META.read(f, chunkHeader)
		default:
			f.Seek(int64(chunkHeader.Size), 1)
		}
	}

	W.physicalTrack = 0
	W.dataTrack = 0
	W.bitStreamPos = 0
}

func (W *WOZFileFormat) Dump(full bool) {
	W.Header.dump()
	W.INFO.dump()
	W.META.dump()
	if full {
		W.TMAP.dump()
		W.TRKS.dump()
	}
}

func (W *WOZFileFormat) DumpTrack(num byte) {
	trkData := W.TRKS.Data[num].RawBytes()
	for _, val := range trkData {
		fmt.Printf("%08b", val)
	}
}

func (W *WOZFileFormat) getNextBit() byte {
	trkData := W.TRKS.Data[W.dataTrack]
	res := trkData.BitAt(W.bitStreamPos)
	W.bitStreamPos++
	if uint32(W.bitStreamPos) > W.TRKS.Tracks[W.dataTrack].BitCount {
		W.bitStreamPos = 0
	}
	// fmt.Printf("%d", res)
	return res
}

func (W *WOZFileFormat) GetNextByte() byte {
	var result byte

	for W.getNextBit() == 0 {
	}
	result = 0x80 // the bit we just retrieved is the high bit
	for i := 6; i >= 0; i-- {
		result |= W.getNextBit() << i
	}
	return result
}

func (W *WOZFileFormat) GoToTrack(num float32) {
	newDataTrack := W.TMAP.Map[num]
	W.bitStreamPos = W.bitStreamPos * int(W.TRKS.Tracks[newDataTrack].BitCount) / int(W.TRKS.Tracks[W.dataTrack].BitCount)

	W.physicalTrack = num
	W.dataTrack = newDataTrack

	fmt.Printf("%d\n", W.dataTrack)
}

func (W *WOZFileFormat) Seek(offset float32) {
	W.physicalTrack += offset
	W.GoToTrack(W.physicalTrack)
}
