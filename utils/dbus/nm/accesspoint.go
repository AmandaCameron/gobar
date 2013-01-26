package nm

import (
	// "fmt"

	"launchpad.net/~jamesh/go-dbus/trunk"
)

type AccessPoint struct {
	*dbus.ObjectProxy
	*dbus.Properties

	cli *Client
	dev *Device
}

const APSecured = 0x01

func (cli *Client) newAccessPoint(device *Device, path dbus.ObjectPath) *AccessPoint {
	if path == dbus.ObjectPath("/") {
		// fmt.Printf("Access Point is / -- Returning Nil.")
		return nil
	}

	obj := cli.conn.Object(NM_UNIQ_NAME, path)

	return &AccessPoint{
		ObjectProxy: obj,
		Properties:  &dbus.Properties{obj},
		cli:         cli,
		dev:         device,
	}
}

func (ap *AccessPoint) Name() (string, error) {
	ret, err := ap.Get(NM_AP_IFACE, "Ssid")
	if err != nil {
		return "", err
	}

	tmp := ret.([]interface{})

	buff := make([]byte, 0, len(tmp))

	for _, c := range tmp {
		buff = append(buff, c.(byte))
	}

	return string(buff), nil
}

func (ap *AccessPoint) Strength() (byte, error) {
	ret, err := ap.Get(NM_AP_IFACE, "Strength")
	if err != nil {
		return 255, err
	}

	return ret.(byte), nil
}

func (ap *AccessPoint) Mac() (string, error) {
	ret, err := ap.Get(NM_AP_IFACE, "HwAddress")
	if err != nil {
		return "", err
	}

	return ret.(string), nil
}

func (ap *AccessPoint) Flags() (uint32, error) {
	ret, err := ap.Get(NM_AP_IFACE, "Flags")
	if err != nil {
		return 0, err
	}

	return ret.(uint32), err
}

func (ap *AccessPoint) WpaFlags() (uint32, error) {
	ret, err := ap.Get(NM_AP_IFACE, "WpaFlags")
	if err != nil {
		return 0, err
	}

	return ret.(uint32), err
}

func (ap *AccessPoint) RsnFlags() (uint32, error) {
	ret, err := ap.Get(NM_AP_IFACE, "RsnFlags")
	if err != nil {
		return 0, err
	}

	return ret.(uint32), err
}

// func (ap *AccessPoint) Connect(extras map[string]map[string]interface{}) (ret dbus.ObjectPath, err error) {
// 	msg, err := ap.cli.Call(NM_BASE_IFACE, "AddAndActivateConnection", extras, ap.dev, ap.GetObjectPath())
// 	if err != nil {
// 		return
// 	}

// 	var path dbus.ObjectPath
// 	err = msg.GetArgs(&path, &ret)

// 	return
// }
