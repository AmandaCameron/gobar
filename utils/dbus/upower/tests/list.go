package main

import (
	"fmt"

	"launchpad.net/~jamesh/go-dbus/trunk"

	"dbus/upower"
)

func failMeMaybe(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	sys, err := dbus.Connect(dbus.SystemBus)
	failMeMaybe(err)

	failMeMaybe(sys.Authenticate())

	up := upower.New(sys)

	devs, err := up.GetDevices()
	failMeMaybe(err)

	for i, dev := range devs {
		charge, err := dev.Charge()
		failMeMaybe(err)

		state, err := dev.State()
		failMeMaybe(err)

		fmt.Printf("-- Device %d --\n", i)
		fmt.Printf("  Type: %s\n", dev.Type())
		fmt.Printf("  State: %s\n", state)
		fmt.Printf("  Percentage: %f\n", charge)
		fmt.Printf("\n")
	}
}
