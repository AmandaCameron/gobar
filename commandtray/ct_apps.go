package commandtray

import (
	"image"
	"strings"

	"github.com/BurntSushi/xgbutil"

	"github.com/AmandaCameron/gobar/utils/xdg"
)

type AppSource struct {
	X          *xgbutil.XUtil
	Xdg        *xdg.XDG
	AppTracker AppTracker
}

type AppCommand struct {
	X          *xgbutil.XUtil
	app        *xdg.Application
	AppTracker AppTracker
}

type AppTracker interface {
	NewApp(*xdg.Application, string)
}

func (as AppSource) GetMatches(inp string, ct *CommandTray) []Command {
	cmds := make([]Command, 0)

	if len(inp) < 3 {
		return []Command{}
	}

	for _, app := range as.Xdg.GetApps() {
		if app.NoDisplay {
			continue
		}

		if strings.Contains(strings.ToLower(app.Name), inp) {
			cmds = append(cmds, AppCommand{app: app, X: as.X, AppTracker: as.AppTracker})
		}
	}

	return cmds
}

func (as AppSource) Open(ct *CommandTray) bool {
	// Do Nothing.
	return true
}

func (as AppSource) Close(ct *CommandTray) {
	// Do Nothing.
}

func (ac AppCommand) GetText() string {
	return "Launch " + ac.app.Name
}

func (ac AppCommand) GetIcon() image.Image {
	return ac.app.FindIcon(16)
}

func (ac AppCommand) Run() {
	deskId := ac.app.Run(ac.X.TimeGet())

	if ac.AppTracker != nil {
		ac.AppTracker.NewApp(ac.app, deskId)
	}
}
