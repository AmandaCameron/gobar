package upower

import (
	"fmt"
)

func (ds DeviceState) String() string {
	if ds == UnknownState {
		return "Unknown"
	} else if ds == Charging {
		return "Charging"
	} else if ds == Discharging {
		return "Discharging"
	} else if ds == Empty {
		return "Empty"
	} else if ds == Full {
		return "Full"
	} else if ds == PendingCharge {
		return "PendingCharge"
	} else if ds == PendingDischarge {
		return "PendingDischarge"
	}
	return fmt.Sprintf("Unknown (%d)", ds)
}

func (dt DeviceType) String() string {
	if dt == Unknown {
		return "Unknown"
	} else if dt == LinePower {
		return "Line Power"
	} else if dt == Battery {
		return "Battery"
	}
	return fmt.Sprintf("Unknown (%d)", dt)
}
