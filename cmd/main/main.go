package main

import (
	"github.com/Djoulzy/gowoz"
)

func main() {
	disk, err := gowoz.InitWozFile("anti-m.woz")
	if err != nil {
		panic(err)
	}
	disk.Dump(true)
	// disk.DumpTrack(0)

	// disk.GoToTrack(1)
	// for x := 0; x < 25000; x++ {
	// 	disk.GetNextByte()
	// }

	// disk.Seek(0.5)
	// disk.Seek(0.5)
	// disk.Seek(0.5)
	// disk.Seek(0.5)
	// disk.Seek(0.5)
	// disk.Seek(0.5)

	// for x := 0; x < 10; x++ {
	// 	fmt.Printf("%d\n", disk.GetNextByte())
	// }

	// disk.GoToTrack(16)
	// for x := 0; x < 10; x++ {
	// 	fmt.Printf("%d\n", disk.GetNextByte())
	// }
}
