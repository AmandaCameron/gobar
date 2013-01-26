package commandtray

import (
	"fmt"
	"image"
	"strings"

	"github.com/AmandaCameron/gobar/utils/xdg"
)

type AppSource struct {
	Xdg *xdg.XDG
}

type AppCommand struct {
	app *xdg.Application
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
			cmds = append(cmds, AppCommand{app: app})
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
	err := ac.app.Run()

	if err != nil {
		fmt.Printf("Error Launching App: ", err.Error())
	}
}
