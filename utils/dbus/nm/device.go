package nm

import (
	"fmt"

	"launchpad.net/~jamesh/go-dbus/trunk"
)

type Device struct {
	*dbus.ObjectProxy
	*dbus.Properties

	cli *Client
	typ DeviceType
}

type DeviceType uint32

type DeviceHandler func(*Device)

const (
	Unknown DeviceType = iota
	Wired
	Wireless
)

func (cli *Client) newDevice(path dbus.ObjectPath) *Device {
	obj := cli.conn.Object(NM_UNIQ_NAME, path)

	if obj == nil {
		return nil
	}

	dev := &Device{
		ObjectProxy: obj,
		Properties:  &dbus.Properties{obj},
		cli:         cli,
	}

	ret, err := dev.Get(NM_DEV_IFACE, "DeviceType")
	if err == nil {
		tmp := ret.(uint32)

		dev.typ = DeviceType(tmp)
	} else {
		dev.typ = Unknown
	}

	return dev
}

func (dev *Device) Disconnect() error {
	_, err := dev.Call(NM_DEV_IFACE, "Disconnect")
	if err != nil {
		return err
	}

	return nil
}

func (dev *Device) Type() DeviceType {
	return dev.typ
}

func (dev *Device) Driver() (string, error) {
	ret, err := dev.Get(NM_DEV_IFACE, "Driver")
	if err != nil {
		return "", err
	}

	return ret.(string), nil
}

func (dev *Device) DriverVersion() (string, error) {
	ret, err := dev.Get(NM_DEV_IFACE, "DriverVersion")
	if err != nil {
		return "", err
	}

	return ret.(string), nil
}

func (dev *Device) FirmwareVersion() (string, error) {
	ret, err := dev.Get(NM_DEV_IFACE, "FirmwareVersion")
	if err != nil {
		return "", err
	}

	return ret.(string), nil
}

func (dev *Device) State() (uint32, error) {
	ret, err := dev.Get(NM_DEV_IFACE, "State")
	if err != nil {
		return 0, err
	}

	return ret.(uint32), err
}

func (dev *Device) PropChanged(handler DeviceHandler) error {
	iface := ""

	if dev.typ == Wired {
		iface = NM_DEV_WIRED_IFACE
	} else if dev.typ == Wireless {
		iface = NM_DEV_WIRELESS_IFACE
	} else {
		return fmt.Errorf("Unknown device type: %s", dev.typ)
	}

	_, err := dev.WatchSignal(iface, "PropertiesChanged", func(_ *dbus.Message) {
		go handler(dev)
	})

	return err
}
