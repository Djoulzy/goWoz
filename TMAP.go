package gowoz

import (
	"fmt"
	"os"
	"sort"
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
	fmt.Printf("Map len: %d\n", len(W.Map))
	W.Sorted = make([]float64, 0)
	for k := range W.Map {
		W.Sorted = append(W.Sorted, float64(k))
	}
	sort.Float64s(W.Sorted)
	fmt.Printf("Sorted len: %d\n", len(W.Sorted))
}

// func (W *WOZTMapChunk) dump() {
// 	var cpt float32

// 	fmt.Printf("== TMap\n")
// 	for cpt = 0; cpt <= 40; cpt += 0.25 {
// 		val, ok := W.Map[cpt]
// 		if ok {
// 			fmt.Printf("Physical Track %0.2f: %d\n", cpt, val)
// 		}
// 	}
// }

func (W *WOZTMapChunk) dump() {
	var tmp float64
	var partieEntiere int
	var partieDecimale string
	var val byte
	var ok bool

	fmt.Printf("== TMap\n")
	for index, trk := range W.Sorted {
		if val, ok = W.Map[float32(trk)]; !ok {
			panic("TMAP: Bad size")
		}
		partieEntiere = int(trk)
		tmp = trk - float64(partieEntiere)
		if tmp != 0 {
			partieDecimale = "+." + fmt.Sprintf("%.2f", trk-float64(partieEntiere))[2:]
		} else {
			partieDecimale = ""
		}

		if val == 0xFF {
			fmt.Printf("%c[%sm%c[%sm", 27, "0;31", 27, "47")
		}

		fmt.Printf("TMAP[0x%02X] track 0x%02X %4s: TRKS track index 0x%02X", index, partieEntiere, partieDecimale, val)

		if val == 0xFF {
			fmt.Printf("%c[%sm%c[%sm", 27, "0;37", 27, "40")
		}

		fmt.Printf("\n")
		// cpt++
	}
}
