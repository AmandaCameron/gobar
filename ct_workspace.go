package main

import (
	"image"
	"strings"

	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/ewmh"
)

type WorkspaceSource struct {
	X *xgbutil.XUtil
}

type WorkspaceCommand struct {
	x    *xgbutil.XUtil
	id   uint
	name string
}

func init() {
	Register(WorkspaceSource{})
}

func (ws WorkspaceSource) GetMatches(inp string, ct *CommandTray) []Command {
	desks, err := ewmh.DesktopNamesGet(ct.X)

	if err != nil {
		return []Command{}
	}

	cmds := make([]Command, 0, len(desks))

	for i, name := range desks {
		if strings.Contains(strings.ToLower(name), inp) {
			cmds = append(cmds, WorkspaceCommand{
				x:    ct.X,
				id:   uint(i),
				name: name,
			})
		}
	}

	return cmds
}

func (wc WorkspaceCommand) GetText() string {
	return "Switch to " + wc.name
}

func (wc WorkspaceCommand) GetIcon() image.Image {
	return nil
}

func (wc WorkspaceCommand) Run() {
	ewmh.CurrentDesktopSet(wc.x, wc.id)
}
