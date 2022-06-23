package main

import "github.com/Djoulzy/gowoz"

func main() {
	disk, err := gowoz.InitWozFile("I.woz")
	if err != nil {
		panic(err)
	}
	disk.Dump()
}
