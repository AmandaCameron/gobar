package main

import (
	"fmt"
	"image"
	"image/draw"

	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/xgraphics"

	"dbus/upower"
	"dbus/wicd"
)

var (
	status_bg = xgraphics.BGRA{R: 64, G: 64, B: 64, A: 255}
)

type StatusBar struct {
	img     *xgraphics.Image
	bar_img *xgraphics.Image

	up *upower.UPower
	w  *wicd.Wicd
}

func NewStatusBar(X *xgbutil.XUtil, up *upower.UPower, w *wicd.Wicd) *StatusBar {
	sb := &StatusBar{
		up:  up,
		w:   w,
		img: xgraphics.New(X, image.Rect(0, 0, 200, bar_size)),
	}

	up.Connect(func(_ *upower.UPower) {
		sb.Draw()
	})

	w.Connect(func(_ *wicd.Wicd) {
		sb.Draw()
	})

	//w.Connect(sb.wicdChanged)

	return sb
}

func (sb *StatusBar) Connect(img *xgraphics.Image) {
	sb.bar_img = img
}

func (sb *StatusBar) Draw() {
	sb.img.For(func(x, y int) xgraphics.BGRA {
		return status_bg
	})

	// Power Stuff.

	batt, err := sb.up.GetBattery()

	if err != nil {
		return
	}

	//		str := "%3.0f%%"

	if batt.State == upower.Charging || batt.State == upower.Full {
		draw.Draw(sb.img, image.Rect(4, 4, 20, 20),
			charging_img, image.Point{0, 0}, draw.Over)
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

	draw.Draw(sb.img, image.Rect(24, 4, 24+16, 20),
		bimg, image.Point{0, 0}, draw.Over)

	// Networking Stuff

	wifi_img := wifi_dc_img

	// WiFi

	if sb.w.IsConnected() {
		strength, err := sb.w.GetCurrentStrength()

		if err == nil {
			wifi_img = wifiStrengthImage(strength)
		} else {
			fmt.Printf("Error getting strength: %s\n", err)
			wifi_img = wifi_dc_img
		}
	} else {
		wifi_img = wifi_dc_img
	}

	draw.Draw(sb.img, image.Rect(40, 4, 40+16, 20), wifi_img, image.Point{0, 0}, draw.Over)

	draw.Draw(sb.bar_img, image.Rect(612, 0, 812, bar_size), sb.img, image.Point{0, 0}, draw.Over)
}
