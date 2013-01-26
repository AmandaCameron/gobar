package statbar

import (
	// "fmt"
	"image"

	"github.com/AmandaCameron/gobar/images"

	"github.com/AmandaCameron/gobar/utils/dbus/nm"
)

type SbNmWifi struct {
	dev *nm.Device
}

func (icon *SbNmWifi) Attach(sb *StatusBar) {
	icon.dev.PropChanged(func(_ *nm.Device) {
		sb.Draw()
	})
}

func (icon *SbNmWifi) Icon() image.Image {
	active, err := icon.dev.GetActive()
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
