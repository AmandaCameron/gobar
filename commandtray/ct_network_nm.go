package commandtray

import (
	"image"
	"image/draw"
	"strings"

	"github.com/AmandaCameron/gobar/images"

	"github.com/AmandaCameron/go.networkmanager"
)

type NmSource struct {
	// cli *nm.Client
	Dev *nm.Device
}

type NmHiddenCommand struct {
	Dev *nm.Device
}

type NmCommand struct {
	ap *nm.AccessPoint
}

func (ns NmSource) GetMatches(inp string, ct *CommandTray) []Command {
	aps, err := ns.Dev.GetAccessPoints()
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
			Dev: ns.Dev,
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
	return images.Wifi_4
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
		return images.WifiDC
	}

	flags, err := nc.ap.Flags()
	if err != nil {
		flags = 0
	}

	if (flags & 1) == 1 {
		img := image.NewRGBA(image.Rect(0, 0, 16, 16))

		draw.Draw(img, image.Rect(0, 0, 16, 16), images.WifiStrengthImage(str), image.Point{0, 0}, draw.Over)
		draw.Draw(img, image.Rect(0, 0, 16, 16), images.WifiEnc, image.Point{0, 0}, draw.Over)

		return img
	}

	return images.WifiStrengthImage(str)
}

func (nc *NmCommand) Run() {
	// TODO.
}
