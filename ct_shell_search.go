package main

import (
	"image"
	"io/ioutil"
	"strings"

	"code.google.com/p/goconf/conf"
	"launchpad.net/~jamesh/go-dbus/trunk"

	"dbus/gs_search"
	"xdg"
)

type ShellSource struct {
	sess_conn *dbus.Connection
	Xdg       *xdg.XDG

	searchers []*gs_search.GsSearch
}

type ShellCommand struct {
	search *gs_search.GsSearch

	id   string
	meta map[string]dbus.Variant
	Xdg  *xdg.XDG
}

func NewShellSource(sess *dbus.Connection) *ShellSource {
	ss := &ShellSource{
		sess_conn: sess,
		Xdg:       xdg.New(),
	}
	srcs, err := ioutil.ReadDir("/usr/share/gnome-shell/search-providers")
	failMeMaybe(err)
	for _, file := range srcs {
		cfg, err := conf.ReadConfigFile("/usr/share/gnome-shell/search-providers/" + file.Name())
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

		ss.searchers = append(ss.searchers, gs_search.New(sess.Object(busName, dbus.ObjectPath(objPath))))
	}

	return ss
}

func (ss *ShellSource) GetMatches(inp string, ct *CommandTray) []Command {
	if ss.sess_conn == nil {
		return []Command{}
	}

	ret := make([]Command, 0)

	for _, searcher := range ss.searchers {
		results, err := searcher.GetResults(strings.Split(inp, " "))
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
			})
		}
	}

	return ret
}

func (sc *ShellCommand) GetText() string {
	tmp, ok := sc.meta["name"]
	if !ok {
		return "** Error **"
	}

	return tmp.Value.(string)
}

func (sc *ShellCommand) GetIcon() image.Image {
	tmp, ok := sc.meta["gicon"]
	if !ok {
		return nil
	}

	icon := tmp.Value.(string)

	if icon[0] == '.' {
		return nil
	}

	return sc.Xdg.GetIcon(icon, 16)
}

func (sc *ShellCommand) Run() {
	sc.search.Activate(sc.id, []string{})
}
