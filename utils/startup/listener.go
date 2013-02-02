package startup

import (
	"fmt"
	"strings"

	"github.com/BurntSushi/xgb/xproto"

	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/xevent"
	"github.com/BurntSushi/xgbutil/xprop"
	"github.com/BurntSushi/xgbutil/xwindow"
)

var (
	typeAtom      xproto.Atom
	typeAtomBegin xproto.Atom
)

type Listener struct {
	X         *xgbutil.XUtil
	Callbacks Contract
	msgBuff   map[xproto.Window][]byte
}

type Contract interface {
	Add(map[string]string)
	Remove(map[string]string)
	Change(map[string]string)
}

func (l *Listener) Initialize() error {

	if l.X == nil {
		return fmt.Errorf("X must not be nil.")
	}

	if typeAtom == 0 && typeAtomBegin == 0 {
		var err error

		typeAtom, err = xprop.Atom(l.X, "_NET_STARTUP_INFO", false)
		if err != nil {
			return err
		}

		typeAtomBegin, err = xprop.Atom(l.X, "_NET_STARTUP_INFO_BEGIN", false)
		if err != nil {
			return err
		}
	}

	if l.Callbacks == nil {
		return fmt.Errorf("Callbacks must not be nil")
	}

	if err := xwindow.New(l.X, l.X.RootWin()).Listen(xproto.EventMaskPropertyChange); err != nil {
		return err
	}

	l.msgBuff = make(map[xproto.Window][]byte)

	xevent.HookFun(func(_ *xgbutil.XUtil, ev interface{}) bool {
		return l.hook(ev)
	}).Connect(l.X)

	return nil
}

func (l *Listener) hook(ev interface{}) bool {
	e, ok := ev.(xproto.ClientMessageEvent)
	if !ok {
		return true
	}

	if e.Type != typeAtomBegin && e.Type != typeAtom {
		return true
	}

	if e.Format != 8 {
		return true
	}

	// Done checking, PROCESS DAMNIT!

	if e.Type == typeAtomBegin {
		l.beginMsg(e)
	} else if e.Type == typeAtom {
		l.continueMsg(e)
	}

	return false
}

func (l *Listener) beginMsg(ev xproto.ClientMessageEvent) {
	l.msgBuff[ev.Window] = ev.Data.Data8
}

func (l *Listener) continueMsg(ev xproto.ClientMessageEvent) {
	buff := l.msgBuff[ev.Window]
	l.msgBuff[ev.Window] = nil

	for _, b := range ev.Data.Data8 {
		if b == 0 {
			l.processMsg(string(buff))
			return
		}
		buff = append(buff, b)
	}

	l.msgBuff[ev.Window] = buff
}

func (l *Listener) processMsg(msg string) {
	parts := strings.SplitN(msg, ": ", 2)

	cmd := parts[0]
	params := parts[1]

	var buff []rune
	var quo, esc bool
	var tmp []string

	for _, b := range params {
		if b == '"' && !esc {
			quo = !quo
		} else if b == '\\' && !esc {
			esc = true
		} else if b == ' ' && !quo && !esc {
			tmp = append(tmp, string(buff))
			buff = nil
		} else {
			buff = append(buff, b)
		}
	}

	tmp = append(tmp, string(buff))

	p := make(map[string]string)

	for _, kv := range tmp {
		if !strings.Contains(kv, "=") {
			return
		}

		parts := strings.SplitN(kv, "=", 2)
		key, val := parts[0], parts[1]

		p[key] = val
	}

	if cmd == "add" {
		l.Callbacks.Add(p)
	} else if cmd == "change" {
		l.Callbacks.Change(p)
	} else if cmd == "remove" {
		l.Callbacks.Remove(p)
	}
}
