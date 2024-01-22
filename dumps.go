package gowoz

import "fmt"

///////////////////////////////////////////////////
//                     DUMPERS                   //
///////////////////////////////////////////////////

func (W *WOZFileFormat) Dump(full bool) {
	W.Header.dump()
	W.INFO.dump()
	W.META.dump()
	W.TMAP.dump()
	if full {
		W.TRKS.dump(W.TMAP.Map)
		W.DumpTracksRaw()
	}
}

func (W *WOZFileFormat) DumpTrack(track float32) {
	var val byte
	var count int = 0

	W.GoToTrack(track)
	for W.GetRevolutionNumber() == 0 {
		val = W.GetNextByte()
		fmt.Printf("%02X ", val)
		count++
		if count%51 == 0 {
			fmt.Printf(" %d \n", W.GetStreamPos())
			count = 0
		}
	}

}

func (W *WOZFileFormat) DumpTracksRaw() {
	for index := range W.TRKS.Data {
		fmt.Printf("TRK index %02X: %08x bytes; %08x bits\n", index, W.TRKS.Tracks[index].ByteCount, W.TRKS.Tracks[index].BitCount)
	}
}
