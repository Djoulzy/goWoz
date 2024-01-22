package gowoz

import (
	"fmt"
	"math/rand"
	"os"
)

var debug bool
var percentOfOne float64 = 30

func InitContainer(fileName string, debugMode bool) (*WOZFileFormat, error) {
	file, err := os.Open(fileName)
	defer file.Close()

	if err != nil {
		return nil, err
	}

	debug = debugMode
	tmp := WOZFileFormat{}
	tmp.init(file)

	return &tmp, err
}

func (W *WOZFileFormat) genRandomBits() {
	var percent float64 = percentOfOne / 100
	var nbOne int = int(percent*256) + 1
	var i int = 0
	var place int

	fmt.Printf("%d\n", nbOne)
	for i < nbOne {
		place = rand.Intn(256)
		if W.randomBits[place] == 1 {
			continue
		} else {
			W.randomBits[place] = 1
			i++
		}
	}

	W.randomBits = [256]byte{
		1, 0, 0, 0, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 1, 0, 0, 0, 1, 0, 0, 1, 0, 0, 0, 0, 0,
		1, 0, 0, 0, 0, 1, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		1, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 1, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0, 1,
		0, 1, 1, 0, 1, 0, 1, 0, 0, 1, 0, 0, 0, 0, 0, 1, 1, 0, 0, 0, 1, 0, 0, 1, 0, 1, 0, 0, 0, 0, 0, 0,
		0, 1, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0, 1, 1, 0, 1, 0, 1, 1, 0,
		0, 0, 0, 0, 0, 1, 0, 1, 0, 1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1, 0, 1, 1, 0, 0, 1, 0, 0, 0, 0, 1, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0,
		1, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 1, 0, 1, 0, 0, 1, 0, 0, 1, 1, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0,
	}
}

func (W *WOZFileFormat) diplayRandomBits() {
	for index, b := range W.randomBits {
		fmt.Printf("%b,", b)
		if (index+1)%32 == 0 {
			fmt.Printf("\n")
		}
	}
	fmt.Printf("\n")
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

	W.physicalTrack = 10
	W.dataTrack = W.TMAP.Map[10]
	W.bitStreamPos = 0
	W.revolution = 0
	W.headWindow = 0
	W.randBitsPos = 0

	W.genRandomBits()
	W.diplayRandomBits()
}
