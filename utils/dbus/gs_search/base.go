package gs_search

import (
	"fmt"

	"launchpad.net/~jamesh/go-dbus/trunk"
)

const IFACE_V1 = "org.gnome.Shell.SearchProvider"
const IFACE_V2 = "org.gnome.Shell.SearchProvider2" // Currently unused.

type GsSearch struct {
	*dbus.ObjectProxy

	// Meta Data ( Not retreived by us. )
	Name string
	Icon string

	prev_results []string
	iface        string
}

func New(prox *dbus.ObjectProxy) *GsSearch {
	// TODO: Make this detect the interfaces it has, and use the proper one.

	return &GsSearch{
		ObjectProxy: prox,
		iface:       IFACE_V1,
	}
}

func (gs *GsSearch) GetResults(terms []string) ([]string, error) {
	var msg *dbus.Message
	var err error

	if gs.prev_results == nil {
		msg, err = gs.Call(gs.iface, "GetInitialResultSet", terms)
	} else {
		msg, err = gs.Call(gs.iface, "GetSubsearchResultSet", gs.prev_results, terms)
	}

	if err != nil {
		return nil, err
	}

	if err = msg.GetArgs(&gs.prev_results); err != nil {
		return nil, err
	}

	return gs.prev_results, nil
}

func (gs *GsSearch) GetResultsMeta() ([]map[string]dbus.Variant, error) {
	if gs.prev_results == nil {
		return nil, fmt.Errorf("No results")
	}

	msg, err := gs.Call(gs.iface, "GetResultMetas", gs.prev_results)
	if err != nil {
		return nil, err
	}

	var ret []map[string]dbus.Variant

	if err = msg.GetArgs(&ret); err != nil {
		return nil, err
	}

	return ret, nil
}

func (gs *GsSearch) Activate(id string, terms []string) error {
	_, err := gs.Call(gs.iface, "ActivateResult", id)
	return err
}

func (gs *GsSearch) Reset() {
	gs.prev_results = nil
}
