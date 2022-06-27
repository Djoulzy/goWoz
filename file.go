package gowoz

import (
	"fmt"
	"math/rand"
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
	switch W.Header.Format {
	case "WOZ1":
		W.Version = 1
	case "WOZ2":
		W.Version = 2
	default:
		panic("Unknown format")
	}

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
			W.TRKS.read(W.TMAP.Map, W.Version, f, chunkHeader)
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
	// Lecture d'un track vide
	// fmt.Printf("DataTrack: %v\n", W.dataTrack)
	if W.dataTrack == 0xFF {
		W.bitStreamPos++
		if W.bitStreamPos > 51200 {
			W.bitStreamPos = 0
		}
		return byte(rand.Intn(2))
	}

	trkData := W.TRKS.Data[W.dataTrack]
	res := trkData.BitAt(int(W.bitStreamPos))
	W.bitStreamPos++
	if W.bitStreamPos > W.TRKS.Tracks[W.dataTrack].BitCount {
		W.bitStreamPos = 0
	}
	// fmt.Printf("%d / \n", W.bitStreamPos)
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
	// fmt.Printf("T:%d Pos:%d -> %02X\n", W.dataTrack, W.bitStreamPos, result)
	return result
}

func (W *WOZFileFormat) GoToTrack(num float32) {
	newDataTrack, ok := W.TMAP.Map[num]
	if !ok {
		panic("bad track")
	}
	if newDataTrack == 0xFF {
		// fmt.Printf("Empty track %02.02f - actual pos: %d\n", num, W.bitStreamPos)
		W.bitStreamPos = W.bitStreamPos * (51200 / W.TRKS.Tracks[W.dataTrack].BitCount)
		W.physicalTrack = num
		W.dataTrack = 0xFF
		// fmt.Printf("new pos: %d\n", W.bitStreamPos)
	} else if newDataTrack == W.dataTrack {
		W.physicalTrack = num
		return
	} else {
		if W.dataTrack != 0xFF {
			W.bitStreamPos = W.bitStreamPos * (W.TRKS.Tracks[newDataTrack].BitCount / W.TRKS.Tracks[W.dataTrack].BitCount)
		}
		W.physicalTrack = num
		W.dataTrack = newDataTrack
		// fmt.Printf("Find Track %d - Pos: %d\n", W.dataTrack, W.bitStreamPos)
	}
	if W.bitStreamPos > 3 {
		W.bitStreamPos -= 4
	}
}

func (W *WOZFileFormat) Seek(offset float32) {
	var maxTrack float32
	destTrack := W.physicalTrack + offset
	fmt.Printf("Seek track %02.2f\n", destTrack)

	if W.Version >= 2 {
		maxTrack = 40
	} else {
		maxTrack = 35
	}
	if destTrack < 0 {
		destTrack = 0
	} else if destTrack > maxTrack {
		destTrack = maxTrack
	} else {
		W.GoToTrack(destTrack)
	}
}
