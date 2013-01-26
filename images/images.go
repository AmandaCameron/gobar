package images

import (
	"image"
	"image/png"

	"github.com/AmandaCameron/gobar/utils"

	"os"
)

var (
	// WiFi
	Wifi_1  image.Image
	Wifi_2  image.Image
	Wifi_3  image.Image
	Wifi_4  image.Image
	WifiDC  image.Image
	WifiEnc image.Image

	// Power Bar.
	Charging  image.Image
	Battery_5 image.Image
	Battery_4 image.Image
	Battery_3 image.Image
	Battery_2 image.Image
	Battery_1 image.Image
)

func openImage(fileName string) image.Image {
	f, err := os.Open("/home/amanda/.local/share/icons/gobar/" + fileName)

	utils.FailMeMaybe(err)

	defer f.Close()

	img, err := png.Decode(f)

	utils.FailMeMaybe(err)

	return img
}

func Init() {
	// WiFi
	WifiEnc = openImage("wifi-enc.png")
	WifiDC = openImage("wifi-dc.png")
	Wifi_1 = openImage("wifi-1.png")
	Wifi_2 = openImage("wifi-2.png")
	Wifi_3 = openImage("wifi-3.png")
	Wifi_4 = openImage("wifi-4.png")

	// Status Bar
	Charging = openImage("ac.png")

	Battery_5 = openImage("battery-5.png")
	Battery_4 = openImage("battery-4.png")
	Battery_3 = openImage("battery-3.png")
	Battery_2 = openImage("battery-2.png")
	Battery_1 = openImage("battery-1.png")
}

func WifiStrengthImage(strength int32) image.Image {
	switch {
	case strength >= 80:
		return Wifi_4
	case strength >= 55:
		return Wifi_3
	case strength >= 30:
		return Wifi_2
	case strength > 5:
		return Wifi_1
	}
	return WifiDC
}
