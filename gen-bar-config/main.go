package main

import (
	"fmt"
	"os"
	"text/template"

	xgbxinerama "github.com/BurntSushi/xgb/xinerama"
	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/xinerama"

	"github.com/AmandaCameron/gobar/utils"
)

var templ *template.Template

func init() {
	templ = template.New("Config File")
	templ.Parse(`
[Main]
Width := {{.Width}}
Size := 24
Position := Top

[Clock]
Width := 200
Position := {{.ClockPos}}
# Additional Keys:
# See http://golang.org/pkg/time/#pkg-constants for more examples.
# Format := "Monday, Jan 02, 2006"

[CommandTray]
Width := {{.CommandWidth}}
Position := 0
Accel := Mod4-b

[App Tracker]
Width := 0
Position := {{.CommandWidth}}

[StatusBar]
Width := 200
Position := {{.StatusPos}}

[Fonts/Clock]
Name := {{.ClockFont}}
Size := 12

[Fonts/CommandTray]
Name := {{.CommandFont}}
Size := 12
`[1:])
}

type Config struct {
	Width        int
	ClockPos     int
	CommandWidth int
	StatusPos    int
	ClockFont    string
	CommandFont  string
}

func main() {
	cfg := Config{}

	X, err := xgbutil.NewConn()
	utils.FailMeMaybe(err)

	xgbXineramaIsActiveCookieReply, err := xgbxinerama.IsActive(X.Conn()).Reply()
	utils.FailMeMaybe(err)

	if xgbXineramaIsActiveCookieReply.State == 0 {
		cfg.Width = int(X.Setup().DefaultScreen(X.Conn()).WidthInPixels)
	} else {
		heads, err := xinerama.PhysicalHeads(X)
		utils.FailMeMaybe(err)

		cfg.Width = heads[0].Width()
	}

	cfg.ClockPos = (cfg.Width / 2) - 100
	cfg.CommandWidth = (cfg.Width / 2) - 124
	cfg.StatusPos = (cfg.Width / 2) + 100

	// Look for ze fonts!

	cfg.ClockFont = "Put your font here!"
	cfg.CommandFont = cfg.ClockFont
	found := false

	for _, fontDir := range []string{"/usr/share/fonts/dejavu/", "/usr/share/fonts/TTF/", "/usr/local/lib/X11/fonts/dejavu/"} {
		if exists(fontDir+"DejaVuSansMono-Bold.ttf") && exists(fontDir+"DejaVuSansMono.ttf") {
			cfg.ClockFont = fontDir + "DejaVuSansMono-Bold.ttf"
			cfg.CommandFont = fontDir + "DejaVuSansMono.ttf"
			found = true
			break
		}
	}

	if !found {
		fmt.Fprintln(os.Stderr, "# We could not locate your fonts directory, please edit the config\n")
		fmt.Fprintln(os.Stderr, "# Accordingly.\n")
	}

	templ.Execute(os.Stdout, cfg)
}

func exists(fileName string) bool {
	f, err := os.Open(fileName)
	if err != nil {
		return false
	}
	f.Close()
	return true
}
