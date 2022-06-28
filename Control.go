package gowoz

import (
	"fmt"
	"math/rand"
)

var count int = 0
var wheel []byte = []byte{'-', '\\', '|', '/'}

func (W *WOZFileFormat) getNextBit() byte {
	// Lecture d'un track vide
	// fmt.Printf("DataTrack: %v\n", W.dataTrack)
	if W.TMAP.Map[W.physicalTrack] == 0xFF {
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
	if W.dataTrack == 0xFF {
		return 0
	}

	for W.getNextBit() == 0 {
	}
	result = 0x80 // the bit we just retrieved is the high bit
	for i := 6; i >= 0; i-- {
		result |= W.getNextBit() << i
	}

	// fmt.Printf("%c T:%02.02f (%d) Pos:%d\r", wheel[count], W.physicalTrack, W.dataTrack, W.bitStreamPos)
	// count++
	// if count >= len(wheel) {
	// 	count = 0
	// }
	return result
}

func (W *WOZFileFormat) GoToTrack(num float32) {
	// var oldDataTrackLength uint32

	newDataTrack, ok := W.TMAP.Map[num]
	if !ok {
		panic("bad track")
	}

	// if W.dataTrack == 0xFF {
	// 	oldDataTrackLength = 51200
	// } else {
	// 	oldDataTrackLength = W.TRKS.Tracks[W.dataTrack].BitCount
	// }

	if newDataTrack == W.dataTrack {
		W.physicalTrack = num
		return
	}

	// if newDataTrack == 0xFF {
	// 	// fmt.Printf("Empty track %02.02f - actual pos: %d\n", num, W.bitStreamPos)
	// 	W.bitStreamPos = W.bitStreamPos * (51200 / oldDataTrackLength)
	// } else {
	// 	W.bitStreamPos = W.bitStreamPos * (W.TRKS.Tracks[newDataTrack].BitCount / oldDataTrackLength)
	// }
	W.bitStreamPos = 0

	W.physicalTrack = num
	W.dataTrack = newDataTrack
	if W.bitStreamPos > 3 {
		W.bitStreamPos -= 4
	}

	fmt.Printf("Move to T:%02.02f (%d) at pos %d\n", W.physicalTrack, W.dataTrack, W.bitStreamPos)
}

func (W *WOZFileFormat) Seek(offset float32) {
	var maxTrack float32
	destTrack := W.physicalTrack + offset
	fmt.Printf("Seek Track offset %.02f -> %d\n", offset, W.TMAP.Map[destTrack])

	if W.Version >= 2 {
		maxTrack = 40
	} else {
		maxTrack = 35
	}

	if destTrack < 0 {
		destTrack = 0
	} else if destTrack > maxTrack {
		destTrack = maxTrack
	}
	W.GoToTrack(destTrack)
}
