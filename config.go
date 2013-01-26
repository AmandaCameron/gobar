package main

import (
	"github.com/BurntSushi/wingo/wini"

	"github.com/AmandaCameron/gobar/utils"
)

type Config struct {
	ClockFont    FontInfo
	CommandFont  FontInfo
	CommandCombo string
	BarSize      int
	Position     string
}

type FontInfo struct {
	Name string
	Size float64
}

func loadConfig(fileName string) *Config {
	dat, err := wini.Parse(fileName)
	utils.FailMeMaybe(err)

	cfg := &Config{}

	loadFont(dat, "clock", &cfg.ClockFont)
	loadFont(dat, "command", &cfg.CommandFont)

	tmp, err := dat.GetKey("theme", "Size").Ints()
	utils.FailMeMaybe(err)

	cfg.BarSize = tmp[0]

	key := dat.GetKey("commandtray", "Accel")
	if key == nil {
		utils.Fail("Must have a [CommandTray] Accel entry.")
	}
	cfg.CommandCombo = key.Strings()[0]

	key = dat.GetKey("theme", "Position")
	if key == nil {
		utils.Fail("Must have a [Theme] Position entry.")
	}

	cfg.Position = key.Strings()[0]

	return cfg
}

func loadFont(dat *wini.Data, name string, target *FontInfo) {
	key := dat.GetKey("fonts/"+name, "Name")
	//utils.FailMeMaybe(err)
	if key == nil {
		utils.Fail("Fonts/" + name + " is missing field: Name")
	}

	tmp3 := key.Strings()
	if len(tmp3) == 0 {
		utils.Fail("Fonts/" + name + " is missing field: Name")
	}

	fntName := tmp3[0]

	key = dat.GetKey("fonts/"+name, "Size")
	//utils.FailMeMaybe(err)

	if key == nil {
		utils.Fail("Fonts/" + name + " is missing field: Size")
	}

	tmp2, err := key.Floats()
	utils.FailMeMaybe(err)

	if len(tmp2) == 0 {
		utils.Fail("Fonts/" + name + " is missing field: Size")
	}

	size := tmp2[0]

	*target = FontInfo{
		fntName, size,
	}
}
