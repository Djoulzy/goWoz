package main

import (
	"fmt"

	"github.com/Djoulzy/gowoz"
)

func main() {
	disk, err := gowoz.InitWozFile("Choplifter.woz")
	if err != nil {
		panic(err)
	}
	disk.Dump(true)
	// disk.DumpTrack(0)

	disk.GoToTrack(34)
	for x := 0; x < 10; x++ {
		fmt.Printf("%d\n", disk.GetNextByte())
	}

	disk.GoToTrack(0)
	for x := 0; x < 10; x++ {
		fmt.Printf("%d\n", disk.GetNextByte())
	}

	disk.GoToTrack(16)
	for x := 0; x < 10; x++ {
		fmt.Printf("%d\n", disk.GetNextByte())
	}
}
