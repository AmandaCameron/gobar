package xsettings

import (
	"fmt"
)

func (typ XSettingType) String() string {
	if typ == XSettingInteger {
		return "Integer"
	} else if typ == XSettingString {
		return "String"
	} else if typ == XSettingColour {
		return "Colour"
	}
	return fmt.Sprintf("Invalid (%d)", typ)
}

func (setting *XSetting) mustBe(typ XSettingType) error {
	if typ == setting.Type {
		return nil
	}
	return fmt.Errorf("Setting must be %s -- is %s", typ, setting.Type)
}
