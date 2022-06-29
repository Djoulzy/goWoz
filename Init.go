package gowoz

import (
	"fmt"
	"os"
)

func InitContainer(fileName string) (*WOZFileFormat, error) {
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
	W.revolution = 0
}

func (W *WOZFileFormat) Dump(full bool) {
	W.Header.dump()
	W.INFO.dump()
	W.META.dump()
	if full {
		W.TRKS.dump(W.TMAP.Map)
	}
}