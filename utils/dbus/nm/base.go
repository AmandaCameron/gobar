package nm

import (
	"launchpad.net/~jamesh/go-dbus/trunk"
)

type Client struct {
	*dbus.ObjectProxy
	*dbus.Properties

	conn *dbus.Connection
}

const (
	NM_UNIQ_NAME = "org.freedesktop.NetworkManager"

	NM_BASE_IFACE         = "org.freedesktop.NetworkManager"
	NM_DEV_IFACE          = "org.freedesktop.NetworkManager.Device"
	NM_DEV_WIRELESS_IFACE = "org.freedesktop.NetworkManager.Device.Wireless"
	NM_DEV_WIRED_IFACE    = "org.freedesktop.NetworkManager.Device.Wired"
	NM_AP_IFACE           = "org.freedesktop.NetworkManager.AccessPoint"
)

func New(conn *dbus.Connection) *Client {
	obj := conn.Object("org.freedesktop.NetworkManager", "/org/freedesktop/NetworkManager")

	cli := &Client{
		ObjectProxy: obj,
		Properties:  &dbus.Properties{obj},
		conn:        conn,
	}

	return cli
}

func (cli *Client) GetDevices() ([]*Device, error) {
	msg, err := cli.Call(NM_BASE_IFACE, "GetDevices")
	if err != nil {
		return nil, err
	}

	var paths []dbus.ObjectPath

	if err = msg.GetArgs(&paths); err != nil {
		return nil, err
	}

	ret := make([]*Device, 0, len(paths))

	for _, path := range paths {
		ret = append(ret, cli.newDevice(path))
	}

	return ret, nil
}

func (cli *Client) GetDeviceByIpIface(iface string) (*Device, error) {
	msg, err := cli.Call(NM_BASE_IFACE, "GetDeviceByIpIface", iface)
	if err != nil {
		return nil, err
	}

	var path dbus.ObjectPath
	if err = msg.GetArgs(&path); err != nil {
		return nil, err
	}

	return cli.newDevice(path), nil
}
