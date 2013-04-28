package systemtray

import (
	"github.com/BurntSushi/xgb/xproto"

	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/xevent"
	"github.com/BurntSushi/xgbutil/xprop"

	"github.com/AmandaCameron/gobar/utils/xembed"
)

type Icon struct {
	Socket *xembed.XEmbedSocket
	Tray   *SystemTray

	classCache string
}

func (tray *SystemTray) newIcon(wid xproto.Window) *Icon {
	sock, err := xembed.NewSocket(tray.X, wid)

	if err != nil {
		tray.Handler.Error(err)
		return nil
	}

	icon := &Icon{
		Socket: sock,
		Tray:   tray,
	}

	return icon
}

func (icon *Icon) Class() string {
	if icon.classCache != "" {
		return icon.classCache
	}

	val, err := xprop.PropValStr(xprop.GetProperty(icon.Socket.X, icon.Socket.Id, "WM_CLASS"))

	if err != nil {
		println("Error getting property:", err)
		return ""
	}

	icon.classCache = val

	return val
}

func (icon *Icon) Embed(x, y int16, parent xproto.Window) error {
	icon.Socket.Listen(xproto.EventMaskStructureNotify)

	xevent.DestroyNotifyFun(func(X *xgbutil.XUtil, ev xevent.DestroyNotifyEvent) {
		icon.Tray.Handler.DelIcon(icon)
	}).Connect(icon.Socket.X, icon.Socket.Id)

	return icon.Socket.Embed(x, y, parent)
}
