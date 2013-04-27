package main

import (
	"image"
	"os"
	"os/signal"
	"syscall"
	"time"

	basedir "github.com/BurntSushi/xdg"
	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/ewmh"
	"github.com/BurntSushi/xgbutil/keybind"
	"github.com/BurntSushi/xgbutil/xevent"
	"github.com/BurntSushi/xgbutil/xgraphics"
	"github.com/BurntSushi/xgbutil/xwindow"

	"launchpad.net/~jamesh/go-dbus/trunk"

	"github.com/AmandaCameron/gobar/commandtray"
	"github.com/AmandaCameron/gobar/images"
	"github.com/AmandaCameron/gobar/statbar"

	"github.com/AmandaCameron/gobar/utils"
	"github.com/AmandaCameron/gobar/utils/dbus/upower"
	"github.com/AmandaCameron/gobar/utils/xdg"
	"github.com/AmandaCameron/gobar/utils/xsettings"

	"github.com/AmandaCameron/go.networkmanager"
)

func main() {
	// Signal Handling.

	sigChan := make(chan os.Signal)

	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT)

	paths := basedir.Paths{
		XDGSuffix:    "gobar",
		GoImportPath: "github.com/AmandaCameron/gobar/data",
	}

	file, err := paths.ConfigFile("config.wini")
	utils.FailMeMaybe(err)

	cfg := loadConfig(file) //os.Getenv("HOME") + "/.config/gobar/config.wini")

	// Load Images.

	images.Init(paths)

	// Setup the X Connection

	X, err := xgbutil.NewConn()
	utils.FailMeMaybe(err)

	win, err := xwindow.Create(X, X.RootWin())
	utils.FailMeMaybe(err)

	win.Resize(cfg.BarWidth, cfg.BarSize)

	// Setup the EWMH Stuff

	utils.FailMeMaybe(ewmh.RestackWindow(X, win.Id))

	var strut *ewmh.WmStrutPartial

	if cfg.Position == "Top" {
		strut = &ewmh.WmStrutPartial{
			Top:       uint(cfg.BarSize),
			TopStartX: 0,
			TopEndX:   uint(cfg.BarWidth),
		}

		win.Move(0, 0)
	} else if cfg.Position == "Bottom" {
		strut = &ewmh.WmStrutPartial{
			Bottom:       uint(cfg.BarSize),
			BottomStartX: 0,
			BottomEndX:   uint(cfg.BarWidth),
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

	win.Map()

	keybind.Initialize(X)

	// Get the DE settings, if we can.

	xs, err := xsettings.New(X)
	if err != nil {
		// Maybe this should be an error, maybe not?
		xs = nil
	}

	// Draw the background

	bg := xgraphics.BGRA{
		R: 64,
		G: 64,
		B: 64,
		A: 255,
	}

	img := xgraphics.New(X, image.Rect(0, 0, cfg.BarWidth, cfg.BarSize))

	img.For(func(x, y int) xgraphics.BGRA {
		return bg
	})

	utils.FailMeMaybe(img.XSurfaceSet(win.Id))
	img.XDraw()
	img.XPaint(win.Id)

	// Connect to DBus

	sys, err := dbus.Connect(dbus.SystemBus)
	utils.FailMeMaybe(err)

	// The session bus, too.

	sess, err := dbus.Connect(dbus.SessionBus)
	utils.FailMeMaybe(err)

	// Blah

	x := xdg.New()

	// TODO: How should this fail? I imagine defaulting to gnome is the wrong thing to do,
	// but I'm not really sure what it should do.
	if xs != nil {
		theme, err := xs.GetString("Net/IconThemeName")
		if err == nil {
			x.SetTheme(theme)
		}
	}

	var dev *nm.Device
	var batt *upower.Device

	up := upower.New(sys)
	cli := nm.New(sys)

	if devs, err := cli.GetDevices(); err == nil {
		for _, d := range devs {
			if d.Type() == nm.Wireless {
				dev = d

				break
			}
		}
	}

	if pdevs, err := up.GetDevices(); err == nil {
		for _, d := range pdevs {
			if d.Type() == upower.Battery {
				batt = d
				break
			}
		}
	}

	// Clock

	clck := &Clock{
		X:          X,
		Position:   cfg.Clock.Position,
		Width:      cfg.Clock.Width,
		Height:     cfg.BarSize,
		Parent:     win,
		Format:     cfg.ClockFormat,
		Background: xgraphics.BGRA{R: 48, G: 48, B: 48, A: 255},
		Foreground: xgraphics.BGRA{R: 255, G: 255, B: 255, A: 255},
		Font:       utils.OpenFont(cfg.Clock.Font.Name),
		FontSize:   cfg.Clock.Font.Size,
	}

	clck.Init()

	// App Launch Tracker

	tracker := &Tracker{
		X:          X,
		Position:   cfg.Tracker.Position,
		Size:       cfg.BarSize,
		Background: bg,
		Parent:     win,
	}

	tracker.Init()

	// Command Tray

	ct := &commandtray.CommandTray{
		X:        X,
		Width:    cfg.Command.Width,
		Height:   cfg.BarSize,
		Position: cfg.Command.Position,
		Parent:   win,
		Font:     utils.OpenFont(cfg.Command.Font.Name),
		FontSize: cfg.Command.Font.Size,
	}

	commandtray.Register(commandtray.AppSource{
		Xdg:        x,
		X:          X,
		AppTracker: tracker,
	})

	if sess != nil {
		commandtray.Register(commandtray.GnomeSessionSource{
			Obj: sess.Object("org.gnome.SessionManager", "/org/gnome/SessionManager"),
			Xdg: x,
		})

		commandtray.Register(commandtray.NewShellSource(sess, x))

		commandtray.Register(&commandtray.AppMenuSource{
			Conn: sess,
		})
	}

	// Done, maybe?

	ct.Init()
	ct.Bind(cfg.CommandAccel)
	ct.Draw()

	// Status Bar

	sb := &statbar.StatusBar{
		X:        X,
		Width:    cfg.StatusBar.Width,
		Position: cfg.StatusBar.Position,
		Height:   cfg.BarSize,
		Parent:   win,
	}

	sb.Init()

	if batt != nil {
		sb.Add(&statbar.SbPower{batt})
	}

	if dev != nil {
		sb.Add(&statbar.SbNmWifi{dev})
		commandtray.Register(commandtray.NmSource{dev})
	}

	sb.Draw()

	// My My this anikin guy...

	go func() {
		for {
			select {
			case <-sigChan:
				sb.Teardown()

				time.Sleep(1 * time.Second)
				// Anybody else?
				xevent.Quit(X)
			}
		}
	}()

	xevent.Main(X)
}
