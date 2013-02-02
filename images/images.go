package images

import (
	"image"
	"image/png"
	"os"

	"github.com/BurntSushi/xdg"

	"github.com/AmandaCameron/gobar/utils"
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

	// App Spinner.
	Tracker_1 image.Image
	Tracker_2 image.Image
	Tracker_3 image.Image
	Tracker_4 image.Image
)

func openImage(paths xdg.Paths, fileName string) image.Image {
	fileName, err := paths.DataFile("images/" + fileName)
	utils.FailMeMaybe(err)

	f, err := os.Open(fileName) //"/home/amanda/.local/share/icons/gobar/" + fileName)

	utils.FailMeMaybe(err)

	defer f.Close()

	img, err := png.Decode(f)

	utils.FailMeMaybe(err)

	return img
}

func Init(paths xdg.Paths) {
	// WiFi
	WifiEnc = openImage(paths, "wifi-enc.png")
	WifiDC = openImage(paths, "wifi-dc.png")
	Wifi_1 = openImage(paths, "wifi-1.png")
	Wifi_2 = openImage(paths, "wifi-2.png")
	Wifi_3 = openImage(paths, "wifi-3.png")
	Wifi_4 = openImage(paths, "wifi-4.png")

	// Status Bar
	Charging = openImage(paths, "ac.png")

	Battery_5 = openImage(paths, "battery-5.png")
	Battery_4 = openImage(paths, "battery-4.png")
	Battery_3 = openImage(paths, "battery-3.png")
	Battery_2 = openImage(paths, "battery-2.png")
	Battery_1 = openImage(paths, "battery-1.png")

	// App Spinner

	Tracker_4 = openImage(paths, "spinner-4.png")
	Tracker_3 = openImage(paths, "spinner-3.png")
	Tracker_2 = openImage(paths, "spinner-2.png")
	Tracker_1 = openImage(paths, "spinner-1.png")
}

func WifiStrengthImage(strength byte) image.Image {
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
