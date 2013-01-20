package main

import (
	"image"
)

type SbWifi struct {
	// Do Nothing.
}

func (icon *SbWifi) Attach(sb *StatusBar) {
	// Do Nothing.
}

func (icon *SbWifi) Icon() image.Image {
	return wifi_dc_img
}
