package main

import (
	"image"
	"image/color"
	"image/draw"

	"strings"

	"code.google.com/p/freetype-go/freetype/truetype"

	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/ewmh"
	"github.com/BurntSushi/xgbutil/keybind"
	"github.com/BurntSushi/xgbutil/xevent"
	"github.com/BurntSushi/xgbutil/xgraphics"
	"github.com/BurntSushi/xgbutil/xwindow"
)

var (
	font      *truetype.Font
	font_size float64 = 12
)

type Command interface {
	GetIcon() image.Image
	GetText() string
	Run()
}

type CommandSource interface {
	GetMatches(string, *CommandTray) []Command
}

type CommandTray struct {
	img        *xgraphics.Image
	bar_img    *xgraphics.Image
	bar_win    *xwindow.Window
	pu_img     *xgraphics.Image
	popup      *xwindow.Window
	is_focused bool
	input      []rune
	mod        string
	selected   int
	cmds       []Command

	X *xgbutil.XUtil
}

var (
	sources []CommandSource
)

func Register(cs CommandSource) {
	sources = append(sources, cs)
}

func NewCommandTray(X *xgbutil.XUtil) *CommandTray {
	font = openFont("/usr/share/fonts/dejavu/DejaVuSans.ttf")

	keybind.Initialize(X)

	ct := &CommandTray{
		X:      X,
		img:    xgraphics.New(X, image.Rect(0, 0, 412, bar_size)),
		pu_img: xgraphics.New(X, image.Rect(0, 0, 412, bar_size*10)),
	}
	return ct
}

func (ct *CommandTray) initPopup() {
	var err error

	ct.popup, err = xwindow.Create(ct.X, ct.X.RootWin())
	failMeMaybe(err)

	ct.pu_img.For(func(x, y int) xgraphics.BGRA {
		if y%bar_size == bar_size-1 || x == 0 || x == 412-1 {
			return xgraphics.BGRA{128, 128, 128, 255}
		}
		return xgraphics.BGRA{64, 64, 64, 255}
	})

	failMeMaybe(ewmh.WmDesktopSet(ct.X, ct.popup.Id, 0xffffffff))

	failMeMaybe(ewmh.WmWindowTypeSet(ct.X, ct.popup.Id, []string{
		"_NET_WM_WINDOW_TYPE_DOCK",
	}))

	failMeMaybe(ct.popup.Listen(xproto.EventMaskKeyPress | xproto.EventMaskStructureNotify))

	ct.popup.Resize(412, bar_size*10)
	ct.popup.Move(0, bar_size)

	ct.pu_img.XSurfaceSet(ct.popup.Id)

	ct.pu_img.XDraw()
	ct.pu_img.XPaint(ct.popup.Id)
}

// API
func (ct *CommandTray) Connect(win *xwindow.Window, img *xgraphics.Image) {
	ct.bar_img = img
	ct.bar_win = win
	ct.Draw()
}

func (ct *CommandTray) Bind(key string) {
	if ct.mod != "" {
		return
	}
	ct.mod = key

	keybind.KeyPressFun(func(_ *xgbutil.XUtil, _ xevent.KeyPressEvent) {
		ct.Focus()
	}).Connect(ct.X, ct.X.RootWin(), ct.mod, true)
}

func (ct *CommandTray) Focus() {
	ct.is_focused = true

	ct.selected = 0

	ct.initPopup()

	xevent.MapNotifyFun(func(X *xgbutil.XUtil, e xevent.MapNotifyEvent) {
		ct.popup.Focus()
	}).Connect(ct.X, ct.popup.Id)

	ct.popup.Stack(xproto.StackModeAbove)
	ct.popup.Map()

	ct.keyPress().Connect(ct.X, ct.popup.Id)

	ct.Draw()
}

func (ct *CommandTray) Blur() {
	ct.is_focused = false

	xevent.Detach(ct.X, ct.popup.Id)

	ct.popup.Destroy()
	ct.popup = nil

	ct.input = []rune{}

	ct.Draw()
}

func (ct *CommandTray) Draw() {
	ct.img.For(func(x, y int) xgraphics.BGRA {
		return xgraphics.BGRA{R: 64, G: 64, B: 64, A: 255}
	})

	if ct.is_focused {
		ct.img.Text(5, 4, color.White, font_size, font, string(ct.input))

		draw.Draw(ct.bar_img, image.Rect(0, 0, 412, bar_size), ct.img, image.Point{0, 0}, draw.Over)

		ct.bar_img.XDraw()
		ct.bar_img.XPaint(ct.bar_win.Id)

		ct.pu_img.For(func(x, y int) xgraphics.BGRA {
			if y%bar_size == bar_size-1 || y == 0 || x == 0 || x == 412-1 {
				return xgraphics.BGRA{128, 128, 128, 255}
			} else if y/bar_size == ct.selected {
				return xgraphics.BGRA{96, 96, 96, 255}
			}
			return xgraphics.BGRA{64, 64, 64, 255}
		})

		if len(ct.cmds) > 0 {

			for i, cmd := range ct.cmds {
				ico := cmd.GetIcon()

				if ico != nil {
					draw.Draw(ct.pu_img, image.Rect(4, (bar_size*i)+4, 20, (bar_size*i)+20),
						ico, image.Point{0, 0}, draw.Over)
				}

				ct.pu_img.Text(24, (bar_size*(i))+4, color.White, font_size, font,
					cmd.GetText())

				if i >= 10 {
					break
				}
			}

			ncmds := len(ct.cmds)
			if ncmds > 10 {
				ncmds = 10
			}

			ct.popup.Resize(412, bar_size*ncmds)

		} else {
			ct.pu_img.Text(5, 4, xgraphics.BGRA{128, 128, 128, 255}, font_size, font, "No Matches")

			ct.popup.Resize(412, bar_size)
		}

		ct.pu_img.XDraw()
		ct.pu_img.XPaint(ct.popup.Id)

	} else {
		ct.img.Text(5, 4, xgraphics.BGRA{128, 128, 128, 255}, font_size, font,
			"Press "+ct.mod+" to type a command")

		draw.Draw(ct.bar_img, image.Rect(0, 0, 412, bar_size), ct.img, image.Point{0, 0}, draw.Over)

	}

}

// Helper Functions

func (ct *CommandTray) getCommands() {
	ct.cmds = make([]Command, 0, 10)

	for _, src := range sources {
		for _, cmd := range src.GetMatches(string(ct.input), ct) {
			ct.cmds = append(ct.cmds, cmd)
		}
	}
}

// Input Handling

func (ct *CommandTray) keyPress() xevent.KeyPressFun {
	f := func(X *xgbutil.XUtil, ev xevent.KeyPressEvent) {

		//		print("Cake:\n")

		if !ct.is_focused {
			//			print("Boom\n")
			return
		}

		mods, kc := keybind.DeduceKeyInfo(ev.State, ev.Detail)

		//		print("Key Press:", mods, kc, "\n")

		switch {
		case keybind.KeyMatch(X, "Escape", mods, kc):
			ct.Blur()

		case keybind.KeyMatch(X, "Return", mods, kc):
			ct.activate()

		case keybind.KeyMatch(X, "Up", mods, kc):
			ct.selected -= 1
			if ct.selected < 0 {
				ct.selected = len(ct.cmds)
			}
			ct.Draw()

		case keybind.KeyMatch(X, "Down", mods, kc):
			if len(ct.cmds) > 0 {
				ct.selected += 1
				ct.selected %= len(ct.cmds)
				ct.Draw()
			}

		case keybind.KeyMatch(X, "BackSpace", mods, kc):
			if len(ct.input) > 0 {

				ct.input = ct.input[:len(ct.input)-1]

				ct.getCommands()

				ct.Draw()
			}

		default:
			s := keybind.LookupString(X, mods, kc)
			if len(s) == 1 {
				ct.input = append(ct.input, rune(strings.ToLower(s)[0]))
				ct.selected = 0

				ct.getCommands()

				ct.Draw()
			}
		}
	}

	return xevent.KeyPressFun(f)
}

func (ct *CommandTray) activate() {
	cmds := ct.cmds

	if len(cmds) > 0 {
		go cmds[ct.selected].Run()
	}

	ct.Blur()
}
