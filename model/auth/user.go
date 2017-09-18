package auth

import (
	"github.com/galaxy-solar/starstore/i18n"
	"github.com/galaxy-solar/starstore/model"
	"github.com/spf13/viper"
)

type User struct {
	model.Base

	Authorization UserAuthorization `gorm:"-" json:"authorization,omitempty"`
}

func (user User) GetMessage() *viper.Viper {
	return i18n.I18NViper.Sub("message.auth.user")
}

func (user *User) GetEntity() interface{} {
	return user
}

type UserMeta struct {
	model.Meta
}
