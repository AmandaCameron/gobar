package statbar

import (
	//"fmt"
	"image"
	"image/draw"

	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/xgraphics"
	"github.com/BurntSushi/xgbutil/xwindow"

	"github.com/AmandaCameron/gobar/utils"
	"github.com/AmandaCameron/gobar/utils/system-tray"
)

var (
	status_bg = xgraphics.BGRA{R: 64, G: 64, B: 64, A: 255}
)

type StatusBar struct {
	Width    int
	Height   int
	Position int
	Parent   *xwindow.Window
	X        *xgbutil.XUtil

	img    *xgraphics.Image
	window *xwindow.Window

	tray_offset int
	tray        *systemtray.SystemTray
	tray_icons  []*systemtray.Icon

	items []StatusItem

	//up *upower.UPower
	//w  *wicd.Wicd
}

type StatusItem interface {
	Icon() image.Image
	Attach(sb *StatusBar)
	Size() int
}

// Start Up / Tear Down

func (sb *StatusBar) Init() {
	sb.img = xgraphics.New(sb.X, image.Rect(0, 0, sb.Width, sb.Height))
	var err error

	sb.window, err = xwindow.Create(sb.X, sb.Parent.Id)
	utils.FailMeMaybe(err)

	sb.window.Move(sb.Position, 0)
	sb.window.Resize(sb.Width, sb.Height)
	sb.window.Map()

	sb.img.XSurfaceSet(sb.window.Id)

	utils.FailMeMaybe(sb.initTray())

	sb.Draw()
}

func (sb *StatusBar) Teardown() {
	utils.FailMeMaybe(sb.teardownTray())

	sb.window.Destroy()
	sb.img.Destroy()
}

// API

func (sb *StatusBar) Add(icon StatusItem) {
	sb.items = append(sb.items, icon)
	icon.Attach(sb)

	sb.tray_offset += icon.Size() + 2
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

	sb.img.XDraw()
	sb.img.XPaint(sb.window.Id)
}
