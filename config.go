package main

import (
	"github.com/BurntSushi/wingo/wini"

	"github.com/AmandaCameron/gobar/utils"
)

type Config struct {
	CommandCombo string
	BarSize      int
	BarWidth     int
	Position     string
	Clock        Section
	StatusBar    Section
	Command      Section
}

type Section struct {
	Width    int
	Position int
	Font     FontInfo
}

type FontInfo struct {
	Name string
	Size float64
}

func loadConfig(fileName string) *Config {
	dat, err := wini.Parse(fileName)
	utils.FailMeMaybe(err)

	cfg := &Config{}

	//loadFont(dat, "Clock", &cfg.ClockFont)
	//loadFont(dat, "CommandTray", &cfg.CommandFont)

	// Main thing.

	tmp, err := dat.GetKey("Main", "Size").Ints()
	utils.FailMeMaybe(err)

	cfg.BarSize = tmp[0]

	tmp, err = dat.GetKey("Main", "Width").Ints()
	utils.FailMeMaybe(err)
	cfg.BarWidth = tmp[0]

	key := dat.GetKey("Main", "Position")
	if key == nil {
		utils.Fail("Must have a [Main] Position entry.")
	}

	cfg.Position = key.Strings()[0]

	// Clock.

	// tmp, err = dat.GetKey("Clock", "Width").Ints()
	// utils.FailMeMaybe(err)

	// cfg.ClockWidth = tmp[0]

	// tmp, err = dat.GetKey("Clock", "Position").Ints()
	// utils.FailMeMaybe(err)

	// cfg.ClockPos = tmp[0]

	loadSection(dat, "Clock", true, &cfg.Clock)

	// Command Tray

	loadSection(dat, "CommandTray", true, &cfg.Command)

	key = dat.GetKey("CommandTray", "Accel")
	if key == nil {
		utils.Fail("Must have a [CommandTray] Accel entry.")
	}
	cfg.CommandCombo = key.Strings()[0]

	// Status Bar

	loadSection(dat, "StatusBar", false, &cfg.StatusBar)

	return cfg
}

func loadSection(dat *wini.Data, name string, fontMeMaybe bool, target *Section) {
	sec := Section{}

	key := dat.GetKey(name, "Width")
	if key == nil {
		utils.Fail("Missing " + name + " entry -- Width")
	}

	tmp, err := key.Ints()
	utils.FailMeMaybe(err)

	sec.Width = tmp[0]

	key = dat.GetKey(name, "Position")
	if key == nil {
		utils.Fail("Missing " + name + " entry -- Position")
	}

	tmp, err = key.Ints()
	utils.FailMeMaybe(err)

	sec.Position = tmp[0]

	if fontMeMaybe {
		loadFont(dat, name, &sec.Font)
	}

	*target = sec
}

func loadFont(dat *wini.Data, name string, target *FontInfo) {
	key := dat.GetKey("Fonts/"+name, "Name")
	//utils.FailMeMaybe(err)
	if key == nil {
		utils.Fail("Fonts/" + name + " is missing field: Name")
	}

	tmp3 := key.Strings()
	if len(tmp3) == 0 {
		utils.Fail("Fonts/" + name + " is missing field: Name")
	}

	fntName := tmp3[0]

	key = dat.GetKey("Fonts/"+name, "Size")
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
