package nm

import (
	"fmt"
)

func (dt DeviceType) String() string {
	if dt == Unknown {
		return "Unknown"
	} else if dt == Wired {
		return "Wired"
	} else if dt == Wireless {
		return "Wireless"
	}
	return fmt.Sprintf("Unknown (%d)", dt)
}

func (dev *Device) devTypeGet(name string) (interface{}, error) {
	iface := ""
	if dev.typ == Wired {
		iface = NM_DEV_WIRED_IFACE
	} else if dev.typ == Wireless {
		iface = NM_DEV_WIRELESS_IFACE
	} else {
		return nil, fmt.Errorf("Unknown device type: %s", dev.typ)
	}

	return dev.Get(iface, "HwAddress")
}

func (dev *Device) mustBe(typ DeviceType) error {
	if dev.typ == typ {
		return nil
	}
	return fmt.Errorf("This function only applys to %s devices, this is a %s", typ, dev.typ)
}
