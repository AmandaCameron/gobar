package xembed

import (
	"fmt"

	"github.com/BurntSushi/xgb/xproto"

	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/xevent"
	"github.com/BurntSushi/xgbutil/xprop"
)

type XEmbedSocket struct {
	X       *xgbutil.XUtil
	Version uint32
	Flags   SocketFlags

	id     xproto.Window
	parent xproto.Window
}

type SocketFlags uint32

type XEmbedMessage uint32

var xembedAtom xproto.Atom

const (
	MsgEmbeddedNotify XEmbedMessage = iota
	MsgWindowActivate
	MsgWindowDeactivate
	MsgRequestFocus
	MsgFocusIn
	MsgFocusOut
	MsgFocusNext
	MsgFocusPrev
	MsgGrabKey   // Deprecated.
	MsgUngrabKey // Deprecated.
	MsgModalityOn
	MsgModalityOff
	MsgRegisterAccelerator
	MsgUnregisterAccelerator
	MsgActivateAccelerator
)

const (
	SocketMapped SocketFlags = 1 << iota
)

func NewSocket(X *xgbutil.XUtil, wid xproto.Window) (*XEmbedSocket, error) {
	sock := &XEmbedSocket{
		id: wid,
		X:  X,
	}

	if err := sock.load(); err != nil {
		return nil, err
	}

	return sock, nil
}

func (sock *XEmbedSocket) load() error {
	vals, err := xprop.PropValNums(xprop.GetProperty(sock.X, sock.id, "_XEMBED_INFO"))
	if err != nil {
		return err
	}

	if len(vals) < 2 {
		return fmt.Errorf("Expected 2 nums in property -- got %d", len(vals))
	}

	sock.Version = uint32(vals[0])
	sock.Flags = SocketFlags(vals[1])

	return nil
}

func (sock *XEmbedSocket) sendMessage(msg XEmbedMessage, detail, data1, data2 int) (err error) {
	if xembedAtom == 0 {
		xembedAtom, err = xprop.Atom(sock.X, "_XEMBED", false)

		if err != nil {
			return
		}
	}

	clientMsg, err := xevent.NewClientMessage(32, sock.id, xembedAtom, detail, data1, data2)

	if err == nil {
		xproto.SendEvent(sock.X.Conn(), false, sock.id, xproto.EventMaskNoEvent, string(clientMsg.Bytes()))
	}

	return
}

func (sock *XEmbedSocket) Embed(x, y int16, window xproto.Window) error {
	if err := xproto.ReparentWindowChecked(sock.X.Conn(), sock.id, window, x, y).Check(); err != nil {
		return err
	}

	sock.sendMessage(MsgEmbeddedNotify, 0, int(window), 0)

	sock.parent = window

	return nil
}

func (sock *XEmbedSocket) Eject() error {
	if err := xproto.UnmapWindowChecked(sock.X.Conn(), sock.id).Check(); err != nil {
		return err
	}

	return xproto.ReparentWindowChecked(sock.X.Conn(), sock.id, sock.X.RootWin(), 1, 1).Check()
}
