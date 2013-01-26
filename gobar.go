package main

import (
	"image"
	"image/color"
	"image/draw"
	"os"
	"time"

	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/ewmh"
	"github.com/BurntSushi/xgbutil/xevent"
	"github.com/BurntSushi/xgbutil/xgraphics"
	"github.com/BurntSushi/xgbutil/xwindow"

	"launchpad.net/~jamesh/go-dbus/trunk"

	"github.com/AmandaCameron/gobar/commandtray"
	"github.com/AmandaCameron/gobar/images"
	"github.com/AmandaCameron/gobar/statbar"

	"github.com/AmandaCameron/gobar/utils"
	"github.com/AmandaCameron/gobar/utils/dbus/nm"
	"github.com/AmandaCameron/gobar/utils/dbus/upower"
	"github.com/AmandaCameron/gobar/utils/xdg"
)

func main() {

	cfg := loadConfig(os.Getenv("HOME") + "/.config/gobar/config.wini")

	// Load Images.

	images.Init()

	// Setup the X Connection

	X, err := xgbutil.NewConn()
	utils.FailMeMaybe(err)

	win, err := xwindow.Create(X, X.RootWin())
	utils.FailMeMaybe(err)

	win.Resize(1024, cfg.BarSize)

	// Setup the EWMH Stuff

	utils.FailMeMaybe(ewmh.RestackWindow(X, win.Id))

	var strut *ewmh.WmStrutPartial

	if cfg.Position == "Top" {
		strut = &ewmh.WmStrutPartial{
			Top:       uint(cfg.BarSize),
			TopStartX: 0,
			TopEndX:   1024,
		}
		win.Move(0, 0)

	} else if cfg.Position == "Bottom" {
		strut = &ewmh.WmStrutPartial{
			Bottom:       uint(cfg.BarSize),
			BottomStartX: 0,
			BottomEndX:   1024,
		}
		win.Move(0, 600-cfg.BarSize)
	} else {
		println("Invalid Position:", cfg.Position)
		os.Exit(1)
	}

	utils.FailMeMaybe(ewmh.WmStrutPartialSet(X, win.Id, strut))

	utils.FailMeMaybe(ewmh.WmWindowTypeSet(X, win.Id, []string{
		"_NET_WM_WINDOW_TYPE_DOCK",
	}))

	// Put us everywhere.
	utils.FailMeMaybe(ewmh.WmDesktopSet(X, win.Id, 0xFFFFFFFF))

	// Show the window?

	win.Map()

	// Draw the background

	bg := xgraphics.BGRA{
		R: 64,
		G: 64,
		B: 64,
		A: 255,
	}

	img := xgraphics.New(X, image.Rect(0, 0, 1024, cfg.BarSize))

	img.For(func(x, y int) xgraphics.BGRA {
		return bg
	})

	utils.FailMeMaybe(img.XSurfaceSet(win.Id))
	img.XDraw()
	img.XPaint(win.Id)

	// Connect to DBus

	sys, err := dbus.Connect(dbus.SystemBus)
	utils.FailMeMaybe(err)

	utils.FailMeMaybe(sys.Authenticate())

	// TODO: This will come in handy sometime.

	sess, err := dbus.Connect(dbus.SessionBus)
	utils.FailMeMaybe(err)

	utils.FailMeMaybe(sess.Authenticate())

	// Blash

	x := xdg.New()
	x.SetTheme("gnome")

	up := upower.New(sys)
	cli := nm.New(sys)

	devs, err := cli.GetDevices()
	utils.FailMeMaybe(err)

	var dev *nm.Device
	var batt *upower.Device

	for _, d := range devs {
		if d.Type() == nm.Wireless {
			dev = d

			break
		}
	}

	pdevs, err := up.GetDevices()
	utils.FailMeMaybe(err)

	for _, d := range pdevs {
		if d.Type() == upower.Battery {
			batt = d
			break
		}
	}

	// Command Tray

	//ct := commandtray.New(X)

	ct := commandtray.CommandTray{
		X:        X,
		Width:    412,
		Height:   cfg.BarSize,
		Font:     utils.OpenFont(cfg.CommandFont.Name),
		FontSize: cfg.CommandFont.Size,
	}

	ct.Init()
	ct.Bind(cfg.CommandCombo)
	ct.Connect(win, img)

	// Register(NetSource{w: w})

	commandtray.Register(commandtray.NewShellSource(sess, x))

	commandtray.Register(&commandtray.AppMenuSource{
		Conn: sess,
	})

	commandtray.Register(commandtray.AppSource{
		Xdg: x,
	})

	// Status Bar

	sb := &statbar.StatusBar{
		X:      X,
		Width:  200,
		Height: cfg.BarSize,
	} //statbar.New(X)

	sb.Init()

	if batt != nil {
		sb.Add(&statbar.SbPower{batt})
	}

	if dev != nil {
		sb.Add(&statbar.SbNmWifi{dev})
		commandtray.Register(commandtray.NmSource{dev})
	}

	sb.Connect(img)
	sb.Draw()

	// My My this anikin guy...

	go drawClock(X, img, win, cfg)
	//go drawWorkspace(X, img, win)

	xevent.Main(X)
}

func drawClock(X *xgbutil.XUtil, bar_img *xgraphics.Image, win *xwindow.Window, cfg *Config) {
	img := xgraphics.New(X, image.Rect(0, 0, 200, cfg.BarSize))

	fnt := utils.OpenFont(cfg.ClockFont.Name)

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

		_, h := xgraphics.TextMaxExtents(fnt, cfg.ClockFont.Size, str)

		img.Text(25, (cfg.BarSize/2)-(h/2), color.White, cfg.ClockFont.Size, fnt, str)

		draw.Draw(bar_img, image.Rect(412, 0, 612, cfg.BarSize), img, image.Point{0, 0}, draw.Over)

		//img.XPaint(win.Id)

		bar_img.XDraw()
		bar_img.XPaint(win.Id)

		time.Sleep(1 * time.Second)
	}
}

// func drawWorkspace(X *xgbutil.XUtil, bar_img graphics.Image, win *xwindow.Window, cfg *Config) {
// 	img := xgraphics.New(X, image.Rect(0, 0, 75, cfg.BarSize))

// 	workspace_bg := xgraphics.BGRA{
// 		R: 32,
// 		G: 32,
// 		B: 128,
// 		A: 255,
// 	}

// 	for {
// 		time.Sleep(500 * time.Millisecond)

// 		img.For(func(x, y int) xgraphics.BGRA {
// 			return workspace_bg
// 		})

// 		dsk, err := ewmh.CurrentDesktopGet(X)

// 		if err != nil {
// 			continue
// 		}

// 		desks, err := ewmh.DesktopNamesGet(X)
// 		if err != nil {
// 			continue
// 		}

// 		if int(dsk) >= len(desks) {
// 			continue
// 		}

// 		str := desks[dsk]

// 		img.Text(5, cfg.BarSize/2-6, color.White, workspace_font_size, fnt, str)

// 		draw.Draw(bar_img, image.Rect(949, 0, 1024, 24), img, image.Point{0, 0}, draw.Over)

// 		bar_img.XDraw()
// 		bar_img.XPaint(win.Id)
// 	}
// }
