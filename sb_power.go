package main

import (
	// "fmt"
	"image"
	"image/draw"

	"dbus/upower"
)

type SbPower struct {
	dev *upower.Device
}

func (icon *SbPower) Attach(sb *StatusBar) {
	icon.dev.Connect(func(_ *upower.Device) {
		// fmt.Printf("[SB] (Power) State Changed.\n")
		sb.Draw()
	})
}

func (icon *SbPower) Icon() image.Image {
	img := image.NewRGBA(image.Rect(0, 0, 34, 16))

	state, err := icon.dev.State()
	if err != nil {
		return img
	}

	charge, err := icon.dev.Charge()
	if err != nil {
		return img
	}

	if state == upower.Charging || state == upower.Full {
		draw.Draw(img, image.Rect(0, 0, 16, 16), charging_img, image.Point{0, 0}, draw.Over)
	}

	var bimg image.Image

	switch {
	case charge > 90:
		bimg = battery_5_img
	case charge > 50:
		bimg = battery_4_img
	case charge > 25:
		bimg = battery_3_img
	case charge > 10:
		bimg = battery_2_img
	case charge <= 10:
		bimg = battery_1_img
	}

	draw.Draw(img, image.Rect(18, 0, 34, 16), bimg, image.Point{0, 0}, draw.Over)

	return img
}
