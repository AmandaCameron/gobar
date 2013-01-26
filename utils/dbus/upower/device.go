package upower

import (
	"launchpad.net/~jamesh/go-dbus/trunk"
)

type Device struct {
	*dbus.ObjectProxy
	*dbus.Properties

	up  *UPower
	typ DeviceType
}

type DeviceType uint32

const (
	Unknown DeviceType = iota
	LinePower
	Battery
)

type DeviceState uint32

const (
	UnknownState DeviceState = iota
	Charging
	Discharging
	Empty
	Full
	PendingCharge
	PendingDischarge
)

func (up *UPower) newDevice(path dbus.ObjectPath) *Device {
	obj := up.conn.Object(UP_NAME, path)
	props := &dbus.Properties{obj}

	tmp, err := props.Get(UP_DEV_IFACE, "Type")
	if err != nil {
		return nil
	}

	return &Device{
		ObjectProxy: obj,
		Properties:  &dbus.Properties{obj},
		up:          up,
		typ:         DeviceType(tmp.(uint32)),
	}
}

func (dev *Device) Connect(handler func(*Device)) {
	dev.WatchSignal(UP_DEV_IFACE, "Changed", func(_ *dbus.Message) {
		go handler(dev)
	})
}

func (dev *Device) Type() DeviceType {
	return dev.typ
}

func (dev *Device) State() (DeviceState, error) {
	tmp, err := dev.Get(UP_DEV_IFACE, "State")
	if err != nil {
		return UnknownState, err
	}

	return DeviceState(tmp.(uint32)), nil
}

func (dev *Device) Charge() (float64, error) {
	tmp, err := dev.Get(UP_DEV_IFACE, "Percentage")
	if err != nil {
		return -1, err
	}

	return tmp.(float64), nil
}
