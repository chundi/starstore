package earth

import (
	"github.com/galaxy-solar/starstore/i18n"
	"github.com/galaxy-solar/starstore/model"
	"github.com/spf13/viper"
)

type Device struct {
	model.Base
}

func (device Device) GetMessage() *viper.Viper {
	return i18n.I18NViper.Sub("message.auth.device")
}

func (device *Device) GetEntity() interface{} {
	return device
}

type DeviceMeta struct {
	model.Meta
}
