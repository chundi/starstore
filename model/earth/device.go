package earth

import (
	"github.com/galaxy-solar/starstore/i18n"
	"github.com/galaxy-solar/starstore/model"
	"github.com/jinzhu/gorm"
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

func GetDeviceByToken(db *gorm.DB, token string, d *Device) bool {
	r := db.Where("token = ?", token).First(d)
	if r.Error != nil {
		flog := model.Logger.WithField("GetDeviceByToken:", token)
		if r.RecordNotFound() {
			flog.Error("Record Not Found!")
		} else {
			flog.Error(r.Error.Error())
		}
		return false
	}
	return true
}
