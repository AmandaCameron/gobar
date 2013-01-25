package main

import (
	"launchpad.net/~jamesh/go-dbus/trunk"
)

// Menus

type GtkMenu struct {
	*dbus.ObjectProxy
}

type GtkMenuOptions struct {
	GroupId uint32
	MenuId  uint32
	Options []map[string]dbus.Variant
}

func (gm *GtkMenu) Start(menus []uint32) (ret []GtkMenuOptions, err error) {
	msg, err := gm.Call("org.gtk.Menus", "Start", menus)
	if err != nil {
		return
	}

	if err = msg.GetArgs(&ret); err != nil {
		return
	}

	return
}

func (gm *GtkMenu) End(menus []uint32) (err error) {
	_, err = gm.Call("org.gtk.Menus", "End", menus)
	return
}

// Actions

type GtkActions struct {
	*dbus.ObjectProxy
}

type ActionDesc struct {
	Enabled   bool
	Signature dbus.Signature
	UnkArray  []dbus.Variant
}

func (ga *GtkActions) Activate(act string, args []interface{}, platform map[string]interface{}) error {
	var var_args []dbus.Variant
	var var_platform map[string]dbus.Variant

	for _, arg := range args {
		var_args = append(var_args, dbus.Variant{arg})
	}

	for key, value := range platform {
		var_platform[key] = dbus.Variant{value}
	}

	_, err := ga.Call("org.gtk.Actions", "Activate", act, var_args, var_platform)
	return err
}

func (ga *GtkActions) List() (ret []string, err error) {
	msg, err := ga.Call("org.gtk.Actions", "List")
	if err != nil {
		return
	}

	err = msg.GetArgs(&ret)

	return
}

func (ga *GtkActions) Describe(act string) (desc ActionDesc, err error) {
	msg, err := ga.Call("org.gtk.Actions", "Describe", act)
	if err != nil {
		return
	}

	err = msg.GetArgs(&desc)

	return
}
