package commandtray

import (
	"image"
	"image/color"
	"image/draw"

	"strings"

	"code.google.com/p/jamslam-freetype-go/freetype/truetype"

	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/ewmh"
	"github.com/BurntSushi/xgbutil/keybind"
	"github.com/BurntSushi/xgbutil/xevent"
	"github.com/BurntSushi/xgbutil/xgraphics"
	"github.com/BurntSushi/xgbutil/xwindow"

	"github.com/AmandaCameron/gobar/utils"
)

var font *truetype.Font

type Command interface {
	GetIcon() image.Image
	GetText() string
	Run()
}

type CommandSource interface {
	GetMatches(string, *CommandTray) []Command
	Open(*CommandTray) bool
}

type CommandTray struct {
	img        *xgraphics.Image
	pu_img     *xgraphics.Image
	popup      *xwindow.Window
	window     *xwindow.Window
	is_focused bool
	input      []rune
	mod        string
	selected   int
	cmds       []Command
	active     []CommandSource

	Font     *truetype.Font
	FontSize float64
	Height   int
	Width    int
	X        *xgbutil.XUtil
	Position int
	Parent   *xwindow.Window
}

var (
	sources []CommandSource
)

func Register(cs CommandSource) {
	sources = append(sources, cs)
}

func (ct *CommandTray) Init() {
	var err error

	ct.img = xgraphics.New(ct.X, image.Rect(0, 0, ct.Width, ct.Height))
	ct.pu_img = xgraphics.New(ct.X, image.Rect(0, 0, ct.Width, ct.Height*10))

	ct.window, err = xwindow.Create(ct.X, ct.Parent.Id)
	utils.FailMeMaybe(err)

	ct.img.XSurfaceSet(ct.window.Id)

	ct.window.Move(ct.Position, 0)
	ct.window.Resize(ct.Width, ct.Height)
	ct.window.Map()

	keybind.Initialize(ct.X)
}

func (ct *CommandTray) initPopup() {
	var err error

	ct.popup, err = xwindow.Create(ct.X, ct.X.RootWin())
	utils.FailMeMaybe(err)

	ct.pu_img.For(func(x, y int) xgraphics.BGRA {
		if y%ct.Height == ct.Height-1 || x == 0 || x == ct.Width-1 {
			return xgraphics.BGRA{128, 128, 128, 255}
		}
		return xgraphics.BGRA{64, 64, 64, 255}
	})

	utils.FailMeMaybe(ewmh.WmDesktopSet(ct.X, ct.popup.Id, 0xffffffff))

	utils.FailMeMaybe(ewmh.WmWindowTypeSet(ct.X, ct.popup.Id, []string{
		"_NET_WM_WINDOW_TYPE_DOCK",
	}))

	utils.FailMeMaybe(ct.popup.Listen(xproto.EventMaskKeyPress | xproto.EventMaskStructureNotify))

	ct.popup.Resize(ct.Width, ct.Height*10)
	ct.popup.Move(0, ct.Height)

	ct.pu_img.XSurfaceSet(ct.popup.Id)

	ct.pu_img.XDraw()
	ct.pu_img.XPaint(ct.popup.Id)
}

// API
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

	for _, src := range sources {
		if src.Open(ct) {
			ct.active = append(ct.active, src)
		}
	}

	xevent.MapNotifyFun(func(X *xgbutil.XUtil, e xevent.MapNotifyEvent) {
		ct.popup.Focus()
	}).Connect(ct.X, ct.popup.Id)

	ct.popup.Stack(xproto.StackModeAbove)
	ct.popup.Map()

	ct.Draw()

	ct.getCommands()

	ct.keyPress().Connect(ct.X, ct.popup.Id)

	ct.Draw()
}

func (ct *CommandTray) Blur() {
	ct.is_focused = false

	xevent.Detach(ct.X, ct.popup.Id)

	ct.popup.Destroy()
	ct.popup = nil

	ct.input = []rune{}

	ct.active = nil
	ct.cmds = nil

	ct.Draw()
}

func (ct *CommandTray) Draw() {
	ct.img.For(func(x, y int) xgraphics.BGRA {
		return xgraphics.BGRA{R: 64, G: 64, B: 64, A: 255}
	})

	if ct.is_focused {
		ct.img.Text(5, 4, color.White, ct.FontSize, ct.Font, string(ct.input))

		//draw.Draw(ct.bar_img, image.Rect(0, 0, ct.Width, ct.Height), ct.img, image.Point{0, 0}, draw.Over)

		//ct.bar_img.XDraw()
		//ct.bar_img.XPaint(ct.bar_win.Id)

		ct.img.XDraw()
		ct.img.XPaint(ct.window.Id)

		ct.pu_img.For(func(x, y int) xgraphics.BGRA {
			if y%ct.Height == ct.Height-1 || y == 0 || x == 0 || x == ct.Width-1 {
				return xgraphics.BGRA{128, 128, 128, 255}
			} else if y/ct.Height == ct.selected {
				return xgraphics.BGRA{96, 96, 96, 255}
			}
			return xgraphics.BGRA{64, 64, 64, 255}
		})

		if len(ct.cmds) > 0 {

			for i, cmd := range ct.cmds {
				ico := cmd.GetIcon()

				if ico != nil {
					draw.Draw(ct.pu_img, image.Rect(4, (ct.Height*i)+4, 20, (ct.Height*i)+20),
						ico, image.Point{0, 0}, draw.Over)
				}

				ct.pu_img.Text(24, (ct.Height*(i))+4, color.White, ct.FontSize, ct.Font,
					cmd.GetText())

				if i >= 10 {
					break
				}
			}

			ncmds := len(ct.cmds)
			if ncmds > 10 {
				ncmds = 10
			}

			ct.popup.Resize(ct.Width, ct.Height*ncmds)

		} else if len(ct.input) > 0 {
			ct.pu_img.Text(5, 4, xgraphics.BGRA{128, 128, 128, 255}, ct.FontSize, ct.Font, "No Matches")

			ct.popup.Resize(ct.Width, ct.Height)
		} else {
			ct.pu_img.Text(5, 4, xgraphics.BGRA{128, 128, 128, 255}, ct.FontSize, ct.Font, "Loading...")

			ct.popup.Resize(ct.Width, ct.Height)
		}

		ct.pu_img.XDraw()
		ct.pu_img.XPaint(ct.popup.Id)

	} else {
		ct.img.Text(5, 4, xgraphics.BGRA{128, 128, 128, 255}, ct.FontSize, ct.Font,
			"Press "+ct.mod+" to type a command")

		ct.img.XDraw()
		ct.img.XPaint(ct.window.Id)
	}

}

// Helper Functions

func (ct *CommandTray) getCommands() {
	ct.cmds = make([]Command, 0, 10)

	for _, src := range ct.active {
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
