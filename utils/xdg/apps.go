package xdg

import (
	"image"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	//"time"

	"code.google.com/p/goconf/conf"
)

type Application struct {
	xdg *XDG
	cfg *conf.ConfigFile

	Exec      string
	Icon      string
	Name      string
	NoDisplay bool
	Type      string
}

func (xdg *XDG) initApps() {
	xdg.apps = make(map[string]*Application)

	for _, path := range []string{
		"/usr/share/applications",
		"/usr/local/share/applications",
		os.Getenv("HOME") + "/.local/share/applications",
	} {
		apps, _ := ioutil.ReadDir(path)
		for _, file := range apps {
			if !strings.HasSuffix(file.Name(), ".desktop") {
				continue // Skip non-desktop files?
			}

			entry, err := xdg.LoadApplication(path + "/" + file.Name())
			if err != nil {
				println("Error Processing "+path+"/"+file.Name(), ": ", err.Error())
				continue
			}

			app, ok := xdg.apps[entry.Name]
			if ok {
				app.Icon = entry.Icon
				app.Exec = entry.Exec
				app.NoDisplay = entry.NoDisplay
				app.Type = entry.Type

				app.cfg = entry.cfg
			} else {
				xdg.apps[entry.Name] = entry
			}
		}
	}
}

func (xdg *XDG) LoadApplication(path string) (*Application, error) {
	cfg, err := conf.ReadConfigFile(path)
	if err != nil {
		return nil, err
	}

	app := &Application{}

	DE := "Desktop Entry"

	app.Exec, err = cfg.GetRawString(DE, "Exec")
	if err != nil {
		return nil, err
	}

	app.Icon, err = cfg.GetRawString(DE, "Icon")
	if err != nil {
		app.Icon = ""
	}

	app.Name, err = cfg.GetRawString(DE, "Name")
	if err != nil {
		return nil, err
	}

	app.Type, err = cfg.GetRawString(DE, "Type")
	if err != nil {
		return nil, err
	}

	app.NoDisplay, err = cfg.GetBool(DE, "NoDisplay")
	if err != nil {
		app.NoDisplay = false
	}

	app.cfg = cfg
	app.xdg = xdg

	return app, nil
}

func (app *Application) Run() error {
	tmp := strings.Split(app.Exec, " ")

	the_cmd := exec.Command("nohup", tmp...)

	return the_cmd.Run()
}

func (app *Application) FindIcon(size int) image.Image {
	return app.xdg.GetIcon(app.Icon, size)
}
