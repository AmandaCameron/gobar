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

	cfg := &Config{}

	loadFont(dat, "Clock", &cfg.ClockFont)
	loadFont(dat, "Command", &cfg.CommandFont)

	tmp, err := dat.GetKey("Theme", "Size").Ints()[0]
	utils.FailMeMaybe(err)

	return cfg
}

func loadFont(dat *wini.Data, name string, target *FontInfo) error {
	tmp, err := dat.GetKey("Fonts/"+name, "Name")
	utils.FailMeMaybe(err)

	*target.Name = tmp.String()

	tmp, err = dat.GetKey("Fonts/"+name, "Size")
	utils.FailMeMaybe(err)

	*target.Size = tmp.Floats()[0]
}
