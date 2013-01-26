package main

import (
	"github.com/BurntSushi/wingo/wini"

	"github.com/AmandaCameron/gobar/utils"
)

type Config struct {
	ClockFont   FontInfo
	CommandFont FontInfo
	BarSize     int
}

type FontInfo struct {
	Name string
	Size float64
}

func loadConfig(fileName string) *Config {
	dat, err := wini.Parse(fileName)
	utils.FailMeMaybe(err)

	cfg := &Config{}

	loadFont(dat, "Clock", &cfg.ClockFont)
	loadFont(dat, "Command", &cfg.CommandFont)

	tmp, err := dat.GetKey("Theme", "Size").Ints()
	utils.FailMeMaybe(err)

	cfg.BarSize = tmp[0]

	return cfg
}

func loadFont(dat *wini.Data, name string, target *FontInfo) {
	tmp := dat.GetKey("Fonts/"+name, "Name")
	//utils.FailMeMaybe(err)

	name = tmp.String()

	tmp = dat.GetKey("Fonts/"+name, "Size")
	//utils.FailMeMaybe(err)

	tmp2, err := tmp.Floats()
	utils.FailMeMaybe(err)

	size := tmp2[0]

	*target = FontInfo{
		name, size,
	}
}
