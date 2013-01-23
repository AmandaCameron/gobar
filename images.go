package main

import (
	"image"
	"image/png"

	"os"
)

var (
	// WiFi
	wifi_1_img     image.Image
	wifi_2_img     image.Image
	wifi_3_img     image.Image
	wifi_4_img     image.Image
	wifi_dc_img    image.Image
	wifi_enc_image image.Image

	// Status Bar
	charging_img  image.Image
	battery_5_img image.Image
	battery_4_img image.Image
	battery_3_img image.Image
	battery_2_img image.Image
	battery_1_img image.Image

	//eth_conn_img  image.Image
	//eth_disc_img  image.Image
)

func openImage(fileName string) image.Image {
	f, err := os.Open("/home/amanda/.local/share/icons/gobar/" + fileName)

	failMeMaybe(err)

	defer f.Close()

	img, err := png.Decode(f)

	failMeMaybe(err)

	return img
}

func initImages() {
	// WiFi
	wifi_enc_image = openImage("wifi-enc.png")
	wifi_dc_img = openImage("wifi-dc.png")
	wifi_1_img = openImage("wifi-1.png")
	wifi_2_img = openImage("wifi-2.png")
	wifi_3_img = openImage("wifi-3.png")
	wifi_4_img = openImage("wifi-4.png")

	// Status Bar
	charging_img = openImage("ac.png")

	battery_5_img = openImage("battery-5.png")
	battery_4_img = openImage("battery-4.png")
	battery_3_img = openImage("battery-3.png")
	battery_2_img = openImage("battery-2.png")
	battery_1_img = openImage("battery-1.png")
}

func wifiStrengthImage(strength int32) image.Image {
	switch {
	case strength >= 80:
		return wifi_4_img
	case strength >= 55:
		return wifi_3_img
	case strength >= 30:
		return wifi_2_img
	case strength > 5:
		return wifi_1_img
	}
	return wifi_dc_img
}
