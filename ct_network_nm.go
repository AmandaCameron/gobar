package main

import (
	"image"
	"image/draw"
	"strings"

	"dbus/nm"
)

type NmSource struct {
	// cli *nm.Client
	dev *nm.Device
}

type NmHiddenCommand struct {
	dev *nm.Device
}

type NmCommand struct {
	ap *nm.AccessPoint
}

func (ns NmSource) GetMatches(inp string, ct *CommandTray) []Command {
	aps, err := ns.dev.GetAccessPoints()
	if err != nil {
		return []Command{}
	}

	var cmds []Command

	for _, ap := range aps {
		name, err := ap.Name()
		if err != nil {
			continue
		}

		if inp == "wifi" || strings.Contains(strings.ToLower(name), inp) {
			cmds = append(cmds, &NmCommand{
				// cli: ns.cli,
				ap: ap,
			})
		}
	}

	if strings.Contains("connect to hidden wifi", inp) {
		cmds = append(cmds, &NmHiddenCommand{
			//cli: ns.cli,
			dev: ns.dev,
		})
	}

	return cmds
}

func (ns NmSource) Open(ct *CommandTray) bool {
	return true
}

func (ns NmSource) Close(ct *CommandTray) {
	// Do Nothing.
}

// Hidden Network

func (nc *NmHiddenCommand) GetText() string {
	return "Connect to hidden network"
}

func (nc *NmHiddenCommand) GetIcon() image.Image {
	return wifi_4_img
}

func (nc *NmHiddenCommand) Run() {
	// TODO.
}

// Regular Networks

func (nc *NmCommand) GetText() string {
	name, err := nc.ap.Name()
	if err != nil {
		name = "** Error **"
	}

	return "Connect to " + name
}

func (nc *NmCommand) GetIcon() image.Image {
	str, err := nc.ap.Strength()
	if err != nil {
		return wifi_dc_img
	}

	flags, err := nc.ap.Flags()
	if err != nil {
		flags = 0
	}

	if (flags & 1) == 1 {
		img := image.NewRGBA(image.Rect(0, 0, 16, 16))

		draw.Draw(img, image.Rect(0, 0, 16, 16), wifiStrengthImage(int32(str)), image.Point{0, 0}, draw.Over)
		draw.Draw(img, image.Rect(0, 0, 16, 16), wifi_enc_image, image.Point{0, 0}, draw.Over)

		return img
	}

	return wifiStrengthImage(int32(str))
}

func (nc *NmCommand) Run() {
	// TODO.
}
