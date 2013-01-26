package main

import (
	"fmt"
	"os"

	"launchpad.net/~jamesh/go-dbus/trunk"

	"dbus/gs_search"
)

func failMeMaybe(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	sess, err := dbus.Connect(dbus.SessionBus)
	failMeMaybe(err)

	sess.Authenticate()

	obj := sess.Object("org.gnome.Nautilus.SearchProvider", "/org/gnome/Nautilus/SearchProvider")

	searcher := gs_search.New(obj)

	res, err := searcher.GetResults(os.Args[1:])
	failMeMaybe(err)

	meta, err := searcher.GetResultsMeta()
	failMeMaybe(err)

	for i, id := range res {
		fmt.Printf("Result %d:\n", i+1)
		fmt.Printf("  Id: %s\n", id)
		fmt.Printf("  Meta:\n")
		for k, m := range meta[i] {
			sig, _ := m.GetVariantSignature()
			fmt.Printf("    %s(%s): %v\n", k, sig, m)
		}
		fmt.Printf("\n")
	}

}
