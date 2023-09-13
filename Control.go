package gowoz

import (
	"fmt"
	"math/rand"
)

var count int = 0
var wheel []byte = []byte{'-', '\\', '|', '/'}
var pickbit = []byte{128, 64, 32, 16, 8, 4, 2, 1}

func (W *WOZFileFormat) IsWriteProtected() bool {
	return W.INFO.WriteProtected == 1
}

func (W *WOZFileFormat) GetMeta() map[string]string {
	return W.META.Metadata
}

func (W *WOZFileFormat) getNextWozBit() byte {
	// Lecture d'un track vide
	// fmt.Printf("DataTrack: %v\n", W.dataTrack)

	// if W.TMAP.Map[W.physicalTrack] == 0xFF {
	// 	W.bitStreamPos++
	// 	if W.bitStreamPos > 51200 {
	// 		W.bitStreamPos = 0
	// 	}
	// 	return byte(rand.Intn(2))
	// }

	W.bitStreamPos = W.bitStreamPos % W.TRKS.Tracks[W.dataTrack].BitCount
	targetByte := W.bitStreamPos >> 3
	targetBit := W.bitStreamPos & 7

	res := (W.TRKS.Data[W.dataTrack][targetByte] & pickbit[targetBit]) >> (7 - targetBit)

	W.bitStreamPos++
	if W.bitStreamPos > W.TRKS.Tracks[W.dataTrack].BitCount {
		W.bitStreamPos = 0
		W.revolution++
	}
	return res
}

func (W *WOZFileFormat) getNextBit() byte {
	W.headWindow = W.headWindow << 1
	W.headWindow |= W.getNextWozBit()
	if (W.headWindow & 0x0f) != 0x00 {
		return (W.headWindow & 0x02) >> 1
	} else {
		return byte(rand.Intn(2))
	}
}

func (W *WOZFileFormat) GetNextByte() byte {
	var bit, result byte

	result = 0
	for bit = 0; bit == 0; bit = W.getNextBit() {
	}
	result = 0x80 // the bit we just retrieved is the high bit
	for i := 6; i >= 0; i-- {
		result |= W.getNextBit() << i
	}

	if debug {
		W.output = fmt.Sprintf("[%c]  %05.02f     %02d    %02d %5d", wheel[count], W.physicalTrack, W.dataTrack, W.revolution, W.bitStreamPos)
	}
	count++
	if count >= len(wheel) {
		count = 0
	}
	// fmt.Printf("Trk: %05.02f = %02X\n", W.physicalTrack, result)
	return result
}

func (W *WOZFileFormat) GoToTrack(num float32) {
	newDataTrack, ok := W.TMAP.Map[num]
	if !ok {
		panic("bad track")
	}

	W.revolution = 0

	if newDataTrack == W.dataTrack {
		W.physicalTrack = num
		return
	}

	W.physicalTrack = num
	W.dataTrack = newDataTrack
	// fmt.Printf("Move to T:%02.02f (%d) at pos %d\n", W.physicalTrack, W.dataTrack, W.bitStreamPos)
}

func (W *WOZFileFormat) Seek(offset float32) {
	var maxTrack float32
	destTrack := W.physicalTrack + offset
	// fmt.Printf("Seek Track offset %.02f -> %d\n", offset, W.TMAP.Map[destTrack])

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

func (W *WOZFileFormat) DumpTrack(track float32) {
	var val byte

	W.GoToTrack(track)
	W.bitStreamPos = 0
	for i := 1; i <= int(W.TRKS.Tracks[W.dataTrack].BlockCount*512); i++ {
		val = W.GetNextByte()
		fmt.Printf("%02X ", val)
		if i%32 == 0 {
			fmt.Printf("\n")
		}
	}
}

func (W *WOZFileFormat) DumpTrackRaw(track float32) {
	W.GoToTrack(track)
	W.bitStreamPos = 0
	// for i := 1; i <= int(W.TRKS.Tracks[W.dataTrack].BitCount); i++ {
	// 	val = W.getNextBit()
	// 	fmt.Printf("%1b", val)
	// 	if i%160 == 0 {
	// 		fmt.Printf("\n")
	// 	}
	// }
	for index, i := range W.TRKS.Data[W.dataTrack] {
		fmt.Printf("%02X ", i)
		if (index+1)%32 == 0 {
			fmt.Printf("\n")
		}
	}
}

func (W *WOZFileFormat) GetCurrentTrack() float32 {
	return W.physicalTrack
}

func (W *WOZFileFormat) GetStatus() string {
	return W.output
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
