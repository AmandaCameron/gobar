package main

import (
	"image"
	"strings"

	"dbus/wicd"
)

type NetCommand struct {
	net *wicd.Network
}

type HiddenNetCommand struct {
	w *wicd.Wicd
}

type NetSource struct {
	w *wicd.Wicd
}

func (ns NetSource) GetMatches(inp string, ct *CommandTray) []Command {
	nets, err := ns.w.GetNetworks()

	if err != nil {
		print("Error: ", err)
		return []Command{}
	}
	cmds := make([]Command, 0)

	for _, net := range nets {
		if strings.Contains(strings.ToLower(net.Name), inp) {
			cmds = append(cmds, NetCommand{
				net: net,
			})
		}
	}

	if strings.Contains("connect to hidden wifi", inp) {
		cmds = append(cmds, HiddenNetCommand{
			w: ns.w,
		})
	}

	return cmds
}

func (ns NetSource) Connect(ct *CommandTray) {
	// Do Nothing.
}

// Known Networks

func (nc NetCommand) GetText() string {
	return "Connect to " + nc.net.Name
}

func (nc NetCommand) GetIcon() image.Image {
	str, err := nc.net.GetStrength()
	if err != nil {
		str = 0
	}
	return wifiStrengthImage(int32(str))
}

func (nc NetCommand) Run() {
	nc.net.Connect()
}

// Unknown Network

func (hnc HiddenNetCommand) GetText() string {
	return "Connect to a hidden network"
}

func (hnc HiddenNetCommand) GetIcon() image.Image {
	return nil
}

func (hnc HiddenNetCommand) Run() {
	print("To be implemented...")
}
