package main

import (
	"image"
	"image/draw"
	"time"

	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/xgraphics"
	"github.com/BurntSushi/xgbutil/xwindow"

	"github.com/AmandaCameron/gobar/images"
	"github.com/AmandaCameron/gobar/utils"
	"github.com/AmandaCameron/gobar/utils/startup"
	"github.com/AmandaCameron/gobar/utils/xdg"
)

type Tracker struct {
	Background xgraphics.BGRA
	Position   int
	Parent     *xwindow.Window
	X          *xgbutil.XUtil
	Size       int

	img      *xgraphics.Image
	window   *xwindow.Window
	stopMe   bool
	pos      int
	launchId string
}

func (t *Tracker) Init() {
	var err error
	t.img = xgraphics.New(t.X, image.Rect(0, 0, t.Size, t.Size))
	t.window, err = xwindow.Create(t.X, t.Parent.Id)
	utils.FailMeMaybe(err)

	t.window.Resize(t.Size, t.Size)
	t.window.Move(t.Position, 0)

	t.img.XSurfaceSet(t.window.Id)

	t.window.Map()

	t.stopMe = true

	l := startup.Listener{
		X:         t.X,
		Callbacks: t,
	}

	l.Initialize()

	t.Draw()
}

// Sastify the commandtray.AppTracker interface

func (t *Tracker) NewApp(app *xdg.Application, launchId string) {
	t.start(launchId)
	go func() {
		time.Sleep(5 * time.Second)
		t.stop()
	}()
}

// Sastify the startup.Contract interface

func (t *Tracker) Add(props map[string]string) {
	t.start(props["ID"])
}

func (t *Tracker) Remove(props map[string]string) {
	if t.launchId == props["ID"] {
		t.stop()
	}
}

func (t *Tracker) Change(props map[string]string) {
	// Do Nothing.
}

// Internal Functions

func (t *Tracker) start(lId string) {
	if !t.stopMe {
		// Refuse to track two apps at once.
		return
	}

	t.launchId = lId

	t.stopMe = false
	go func() {
		for !t.stopMe {
			t.Draw()
			time.Sleep(250 * time.Millisecond)
		}

		t.Draw()
	}()
}

func (t *Tracker) stop() {
	t.stopMe = true
}

func (t *Tracker) Draw() {
	t.img.For(func(x, y int) xgraphics.BGRA {
		return t.Background
	})

	img := images.Tracker_1

	if t.pos == 1 {
		img = images.Tracker_2
	} else if t.pos == 2 {
		img = images.Tracker_3
	} else if t.pos == 3 {
		img = images.Tracker_4
	}
	t.pos++
	t.pos %= 4

	if !t.stopMe {
		draw.Draw(t.img, image.Rect(4, 4, 20, 20), img, image.Point{0, 0}, draw.Over)
	}

	t.img.XDraw()
	t.img.XPaint(t.window.Id)
}
