package main

import (
	"image"
	"image/color"
	"image/draw"
	"os"
	"time"

	"code.google.com/p/freetype-go/freetype/truetype"

	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/ewmh"
	"github.com/BurntSushi/xgbutil/xevent"
	"github.com/BurntSushi/xgbutil/xgraphics"
	"github.com/BurntSushi/xgbutil/xwindow"

	// "github.com/norisatir/go-dbus"
	"launchpad.net/~jamesh/go-dbus/trunk"

	"dbus/upower"
	//"dbus/wicd"
)

var (
	clock_font      string  = "/usr/share/fonts/dejavu/DejaVuSans-Bold.ttf"
	clock_font_size float64 = 12

	workspace_font      string  = "/usr/share/fonts/dejavu/DejaVuSans.ttf"
	workspace_font_size float64 = 12

	title_font      string  = "/usr/share/fonts/dejavu/DejaVuSans.ttf"
	title_font_size float64 = 12

	status_font      string  = "/usr/share/fonts/dejavu/DejaVuSans.ttf"
	status_font_size float64 = 12

	bar_size int = 24
)

func openFont(fileName string) *truetype.Font {
	f, err := os.Open(fileName)

	failMeMaybe(err)

	defer f.Close()

	fnt, err := xgraphics.ParseFont(f)

	failMeMaybe(err)

	return fnt
}

func main() {

	// Load Images.

	initImages()

	// Setup the X Connection

	X, err := xgbutil.NewConn()
	failMeMaybe(err)

	win, err := xwindow.Create(X, X.RootWin())
	failMeMaybe(err)

	win.Move(0, 0)
	win.Resize(1024, bar_size)

	// Setup the EWMH Stuff

	failMeMaybe(ewmh.RestackWindow(X, win.Id))

	strut := &ewmh.WmStrutPartial{
		Top:       uint(bar_size),
		TopStartX: 0,
		TopEndX:   1024,
	}

	failMeMaybe(ewmh.WmStrutPartialSet(X, win.Id, strut))

	failMeMaybe(ewmh.WmWindowTypeSet(X, win.Id, []string{
		"_NET_WM_WINDOW_TYPE_DOCK",
	}))

	// Put us everywhere.
	failMeMaybe(ewmh.WmDesktopSet(X, win.Id, 0xFFFFFFFF))

	// Show the window?

	win.Map()

	// Draw the background

	bg := xgraphics.BGRA{
		R: 64,
		G: 64,
		B: 64,
		A: 255,
	}

	img := xgraphics.New(X, image.Rect(0, 0, 1024, bar_size))

	img.For(func(x, y int) xgraphics.BGRA {
		return bg
	})

	failMeMaybe(img.XSurfaceSet(win.Id))
	img.XDraw()
	img.XPaint(win.Id)

	// Connect to DBus

	conn, err := dbus.Connect(dbus.SystemBus)
	failMeMaybe(err)

	failMeMaybe(conn.Authenticate())

	up, err := upower.New(conn)
	failMeMaybe(err)

	//w, err := wicd.New(conn)
	//failMeMaybe(err)
	//var w *wicd.Wicd

	// Command Tray

	ct := NewCommandTray(X)

	ct.Bind("Mod4-n")
	ct.Connect(win, img)

	// Register(NetSource{w: w})

	// Status Bar

	sb := NewStatusBar(X)

	sb.Add(&SbPower{up})
	sb.Add(&SbWifi{})

	sb.Connect(img)
	sb.Draw()

	// My My this anikin guy...

	go drawClock(X, img, win)
	go drawWorkspace(X, img, win)
	//	go drawStatus(X, img, up, w)
	//	go ct.Draw()

	// print("Oh myyyy.\n")

	xevent.Main(X)
}

func drawClock(X *xgbutil.XUtil, bar_img *xgraphics.Image, win *xwindow.Window) {
	//img := bar_img.SubImage(image.Rect(0, 412, 612, bar_size))
	img := xgraphics.New(X, image.Rect(0, 0, 200, bar_size))
	// failMeMaybe(err)
	fnt := openFont(clock_font)

	for {
		clock_bg := xgraphics.BGRA{
			R: 48,
			G: 48,
			B: 48,
			A: 255,
		}

		img.For(func(x, y int) xgraphics.BGRA {
			return clock_bg
		})

		now := time.Now()
		str := now.Format("2006-01-02 15:04:05")

		_, h := xgraphics.TextMaxExtents(fnt, clock_font_size, str)

		img.Text(25, (bar_size/2)-(h/2), color.White, clock_font_size, fnt, str)

		draw.Draw(bar_img, image.Rect(412, 0, 612, bar_size), img, image.Point{0, 0}, draw.Over)

		//img.XPaint(win.Id)

		time.Sleep(1 * time.Second)
	}
}

func drawWorkspace(X *xgbutil.XUtil, bar_img *xgraphics.Image, win *xwindow.Window) {
	img := xgraphics.New(X, image.Rect(0, 0, 75, bar_size))
	//img := bar_img.SubImage(image.Rect(949, 0, 1024, bar_size))

	// failMeMaybe(err)
	fnt := openFont(workspace_font)

	workspace_bg := xgraphics.BGRA{
		R: 32,
		G: 32,
		B: 128,
		A: 255,
	}

	for {
		img.For(func(x, y int) xgraphics.BGRA {
			return workspace_bg
		})

		dsk, err := ewmh.CurrentDesktopGet(X)

		if err != nil {
			continue
		}

		desks, err := ewmh.DesktopNamesGet(X)

		if err != nil {
			continue
		}

		str := desks[dsk]

		img.Text(5, bar_size/2-6, color.White, workspace_font_size, fnt, str)

		draw.Draw(bar_img, image.Rect(949, 0, 1024, 24), img, image.Point{0, 0}, draw.Over)

		//		img.XPaint(win.Id)

		bar_img.XDraw()
		bar_img.XPaint(win.Id)

		time.Sleep(500 * time.Millisecond)
	}
}
