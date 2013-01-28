package xsettings

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/xprop"

	"github.com/BurntSushi/xgb/xproto"
)

type XSettings struct {
	X      *xgbutil.XUtil
	Serial uint32
	props  map[string]*XSetting
}

type XSetting struct {
	Type   XSettingType
	Serial uint32
	// Values -- Note it's probably better to use the 
	// getters in XSettings, than to use this.
	String  string
	Integer int32
	Colour  XSColour
}

type XSColour struct {
	R, G, B, A uint16
}

type XSettingType uint8

const (
	XSettingInteger XSettingType = iota
	XSettingString
	XSettingColour
)

var (
	atomXSettings xproto.Atom
)

// Creates a new XSettings object, and initializes it with data populated 
// from the settings manager that is currently active.
func New(X *xgbutil.XUtil) (*XSettings, error) {
	var err error

	if atomXSettings == 0 {
		atomXSettings, err = xprop.Atom(X, "_XSETTINGS_S0", false)
		if err != nil {
			return nil, err
		}
	}

	settings := &XSettings{
		X: X,
	}

	if err = settings.Refresh(); err != nil {
		return nil, err
	}

	return settings, nil
}

// Refreshes the internal cache of properties.
//
// NOTE: you do not have to call this, unless you want a more up-to-date value
// set than when you first created it with New -- or the last Refresh call.
func (xs *XSettings) Refresh() error {
	xs.props = make(map[string]*XSetting)

	gso := xproto.GetSelectionOwner(xs.X.Conn(), atomXSettings)

	reply, err := gso.Reply()
	if err != nil {
		return err
	}

	prop, err := xprop.GetProperty(xs.X, reply.Owner, "_XSETTINGS_SETTINGS")
	if err != nil {
		return err
	}

	data := bytes.NewReader(prop.Value)

	var endian binary.ByteOrder

	var tmp byte

	if tmp, err = data.ReadByte(); err != nil {
		return err
	}

	if tmp == 0 {
		endian = binary.LittleEndian
	} else {
		endian = binary.BigEndian
	}

	buff := make([]byte, 3, 3)

	if n, err := data.Read(buff); n != 3 || err != nil {
		return err
	}

	if err = binary.Read(data, endian, &xs.Serial); err != nil {
		return err
	}

	var num uint32

	if err = binary.Read(data, endian, &num); err != nil {
		return err
	}

	for i := uint32(0); i < num; i++ {
		setting, name, err := xs.readSetting(data, endian)
		if err != nil {
			return err
		}
		xs.props[name] = setting
	}

	return nil
}

func (xs *XSettings) readSetting(data *bytes.Reader, endian binary.ByteOrder) (*XSetting, string, error) {
	var name string
	var err error

	setting := &XSetting{}

	var tmp byte

	if tmp, err = data.ReadByte(); err != nil {
		return nil, "", err
	}

	if tmp > 2 {
		return nil, "", fmt.Errorf("Invalid type identifier %d", tmp)
	}

	setting.Type = XSettingType(tmp)

	if _, err = data.ReadByte(); err != nil {
		return nil, "", err
	}

	var l16 uint16

	if err = binary.Read(data, endian, &l16); err != nil {
		return nil, "", err
	}

	buff := make([]byte, int(l16), int(l16))

	if n, err := data.Read(buff); err != nil || n != int(l16) {
		return nil, "", err
	}

	name = string(buff)

	pad(data, 4)

	if err = binary.Read(data, endian, &setting.Serial); err != nil {
		return nil, "", err
	}

	switch setting.Type {
	case XSettingInteger:
		if err = binary.Read(data, endian, &setting.Integer); err != nil {
			return nil, "", err
		}
	case XSettingString:
		var l32 uint32
		if err = binary.Read(data, endian, &l32); err != nil {
			return nil, "", err
		}

		buff := make([]byte, int(l32), int(l32))

		if n, err := data.Read(buff); err != nil || n != int(l32) {
			return nil, "", err
		}

		setting.String = string(buff)
		pad(data, 4)
	case XSettingColour:
		var r, g, b, a uint16
		if err = binary.Read(data, endian, &r); err != nil {
			return nil, "", err
		}
		if err = binary.Read(data, endian, &g); err != nil {
			return nil, "", err
		}
		if err = binary.Read(data, endian, &b); err != nil {
			return nil, "", err
		}
		if err = binary.Read(data, endian, &a); err != nil {
			return nil, "", err
		}

		setting.Colour = XSColour{r, g, b, a}
	default:
		panic("This shouldn't be happening! D:")
	}

	return setting, name, nil
}

func pad(data *bytes.Reader, align int64) {
	for pos, _ := data.Seek(0, 1); pos%align != 0; pos++ {
		data.ReadByte()
	}
}
