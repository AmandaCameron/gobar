package nm

import (
	"launchpad.net/~jamesh/go-dbus/trunk"
)

func (dev *Device) RequestScan(opts map[string]interface{}) error {
	if err := dev.mustBe(Wireless); err != nil {
		return err
	}

	_, err := dev.Call(NM_DEV_WIRELESS_IFACE, "RequestScan", opts)
	if err != nil {
		return err
	}

	return nil
}

func (dev *Device) GetAccessPoints() ([]*AccessPoint, error) {
	if err := dev.mustBe(Wireless); err != nil {
		return nil, err
	}

	msg, err := dev.Call(NM_DEV_WIRELESS_IFACE, "GetAccessPoints")
	if err != nil {
		return nil, err
	}

	var paths []dbus.ObjectPath

	if err = msg.GetArgs(&paths); err != nil {
		return nil, err
	}

	ret := make([]*AccessPoint, 0, len(paths))

	for _, path := range paths {
		ret = append(ret, dev.cli.newAccessPoint(dev, path))
	}

	return ret, nil
}

func (dev *Device) GetActive() (*AccessPoint, error) {
	if err := dev.mustBe(Wireless); err != nil {
		return nil, err
	}

	tmp, err := dev.Get(NM_DEV_WIRELESS_IFACE, "ActiveAccessPoint")
	if err != nil {
		return nil, err
	}

	return dev.cli.newAccessPoint(dev, tmp.(dbus.ObjectPath)), nil
}
