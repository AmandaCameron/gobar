package main

import (
	"fmt"
	"os"

	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/xevent"
	"github.com/BurntSushi/xgbutil/xwindow"
)

func failMeMaybe(err error) {
	if err != nil {
		fail(err.Error())
	}
}

func fail(err string) {
	fmt.Println("Error: ", err)
	os.Exit(1)
}

type Handler struct {
	X     *xgbutil.XUtil
	Icons []STIcon

	Window *xwindow.Window
}

func (hdlr *Handler) NewIcon(icon STIcon) {
	fmt.Printf("Got new icon! %p", icon)

	hdlr.Icons = append(hdlr.Icons, icon)

	hdlr.Window.Resize(24, len(hdlr.Icons)*24)

	icon.Socket.Embed(int16(len(hdlr.Icons)*24-24), 4, hdlr.Window.Id)
	icon.Window.Resize(16, 16)
	icon.Window.Map()
}

func (hdlr *Handler) Error(err error) {
	failMeMaybe(err)
}

func main() {
	X, err := xgbutil.NewConn()
	failMeMaybe(err)

	tray, err := NewSystemTray(X)
	failMeMaybe(err)

	hdlr := &Handler{
		X: X,
	}

	hdlr.Window, err = xwindow.Create(X, X.RootWin())
	failMeMaybe(err)

	hdlr.Window.Move(0, 24)
	hdlr.Window.Map()

	tray.Handler = hdlr

	xevent.Main(X)
}
