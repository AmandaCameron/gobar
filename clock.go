package main

import (
	"image"
	"time"

	"code.google.com/p/jamslam-freetype-go/freetype/truetype"

	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/xgraphics"
	"github.com/BurntSushi/xgbutil/xwindow"

	"github.com/AmandaCameron/gobar/utils"
)

type Clock struct {
	Background xgraphics.BGRA
	Foreground xgraphics.BGRA
	Width      int
	Height     int
	Position   int
	Font       *truetype.Font
	FontSize   float64
	X          *xgbutil.XUtil
	Parent     *xwindow.Window

	img    *xgraphics.Image
	window *xwindow.Window
}

func (c *Clock) Init() {
	var err error
	c.img = xgraphics.New(c.X, image.Rect(0, 0, c.Width, c.Height))
	c.window, err = xwindow.Create(c.X, c.Parent.Id)
	utils.FailMeMaybe(err)

	c.window.Resize(c.Width, c.Height)
	c.window.Move(c.Position, 0)
	c.img.XSurfaceSet(c.window.Id)

	c.window.Map()

	c.Draw()

	go c.tickTock()
}

func (c *Clock) Draw() {
	c.img.For(func(x, y int) xgraphics.BGRA {
		return c.Background
	})

	now := time.Now().Format("2006-01-02 15:04:05")

	w, h, _ := c.img.Text(0, c.Height*2, c.Foreground, c.FontSize, c.Font, now)

	//println("OH MYYYY: ", w, h)

	c.img.Text((c.Width/2)-(w/2), h-(c.Height*2)-int(c.FontSize/2), c.Foreground, c.FontSize, c.Font, now)

	c.img.XDraw()
	c.img.XPaint(c.window.Id)
}

func (c *Clock) tickTock() {
	ticker := time.Tick(1 * time.Second)
	for {
		select {
		case <-ticker:
			c.Draw()
		}
	}
}
