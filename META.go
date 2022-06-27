package gowoz

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

///////////////////////////////////////////
//                  META                 //
///////////////////////////////////////////

func (W *WOZChunkMeta) read(f *os.File, header WOZChunkHeader) {
	var tmp []byte
	W.Header = header

	tmp = make([]byte, header.Size)
	f.Read(tmp)
	buff := fmt.Sprintf("%s", tmp)

	if len(buff) > 0 {
		W.Metadata = make(map[string]string)
		r := csv.NewReader(strings.NewReader(buff))
		r.Comma = '\t'
		r.FieldsPerRecord = 2
		for {
			record, err := r.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatal(err)
			}

			W.Metadata[record[0]] = record[1]
		}
	}
}

func (W *WOZChunkMeta) dump() {
	fmt.Printf("== Meta\n")
	for label, txt := range W.Metadata {
		fmt.Printf("\t%s: %s\n", label, txt)
	}
}
