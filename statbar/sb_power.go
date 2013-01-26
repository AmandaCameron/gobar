package statbar

import (
	// "fmt"
	"image"
	"image/draw"

	"github.com/AmandaCameron/gobar/images"

	"github.com/AmandaCameron/gobar/utils/dbus/upower"
)

type SbPower struct {
	Dev *upower.Device
}

func (icon *SbPower) Attach(sb *StatusBar) {
	icon.Dev.Connect(func(_ *upower.Device) {
		// fmt.Printf("[SB] (Power) State Changed.\n")
		sb.Draw()
	})
}

func (icon *SbPower) Icon() image.Image {
	img := image.NewRGBA(image.Rect(0, 0, 34, 16))

	state, err := icon.Dev.State()
	if err != nil {
		return img
	}

	charge, err := icon.Dev.Charge()
	if err != nil {
		return img
	}

	if state == upower.Charging || state == upower.Full {
		draw.Draw(img, image.Rect(0, 0, 16, 16), images.Charging, image.Point{0, 0}, draw.Over)
	}

	var bimg image.Image

	switch {
	case charge > 90:
		bimg = images.Battery_5
	case charge > 50:
		bimg = images.Battery_4
	case charge > 25:
		bimg = images.Battery_3
	case charge > 10:
		bimg = images.Battery_2
	case charge <= 10:
		bimg = images.Battery_1
	}

	draw.Draw(img, image.Rect(18, 0, 34, 16), bimg, image.Point{0, 0}, draw.Over)

	return img
}
