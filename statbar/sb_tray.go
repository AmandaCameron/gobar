package statbar

import (
	"github.com/AmandaCameron/gobar/utils"
	"github.com/AmandaCameron/gobar/utils/system-tray"
)

func (sb *StatusBar) initTray() (err error) {
	sb.tray, err = systemtray.New(sb.X)

	if err != nil {
		return err
	}

	sb.tray.Handler = sb

	return nil
}

func (sb *StatusBar) teardownTray() (err error) {
	for _, icon := range sb.tray_icons {
		if err := icon.Socket.Eject(); err != nil {
			println("Error Ejecting:", err.Error())
		}
	}

	return
}

// Tray Handlers

func (sb *StatusBar) NewIcon(icon systemtray.Icon) {
	sb.tray_icons = append(sb.tray_icons, icon)

	utils.FailMeMaybe(icon.Socket.Embed(-16, 4, sb.window.Id))

	icon.Window.Resize(16, 16)

	sb.layout()

	icon.Window.Map()
}

func (sb *StatusBar) Error(err error) {
	utils.Fail(err.Error())
}

// End Handlers.

func (sb *StatusBar) layout() {
	x := sb.tray_offset

	for _, icon := range sb.tray_icons {
		icon.Window.Move(x, 4)
		x += 16
	}
}
