package main

import "github.com/Djoulzy/gowoz"

func main() {
	disk, err := gowoz.InitWozFile("Choplifter.woz")
	if err != nil {
		panic(err)
	}
	disk.Dump()
}
