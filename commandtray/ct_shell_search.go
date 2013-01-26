package commandtray

import (
	"image"
	"io/ioutil"
	"os"
	"strings"

	"code.google.com/p/goconf/conf"
	"launchpad.net/~jamesh/go-dbus/trunk"

	"github.com/AmandaCameron/gobar/utils/dbus/gs_search"
	"github.com/AmandaCameron/gobar/utils/xdg"
)

type ShellSource struct {
	sess_conn *dbus.Connection
	Xdg       *xdg.XDG

	searchers []*gs_search.GsSearch
}

type ShellCommand struct {
	search *gs_search.GsSearch
	id     string
	meta   map[string]dbus.Variant
	terms  []string
	Xdg    *xdg.XDG
}

func NewShellSource(sess *dbus.Connection, x *xdg.XDG) *ShellSource {
	ss := &ShellSource{
		sess_conn: sess,
		Xdg:       x,
	}
	for _, dir := range []string{
		"/usr/share/gnome-shell/search-providers",
		"/usr/local/share/gnoem-shell/search-providers",
	} {

		srcs, err := ioutil.ReadDir(dir)
		//failMeMaybe(err)
		if err != nil {
			continue
		}
		for _, file := range srcs {
			cfg, err := conf.ReadConfigFile(dir + "/" + file.Name())
			failMeMaybe(err)

			SSP := "Shell Search Provider"

			objPath, err := cfg.GetString(SSP, "ObjectPath")
			if err != nil {
				continue
			}

			busName, err := cfg.GetString(SSP, "BusName")
			if err != nil {
				continue
			}

			var name, icon string

			name, err = cfg.GetString(SSP, "Name")
			if err != nil {
				did, err := cfg.GetString(SSP, "DesktopId")
				if err == nil {
					name, icon = getName(did)
				}
			}

			if icon == "" {
				if tmp, err := cfg.GetString(SSP, "Icon"); err == nil {
					icon = tmp
				}
			}

			searcher := gs_search.New(sess.Object(busName, dbus.ObjectPath(objPath)))
			searcher.Name = name
			searcher.Icon = icon

			ss.searchers = append(ss.searchers, searcher)
		}
	}

	return ss
}

func getName(did string) (string, string) {
	for _, dir := range []string{
		"/usr/share/applications",
		"/usr/local/share/applications",
		os.Getenv("HOME") + "/.local/share/applications",
	} {
		cfg, err := conf.ReadConfigFile(dir + "/" + did)
		if err == nil {
			name, err := cfg.GetString("Desktop Entry", "Name")
			if err != nil {
				return strings.Split(did, ".")[0], ""
			}

			icon, err := cfg.GetString("Desktop Entry", "Icon")
			if err != nil {
				return strings.Split(did, ".")[0], ""
			}

			return name, icon
		}
	}

	return strings.Split(did, ".")[0], ""
}

func (ss *ShellSource) GetMatches(inp string, ct *CommandTray) []Command {
	if ss.sess_conn == nil {
		return []Command{}
	}

	if !strings.HasPrefix(inp, "gs ") {
		return []Command{}
	}

	inp = inp[3:]

	terms := strings.Split(inp, " ")

	ret := make([]Command, 0)

	for _, searcher := range ss.searchers {
		results, err := searcher.GetResults(terms)
		if err != nil {
			continue
		}

		meta, err := searcher.GetResultsMeta()
		if err != nil {
			continue
		}

		for i := range results {
			ret = append(ret, &ShellCommand{
				search: searcher,
				id:     results[i],
				meta:   meta[i],
				Xdg:    ss.Xdg,
				terms:  terms,
			})
		}
	}

	return ret
}

func (ss *ShellSource) Open(ct *CommandTray) bool {
	return true
}

func (ss *ShellSource) Close(ct *CommandTray) {
	// Do Nothing.
}

func (sc *ShellCommand) GetText() string {
	tmp, ok := sc.meta["name"]
	if !ok {
		return "[" + sc.search.Name + "] ** Error **"
	}

	return "[" + sc.search.Name + "] " + tmp.Value.(string)
}

func (sc *ShellCommand) GetIcon() image.Image {
	tmp, ok := sc.meta["gicon"]
	if !ok {
		return sc.Xdg.GetIcon(sc.search.Icon, 16)
	}

	icon := tmp.Value.(string)

	if icon[0] == '.' {
		return sc.Xdg.GetIcon(sc.search.Icon, 16)
	}

	return sc.Xdg.GetIcon(icon, 16)
}

func (sc *ShellCommand) Run() {
	sc.search.Activate(sc.id, sc.terms)
}
