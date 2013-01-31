package statbar

import (
	// "fmt"
	"image"

	"github.com/AmandaCameron/gobar/images"

	"github.com/AmandaCameron/go.networkmanager"
)

type SbNmWifi struct {
	Dev *nm.Device
}

func (icon *SbNmWifi) Attach(sb *StatusBar) {
	icon.Dev.PropChanged(func(_ *nm.Device) {
		sb.Draw()
	})
}

func (icon *SbNmWifi) Icon() image.Image {
	active, err := icon.Dev.GetActive()
	if err != nil {
		return images.WifiDC
	}

	if active == nil {
		return images.WifiDC
	}

	str, err := active.Strength()
	if err != nil {
		str = 0
	}

	return images.WifiStrengthImage(str)
}
