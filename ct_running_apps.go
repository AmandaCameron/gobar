package main

import (
	"image"
	"strings"

	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/ewmh"
	"github.com/BurntSushi/xgbutil/xgraphics"
	"github.com/BurntSushi/xgbutil/xprop"
)

type RunningSource struct {
	X *xgbutil.XUtil
}

type RunningCommand struct {
	X    *xgbutil.XUtil
	win  xproto.Window
	name string
}

func init() {
	Register(RunningSource{})
}

func (rs RunningSource) GetMatches(inp string, ct *CommandTray) []Command {
	cmds := make([]Command, 0)

	if ct.X == nil {
		return []Command{}
	}

	clis, err := ewmh.ClientListGet(ct.X)

	if err != nil {
		return cmds
	}

	for _, xwin := range clis {
		dt, err := ewmh.CurrentDesktopGet(ct.X)

		if err != nil {
			dt = 0xFFFFFFFF
		}

		wdt, err := ewmh.WmDesktopGet(ct.X, xwin)
		if err != nil {
			wdt = dt
		}

		if dt != wdt {
			continue
		}

		name, err := xprop.PropValStr(xprop.GetProperty(ct.X, xwin, "_NET_WM_NAME"))
		if err != nil {
			//print("Err1: ", err.Error(), "\n")
			name, err = xprop.PropValStr(xprop.GetProperty(ct.X, xwin, "WM_NAME"))
			if err != nil {
				//print("Err2: ", err.Error(), "\n")
				name = "Unnamed Window"
			}
		}

		if strings.Contains(strings.ToLower(name), inp) {
			cmds = append(cmds, RunningCommand{X: ct.X, win: xwin, name: name})
		}
	}

	return cmds
}

func (rc RunningCommand) GetIcon() image.Image {
	ico, err := xgraphics.FindIcon(rc.X, rc.win, 16, 16)
	if err != nil {
		return nil
	}

	return ico
}

func (rc RunningCommand) GetText() string {
	return rc.name
}

func (rc RunningCommand) Run() {
	ewmh.ActiveWindowReq(rc.X, rc.win)
}
