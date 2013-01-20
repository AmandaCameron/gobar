package main

import (
	"image"
	"image/draw"

	"dbus/upower"
)

type SbPower struct {
	up *upower.UPower
}

func (icon *SbPower) Attach(sb *StatusBar) {
	icon.up.Connect(func(_ *upower.UPower) {
		sb.Draw()
	})
}

func (icon *SbPower) Icon() image.Image {
	img := image.NewRGBA(image.Rect(0, 0, 34, 16))

	batt, err := icon.up.GetBattery()
	if err != nil {
		return nil
	}

	if batt.State == upower.Charging || batt.State == upower.Full {
		draw.Draw(img, image.Rect(0, 0, 16, 16), charging_img, image.Point{0, 0}, draw.Over)
	}

	var bimg image.Image

	switch {
	case batt.Charge > 90:
		bimg = battery_5_img
	case batt.Charge > 50:
		bimg = battery_4_img
	case batt.Charge > 25:
		bimg = battery_3_img
	case batt.Charge > 10:
		bimg = battery_2_img
	case batt.Charge <= 10:
		bimg = battery_1_img
	}

	draw.Draw(img, image.Rect(18, 0, 34, 16), bimg, image.Point{0, 0}, draw.Over)

	return img
}
