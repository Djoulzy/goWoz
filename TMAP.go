package gowoz

import (
	"fmt"
	"os"
)

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
