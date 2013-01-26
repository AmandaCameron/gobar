package statbar

import (
	//"fmt"
	"image"
	"image/draw"

	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/xgraphics"

	// "dbus/upower"
	// "dbus/wicd"
)

var (
	status_bg = xgraphics.BGRA{R: 64, G: 64, B: 64, A: 255}
)

type StatusBar struct {
	Width  int
	Height int
	X      *xgbutil.XUtil

	img     *xgraphics.Image
	bar_img *xgraphics.Image

	items []StatusItem

	//up *upower.UPower
	//w  *wicd.Wicd
}

type StatusItem interface {
	Icon() image.Image
	Attach(sb *StatusBar)
}

func (sb *StatusBar) Init() {
	sb.img = xgraphics.New(sb.X, image.Rect(0, 0, sb.Width, sb.Height))
}

func (sb *StatusBar) Add(icon StatusItem) {
	sb.items = append(sb.items, icon)
	icon.Attach(sb)
}

func (sb *StatusBar) Connect(img *xgraphics.Image) {
	sb.bar_img = img
}

func (sb *StatusBar) Draw() {
	sb.img.For(func(x, y int) xgraphics.BGRA {
		return status_bg
	})

	x := 0

	for _, item := range sb.items {
		if x > sb.img.Bounds().Dx() {
			break
		}

		img := item.Icon()
		img_width := img.Bounds().Dx()

		draw.Draw(sb.img, image.Rect(x, 4, x+img_width, 20), img, image.Point{0, 0}, draw.Over)

		x += img_width + 2
	}

	draw.Draw(sb.bar_img, image.Rect(612, 0, 812, sb.Height), sb.img, image.Point{0, 0}, draw.Over)
}
