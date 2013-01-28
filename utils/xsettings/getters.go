package xsettings

import (
	"fmt"
)

// Gets a speffic XSetting out of the parsed list, if it can not be found,
// it will return nil
func (xs *XSettings) Get(key string) *XSetting {
	set, ok := xs.props[key]
	if ok {
		return set
	}
	return nil
}

// Gets all the XSetting values.
func (xs *XSettings) GetAll() map[string]*XSetting {
	ret := make(map[string]*XSetting)

	for k, val := range xs.props {
		ret[k] = val
	}

	return ret
}

// Gets a string by the name of `key` -- if the key does not exist, it will
// return "" and an error of "No Such Key {key-name}"
func (xs *XSettings) GetString(key string) (string, error) {
	setting := xs.Get(key)
	if setting == nil {
		return "", fmt.Errorf("No such key: %s", key)
	}

	if err := setting.mustBe(XSettingString); err != nil {
		return "", err
	}

	return setting.String, nil
}

// Gets an integer by the name of `key` -- if the key does not exist, it will
// return -1 and an error of "No such key: {key-name}"
func (xs *XSettings) GetInteger(key string) (int32, error) {
	setting := xs.Get(key)
	if setting == nil {
		return -1, fmt.Errorf("No such key: %s", key)
	}

	if err := setting.mustBe(XSettingInteger); err != nil {
		return -1, err
	}

	return setting.Integer, nil
}

// TODO: Get Colours, maybe?
