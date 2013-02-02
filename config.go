package main

import (
	"github.com/BurntSushi/wingo/wini"

	"github.com/AmandaCameron/gobar/utils"
)

type Config struct {
	CommandAccel string
	BarSize      int
	BarWidth     int
	Position     string
	ClockFormat  string
	Clock        Section
	StatusBar    Section
	Command      Section
	Tracker      Section
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

	// Main thing.

	loadInt(dat, "Main", "Size", &cfg.BarSize)
	loadInt(dat, "Main", "Width", &cfg.BarWidth)
	loadString(dat, "Main", "Position", "", &cfg.Position)

	// Clock.

	loadSection(dat, "Clock", true, &cfg.Clock)
	loadString(dat, "Clock", "Format", "2006-01-02 15:04:05", &cfg.ClockFormat)

	// Tracker

	loadSection(dat, "App Tracker", false, &cfg.Tracker)

	// Command Tray

	loadSection(dat, "CommandTray", true, &cfg.Command)
	loadString(dat, "CommandTray", "Accel", "", &cfg.CommandAccel)

	// Status Bar

	loadSection(dat, "StatusBar", false, &cfg.StatusBar)

	return cfg
}

func loadSection(dat *wini.Data, name string, fontMeMaybe bool, target *Section) {
	sec := Section{}

	loadInt(dat, name, "Width", &sec.Width)
	loadInt(dat, name, "Position", &sec.Position)

	if fontMeMaybe {
		loadFont(dat, name, &sec.Font)
	}

	*target = sec
}

func loadFont(dat *wini.Data, name string, target *FontInfo) {
	font := FontInfo{}

	loadString(dat, "Fonts/"+name, "Name", "", &font.Name)
	loadFloat(dat, "Fonts/"+name, "Size", &font.Size)

	*target = font
}

func loadString(dat *wini.Data, name, key, def string, target *string) {
	k := dat.GetKey(name, key)
	if k == nil {
		if def == "" {
			utils.Fail("Missing key " + key + " in section " + name)
		} else {
			*target = def
			return
		}
	}

	tmp := k.Strings()

	if len(tmp) > 0 {
		*target = tmp[0]
	} else {
		if def == "" {
			utils.Fail("Missing key " + key + " in section " + name)
		} else {
			*target = def
			return
		}
	}
}

func loadInt(dat *wini.Data, name, key string, target *int) {
	k := dat.GetKey(name, key)
	if k == nil {
		utils.Fail("Missing key " + key + " in section " + name)
	}

	tmp, err := k.Ints()
	utils.FailMeMaybe(err)

	if len(tmp) > 0 {
		*target = tmp[0]
	} else {
		utils.Fail("Missing key " + key + " in section " + name)
	}
}

func loadFloat(dat *wini.Data, name, key string, target *float64) {
	k := dat.GetKey(name, key)
	if k == nil {
		utils.Fail("Missing key " + key + " in section " + name)
	}

	tmp, err := k.Floats()
	utils.FailMeMaybe(err)

	if len(tmp) > 0 {
		*target = tmp[0]
	} else {
		utils.Fail("Missing key " + key + " in section " + name)
	}
}
