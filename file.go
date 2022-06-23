package gowoz

import "os"

func InitWozFile(fileName string) (*WOZFileFormat, error) {
	file, err := os.Open(fileName)
	defer file.Close()

	if err != nil {
		return nil, err
	}

	tmp := WOZFileFormat{}
	tmp.init(file)

	return &tmp, err
}
