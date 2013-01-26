package xdg

import (
	"fmt"
	"image"
	"image/png"
	//"io/ioutil"
	"os"
	"strings"

	"code.google.com/p/goconf/conf"
	"code.google.com/p/graphics-go/graphics"
)

var XDG_THEME_DIRS = []string{
	"/usr/share/icons",
	"/usr/local/share/icons",
}

type Icon struct {
	name string
	size int
}

type iconDirType uint8

const (
	iconFixed iconDirType = iota
	iconScalable
	iconThreshold
)

type iconTheme struct {
	name string
	desc string
	dirs []iconDir
}

type iconDir struct {
	iconType    iconDirType
	iconSize    int
	iconMinSize int
	iconMaxSize int
	dirName     string
}

func (dt iconDirType) String() string {
	if dt == iconFixed {
		return "Fixed"
	} else if dt == iconScalable {
		return "Scalable"
	} else if dt == iconThreshold {
		return "Threshold"
	}
	return fmt.Sprintf("Unknown (%d)", dt)
}

func (xdg *XDG) initIcons() error {
	xdg.icons = make(map[Icon]image.Image)

	xdg.iconTheme = nil

	if err := xdg.readIndex(xdg.theme); err != nil {
		return err
	}

	if err := xdg.readIndex("hicolor"); err != nil {
		return err
	}

	return nil
}

func (xdg *XDG) readIndex(theme string) error {
	foundValid := false

	for _, dir := range XDG_THEME_DIRS {
		// fmt.Printf("Reading %s/%s\n", dir, theme)

		cfg, err := conf.ReadConfigFile(dir + "/" + theme + "/index.theme")
		if err != nil {
			// fmt.Printf("%s does not have the file we are looking for.\n", dir)
			continue
		}

		foundValid = true

		var it iconTheme

		if it.name, err = cfg.GetString("Icon Theme", "Name"); err != nil {
			// fmt.Printf("Invalid Theme ( Missing Name ) %s\n", theme)
			return fmt.Errorf("Invalid theme %s (Missing: Name)", theme)
		}

		if it.desc, err = cfg.GetString("Icon Theme", "Comment"); err != nil {
			// fmt.Printf("Invalid Theme ( Missing Comment ) %s\n", theme)
			return fmt.Errorf("Invalid theme %s (Missing: Comment)", theme)
		}

		var dirs string

		if dirs, err = cfg.GetString("Icon Theme", "Directories"); err != nil {
			//fmt.Printf("Invalid theme ( Missing Directories ) %s\n", theme)
			return fmt.Errorf("Invalid theme %s (Missing: Directories)", theme)
		}

		for _, d := range strings.Split(dirs, ",") {
			var id iconDir

			var tmp string

			if id.iconSize, err = cfg.GetInt(d, "Size"); err != nil {
				return fmt.Errorf("Invalid theme %s (Missing: Size)", theme)
			}

			if tmp, err = cfg.GetString(d, "Type"); err != nil {
				id.iconType = iconThreshold
			} else {
				if tmp == "Threshold" {
					id.iconType = iconThreshold
				} else if tmp == "Scalable" {
					id.iconType = iconScalable
				} else if tmp == "Fixed" {
					id.iconType = iconFixed
				} else {
					return fmt.Errorf("Invalid theme %s (Invalid IT: %s)", theme, tmp)
				}
			}

			if id.iconType == iconThreshold {
				var tmpInt int

				if tmpInt, err = cfg.GetInt(d, "Threshold"); err != nil {
					tmpInt = 2
				}

				id.iconMinSize = id.iconSize - tmpInt
				id.iconMaxSize = id.iconSize + tmpInt
			} else if id.iconType == iconScalable {
				if id.iconMinSize, err = cfg.GetInt(d, "MinSize"); err != nil {
					id.iconMinSize = id.iconSize
				}
				if id.iconMaxSize, err = cfg.GetInt(d, "MaxSize"); err != nil {
					id.iconMaxSize = id.iconSize
				}
			} else if id.iconType == iconFixed {
				id.iconMinSize = id.iconSize
				id.iconMaxSize = id.iconSize
			}

			id.dirName = dir + "/" + theme + "/" + d

			it.dirs = append(it.dirs, id)
		}

		xdg.iconTheme = append(xdg.iconTheme, it)
	}

	if !foundValid {
		return fmt.Errorf("No such theme: %s", theme)
	}

	return nil
}

func (xdg *XDG) scaleIcon(srcImg image.Image, name string, size int) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, size, size))

	err := graphics.Scale(img, srcImg)
	if err != nil {
		fmt.Printf("Failed to scale: %s\n", err.Error())
		return nil
	}

	xdg.icons[Icon{name, size}] = img

	return img
}

func (xdg *XDG) lookupIcon(name string, size int) image.Image {
	// Cache Check!

	img, ok := xdg.icons[Icon{name, size}]
	if ok {
		return img
	}

	//fmt.Printf("Looking for %s  (%dx%d)\n", name, size, size)

	var found []iconDir

	for _, it := range xdg.iconTheme {
		for _, dir := range it.dirs {
			if dir.iconSize == size {
				f, err := os.Open(dir.dirName + "/" + name + ".png")
				if err != nil {
					continue
				}
				defer f.Close()

				img, err := png.Decode(f)
				if err != nil {
					continue
				}

				xdg.icons[Icon{name, size}] = img
				return img
			} else if dir.iconMinSize <= size && dir.iconMaxSize >= size {
				found = append(found, dir)
			}
		}
	}

	for _, dir := range found {
		// fmt.Printf(" -- Checking %s -- Type: %s -- Min/Max Size: %d/%d -- Size: %d\n",
		// 	dir.dirName, dir.iconType, dir.iconMinSize, dir.iconMaxSize, dir.iconSize)
		f, err := os.Open(dir.dirName + "/" + name + ".png")
		if err != nil {
			continue
		}
		defer f.Close()
		img, err := png.Decode(f)
		if err != nil {
			continue
		}

		return xdg.scaleIcon(img, name, size)
	}

	// Last Resort:

	// fmt.Printf("Last Resort!\n")

	f, err := os.Open("/usr/share/pixmaps/" + name + ".png")
	if err != nil {
		return nil
	}
	defer f.Close()

	img, err = png.Decode(f)
	if err != nil {
		return nil
	}

	return xdg.scaleIcon(img, name, size)
}
