package nm

func (dev *Device) Mac() (string, error) {
	tmp, err := dev.devTypeGet("HwAddress")
	if err != nil {
		return "", err
	}

	return tmp.(string), nil
}

func (dev *Device) PermMac() (string, error) {
	tmp, err := dev.devTypeGet("PermHwAddress")
	if err != nil {
		return "", err
	}

	return tmp.(string), nil
}
