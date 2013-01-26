package upower

import (
	"errors"

	"launchpad.net/~jamesh/go-dbus/trunk"
)

type UPower struct {
	*dbus.ObjectProxy
	*dbus.Properties

	conn *dbus.Connection
}

const (
	UP_IFACE     = "org.freedesktop.UPower"
	UP_DEV_IFACE = "org.freedesktop.UPower.Device"
	UP_NAME      = "org.freedesktop.UPower"
)

func New(conn *dbus.Connection) *UPower {
	obj := conn.Object(UP_NAME, "/org/freedesktop/UPower")

	return &UPower{
		ObjectProxy: obj,
		Properties:  &dbus.Properties{obj},
		conn:        conn,
	}
}

func (up *UPower) GetDevices() ([]*Device, error) {
	msg, err := up.Call(UP_IFACE, "EnumerateDevices")
	if err != nil {
		return nil, err
	}

	var tmp []dbus.ObjectPath

	if err = msg.GetArgs(&tmp); err != nil {
		return nil, err
	}

	var ret []*Device

	for _, path := range tmp {
		ret = append(ret, up.newDevice(path))
	}

	return ret, nil
}

func (up *UPower) GetBattery() (*Device, error) {
	devs, err := up.GetDevices()
	if err != nil {
		return nil, err
	}

	for _, dev := range devs {
		if dev.Type() == Battery {
			return dev, nil
		}
	}

	return nil, errors.New("No battery in machine.")
}
