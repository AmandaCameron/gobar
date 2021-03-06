package systemtray

import (
	"fmt"

	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/xevent"
	"github.com/BurntSushi/xgbutil/xprop"
	"github.com/BurntSushi/xgbutil/xwindow"
)

type SystemTray struct {
	X       *xgbutil.XUtil
	Handler Callbacks

	wid *xwindow.Window
}

type Callbacks interface {
	NewIcon(*Icon)
	DelIcon(*Icon)
	Error(error)
}

var sysTrayAtom xproto.Atom
var sysTrayMsgAtom xproto.Atom
var managerEventAtom xproto.Atom

func New(X *xgbutil.XUtil) (*SystemTray, error) {
	tray := &SystemTray{
		X: X,
	}

	var err error

	if sysTrayAtom == 0 {
		sysTrayAtom, err = xprop.Atom(X, "_NET_SYSTEM_TRAY_S0", false)

		if err != nil {
			return nil, err
		}
	}

	if sysTrayMsgAtom == 0 {
		sysTrayMsgAtom, err = xprop.Atom(X, "_NET_SYSTEM_TRAY_OPCODE", false)

		if err != nil {
			return nil, err
		}
	}

	if managerEventAtom == 0 {
		managerEventAtom, err = xprop.Atom(X, "MANAGER", false)

		if err != nil {
			return nil, err
		}
	}

	tray.wid, err = xwindow.Create(X, X.RootWin())

	if err != nil {
		return nil, err
	}

	ts, err := currentTime(X)

	if err != nil {
		return nil, err
	}

	X.TimeSet(ts)

	// tray.wid.Listen(xproto.EventMaskNoEvent | xproto.EventMaskPropertyChange)

	err = xproto.SetSelectionOwnerChecked(tray.X.Conn(), tray.wid.Id, sysTrayAtom, tray.X.TimeGet()).Check()

	if err != nil {
		return nil, err
	}

	reply, err := xproto.GetSelectionOwner(X.Conn(), sysTrayAtom).Reply()
	if err != nil {
		return nil, err
	}

	if reply.Owner != tray.wid.Id {
		return nil, fmt.Errorf("Could not get ownership of the thingy-thing.")
	}

	evt, err := xevent.NewClientMessage(32, X.RootWin(), managerEventAtom,
		int(X.TimeGet()), int(sysTrayAtom), int(tray.wid.Id))

	if err != nil {
		return nil, err
	}

	if err = xevent.SendRootEvent(X, evt, xproto.EventMaskStructureNotify); err != nil {
		return nil, err
	}

	xevent.ClientMessageFun(func(_ *xgbutil.XUtil, ev xevent.ClientMessageEvent) {
		tray.event(ev)
	}).Connect(tray.X, tray.wid.Id)

	return tray, nil
}

func (tray *SystemTray) Teardown() error {
	// TODO: Tray Icon Tracking?

	currentTime(tray.X)
	return xproto.SetSelectionOwnerChecked(tray.X.Conn(), xproto.WindowNone, sysTrayAtom, tray.X.TimeGet()).Check()
}

func (tray *SystemTray) event(ev xevent.ClientMessageEvent) {
	if ev.Format != 32 {
		return
	}

	if ev.Type != sysTrayMsgAtom {
		return
	}

	opCode := ev.Data.Data32[1]

	if opCode == 0 {
		// SYSTEM_TRAY_REQUEST_DOCK

		icon := tray.newIcon(xproto.Window(ev.Data.Data32[2]))

		if icon != nil {
			tray.Handler.NewIcon(icon)
		}
	} else if opCode == 1 {
		// Not Implemented.

	} else if opCode == 2 {
		// Not Implemented.

	} else {
		// Do Nothing for now.
		fmt.Printf("[SystemTray] Got unknown opcode: %d", opCode)
	}
}
