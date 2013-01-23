package main

import (
	// "fmt"
	"image"

	"dbus/nm"
)

type SbNmWifi struct {
	dev *nm.Device
}

func (icon *SbNmWifi) Attach(sb *StatusBar) {
	icon.dev.PropChanged(func(_ *nm.Device) {
		//fmt.Printf("[SB] (Wifi) Props Changed.\n")
		sb.Draw()
	})
}

func (icon *SbNmWifi) Icon() image.Image {
	//return wifi_dc_img
	// fmt.Printf("[SB] (Wifi) Getting Active.\n")

	active, err := icon.dev.GetActive()
	if err != nil {
		// fmt.Printf("[SB] (Wifi) Error: %s\n", err.Error())
		return wifi_dc_img
	}

	if active == nil {
		// fmt.Printf("[SB] (Wifi) Active == nil\n")
		return wifi_dc_img
	}

	// fmt.Printf("[SB] (Wifi) Getting Strength.\n")

	str, err := active.Strength()
	if err != nil {
		// fmt.Printf("[SB] (Wifi) Error: %s\n", err.Error())
		str = 0
	}

	return wifiStrengthImage(int32(str))
}
