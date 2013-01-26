package utils

import (
	"os"

	"code.google.com/p/freetype-go/freetype/truetype"

	"github.com/BurntSushi/xgbutil/xgraphics"
)

func OpenFont(fileName string) *truetype.Font {
	f, err := os.Open(fileName)

	FailMeMaybe(err)

	defer f.Close()

	fnt, err := xgraphics.ParseFont(f)

	FailMeMaybe(err)

	return fnt
}

func FailMeMaybe(err error) {
	if err != nil {
		panic(err)
	}
}
