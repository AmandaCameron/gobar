package main

import (
	"image"
	"strings"

	"github.com/BurntSushi/xgbutil/ewmh"
	"github.com/BurntSushi/xgbutil/xprop"
	"launchpad.net/~jamesh/go-dbus/trunk"
)

type AppMenuSource struct {
	sessConn *dbus.Connection
	menu     *GtkMenu
	app      *GtkActions
}

type AppMenuCommand struct {
	AppCommand *GtkActions
	opt        map[string]dbus.Variant
}

func (ams *AppMenuSource) GetMatches(inp string, ct *CommandTray) []Command {
	if ams.app == nil || ams.menu == nil {
		return []Command{}
	}

	menus, err := ams.menu.Start([]uint32{0})
	if err != nil {
		return nil
	}

	ret := make([]Command, 0)

	for _, group := range menus[1:] {
		// if group.MenuId != 1 {
		// 	ret = append(ret, &AppMenuSeperator{})
		// }

		for _, opt := range group.Options {
			cmd := &AppMenuCommand{
				ams.app,
				opt,
			}

			if strings.Contains(strings.ToLower(cmd.GetText()), inp) || strings.Contains("menu", inp) {
				ret = append(ret, cmd)
			}
		}
	}

	return ret
}

func (ams *AppMenuSource) Open(ct *CommandTray) bool {
	ams.app = nil
	ams.menu = nil

	active, err := ewmh.ActiveWindowGet(ct.X)
	if err != nil {
		return false
	}

	uniqName, err := xprop.PropValStr(xprop.GetProperty(ct.X, active, "_GTK_UNIQUE_BUS_NAME"))
	if err != nil {
		return false
	}

	pathMenu, err := xprop.PropValStr(xprop.GetProperty(ct.X, active, "_GTK_APP_MENU_OBJECT_PATH"))
	if err != nil {
		return false
	}

	pathApp, err := xprop.PropValStr(xprop.GetProperty(ct.X, active, "_GTK_APPLICATION_OBJECT_PATH"))
	if err != nil {
		return false
	}

	// Done parsing props! Yay!

	ams.app = &GtkActions{
		ams.sessConn.Object(uniqName, dbus.ObjectPath(pathApp)),
	}

	ams.menu = &GtkMenu{
		ams.sessConn.Object(uniqName, dbus.ObjectPath(pathMenu)),
	}

	return true
}

func (ams *AppMenuSource) Close(ct *CommandTray) {
	ams.menu.End([]uint32{0})
}

func (amc *AppMenuCommand) GetText() string {
	tmp, ok := amc.opt["label"]
	if !ok {
		return "[Menu] ** Error **"
	}

	label, ok := tmp.Value.(string)
	if !ok {
		return "[Menu] ** Error **"
	}

	return "[Menu] " + label
}

func (amc *AppMenuCommand) GetIcon() image.Image {
	return nil
}

func (amc *AppMenuCommand) Run() {
	tmp, ok := amc.opt["action"]
	if !ok {
		return
	}

	act, ok := tmp.Value.(string)
	if !ok {
		return
	}

	if act[0:4] == "app." {
		act = act[4:]
	}

	amc.AppCommand.Activate(act, nil, nil)
}
