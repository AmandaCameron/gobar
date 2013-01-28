package commandtray

import (
	"image"
	"strings"

	"launchpad.net/~jamesh/go-dbus/trunk"

	"github.com/AmandaCameron/gobar/utils/xdg"
)

type GnomeSessionSource struct {
	Obj *dbus.ObjectProxy
	Xdg *xdg.XDG
}

type GnomeLogoutCommand struct {
	Obj *dbus.ObjectProxy
	Xdg *xdg.XDG
}

type GnomeShutdownCommand struct {
	Obj *dbus.ObjectProxy
	Xdg *xdg.XDG
}

type GnomeRestartCommand struct {
	Obj *dbus.ObjectProxy
	Xdg *xdg.XDG
}

func (gss GnomeSessionSource) GetMatches(inp string, ct *CommandTray) []Command {
	if gss.Obj == nil || inp == "" {
		return nil
	}

	var ret []Command

	if strings.Contains("log out", inp) {
		ret = append(ret, GnomeLogoutCommand{
			gss.Obj, gss.Xdg,
		})
	}

	if strings.Contains("shutdown", inp) {
		ret = append(ret, GnomeShutdownCommand{
			gss.Obj, gss.Xdg,
		})
	}

	if strings.Contains("restart", inp) {
		// Includt the shutdown command for funsies.

		ret = append(ret, GnomeShutdownCommand{
			gss.Obj, gss.Xdg,
		})

		ret = append(ret, GnomeRestartCommand{
			gss.Obj, gss.Xdg,
		})
	}

	return ret
}

func (gss GnomeSessionSource) Open(ct *CommandTray) bool {
	return true
}

func (gss GnomeSessionSource) Close() {
	// Do Nothing.
}

// First up -- Logout

func (cmd GnomeLogoutCommand) GetText() string {
	return "Log Out"
}

func (cmd GnomeLogoutCommand) GetIcon() image.Image {
	return cmd.Xdg.GetIcon("system-log-out", 16)
}

func (cmd GnomeLogoutCommand) Run() {
	cmd.Obj.Call("org.gnome.SessionManager", "Logout", uint32(0))
}

// And next -- shutdown.

func (cmd GnomeShutdownCommand) GetText() string {
	return "Shutdown"
}

func (cmd GnomeShutdownCommand) GetIcon() image.Image {
	return cmd.Xdg.GetIcon("system-shutdown", 16)
}

func (cmd GnomeShutdownCommand) Run() {
	cmd.Obj.Call("org.gnome.SessionManager", "Shutdown")
}

// And last, but not least, restart!

func (cmd GnomeRestartCommand) GetText() string {
	return "Restart"
}

func (cmd GnomeRestartCommand) GetIcon() image.Image {
	return cmd.Xdg.GetIcon("system-reboot", 16)
}

func (cmd GnomeRestartCommand) Run() {
	cmd.Obj.Call("org.gnome.SessionManager", "Reboot")
}
