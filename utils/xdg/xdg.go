package xdg

import (
	//"os"

	"image"
)

type XDG struct {
	theme       string
	iconTheme   []iconTheme
	config_path string

	apps  map[string]*Application
	icons map[Icon]image.Image
}

func New() *XDG {
	xdg := &XDG{}

	xdg.initApps()

	return xdg
}

func (xdg *XDG) SetTheme(name string) error {
	xdg.theme = name

	return xdg.initIcons()
}

// func (xdg *XDG) GetIcons() Icons {
// 	return xdg.icons
// }

func (xdg *XDG) GetApps() []*Application {
	apps := make([]*Application, 0)
	for _, app := range xdg.apps {
		apps = append(apps, app)
	}
	return apps
}

func (xdg *XDG) GetIcon(name string, size int) image.Image {
	if xdg.theme == "" {
		return nil
	}

	img := xdg.lookupIcon(name, size)

	return img
}
