package gowoz

// Blaock size: 512 Bytes ( bits << 9)
// LSS:
//

import (
	"fmt"
	"math/rand"
)

var count int = 0
var wheel []byte = []byte{'-', '\\', '|', '/'}
var pickbit = []byte{128, 64, 32, 16, 8, 4, 2, 1}
var randomBits uint64 = 0b0100101001001001001001001001000101001010010010010010010010010001
var randBitsPos uint64 = 0x01

func (W *WOZFileFormat) makeDebug() {
	W.output = fmt.Sprintf("[%c]  %05.02f     %02d    %02d %5d", wheel[count], W.physicalTrack, W.dataTrack, W.revolution, W.bitStreamPos)
}

func (W *WOZFileFormat) IsWriteProtected() bool {
	return W.INFO.WriteProtected == 1
}

func (W *WOZFileFormat) GetMeta() map[string]string {
	return W.META.Metadata
}

func (W *WOZFileFormat) GetStreamPos() uint32 {
	return W.bitStreamPos
}

func (W *WOZFileFormat) GetRevolutionNumber() int {
	return W.revolution
}

func (W *WOZFileFormat) getNoise() byte {
	res := randomBits & randBitsPos
	randBitsPos <<= 1
	if res > 0 {
		return 1
	}
	return 0
}

func (W *WOZFileFormat) getNextWozBit() byte {
	var currentLength uint32
	var res byte

	if W.dataTrack == 0xFF {
		currentLength = 51200
		res = byte(rand.Intn(2))
	} else {
		currentLength = W.TRKS.Tracks[W.dataTrack].BitCount

		newPos := W.bitStreamPos % currentLength
		targetByte := newPos >> 3
		targetBit := newPos & 7

		res = (W.TRKS.Data[W.dataTrack][targetByte] & pickbit[targetBit]) >> (7 - targetBit)
	}

	W.bitStreamPos++
	if W.bitStreamPos > currentLength {
		W.bitStreamPos = 0
		W.revolution++
	}
	return res
}

func (W *WOZFileFormat) getNextBit() byte {
	// time.Sleep(4 * time.Microsecond)
	W.headWindow = W.headWindow << 1
	W.headWindow |= W.getNextWozBit()
	if (W.headWindow & 0x0f) != 0x00 {
		return (W.headWindow & 0x02) >> 1
	} else {
		return W.getNoise()
	}
}

// func (W *WOZFileFormat) LSSRead() byte {
// 	// time.Sleep(4 * time.Microsecond)
// 	W.getNextBit()
// 	return W.getNextBit()
// }

func (W *WOZFileFormat) GetNextByte() byte {
	var bit byte
	var result byte

	result = 0
	for bit = 0; bit == 0; bit = W.getNextBit() {
	}
	result = 0x80 // the bit we just retrieved is the high bit
	for i := 6; i >= 0; i-- {
		result |= W.getNextBit() << i
	}

	if debug {
		W.makeDebug()
	}
	count++
	if count >= len(wheel) {
		count = 0
	}
	// fmt.Printf("Trk: %05.02f = %02X\n", W.physicalTrack, result)
	return result
}

func (W *WOZFileFormat) GoToTrack(num float32) {
	var currentLength uint32

	newDataTrack, ok := W.TMAP.Map[num]
	if !ok {
		panic("bad track")
	}

	if newDataTrack == 0xFF {
		W.physicalTrack = num
		W.dataTrack = newDataTrack
		return
	}

	W.revolution = 0

	if newDataTrack == W.dataTrack {
		W.physicalTrack = num
		return
	}

	if W.dataTrack == 0xFF {
		currentLength = 51200
	} else {
		currentLength = W.TRKS.Tracks[W.dataTrack].BitCount
	}
	W.bitStreamPos = W.bitStreamPos * W.TRKS.Tracks[newDataTrack].BitCount / currentLength
	W.physicalTrack = num
	W.dataTrack = newDataTrack
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
	if debug {
		W.makeDebug()
	}
}

func (W *WOZFileFormat) GetCurrentTrack() float32 {
	return W.physicalTrack
}

func (W *WOZFileFormat) GetCurrentDataTrack() byte {
	return W.dataTrack
}

func (W *WOZFileFormat) GetStatus() string {
	return W.output
}
