package main

import (
	"fmt"

	"launchpad.net/~jamesh/go-dbus/trunk"

	"dbus/nm"
)

func failMeMaybe(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	conn, err := dbus.Connect(dbus.SystemBus)
	failMeMaybe(err)

	failMeMaybe(conn.Authenticate())

	cli := nm.New(conn)

	devs, err := cli.GetDevices()
	failMeMaybe(err)

	for _, dev := range devs {
		ipiface, err := dev.Get(nm.NM_DEV_IFACE, "Interface")
		failMeMaybe(err)

		driver, err := dev.Driver()
		failMeMaybe(err)

		driver_ver, err := dev.DriverVersion()
		failMeMaybe(err)

		firmware_version, err := dev.FirmwareVersion()
		failMeMaybe(err)

		mac, err := dev.Mac()
		failMeMaybe(err)

		typ := dev.Type()

		fmt.Printf("---- Device %s ----\n", ipiface)
		fmt.Printf("  Type:     %s\n", typ)
		fmt.Printf("  Driver:   %s (Version: %s)\n", driver, driver_ver)
		fmt.Printf("  Firmware: %s\n", firmware_version)
		fmt.Printf("  MAC:      %s\n", mac)
		if typ == nm.Wireless {
			active, err := dev.GetActive()
			failMeMaybe(err)

			active_mac, err := active.Mac()
			failMeMaybe(err)

			fmt.Printf("  Access Points:\n")
			aps, err := dev.GetAccessPoints()

			failMeMaybe(err)
			for _, ap := range aps {
				ssid, err := ap.Name()
				failMeMaybe(err)

				strength, err := ap.Strength()
				failMeMaybe(err)

				mac, err := ap.Mac()
				failMeMaybe(err)

				maybeActive := ""
				if mac == active_mac {
					maybeActive = "(Active)"
				}

				fmt.Printf("    %s(%s) -- %d%% %s\n", ssid, mac, strength, maybeActive)
				fmt.Printf("       Path: %s\n", ap.GetObjectPath())
			}
		}
		fmt.Printf("\n")
	}

}
