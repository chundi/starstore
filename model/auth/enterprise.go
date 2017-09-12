package auth

import (
	"github.com/galaxy-solar/starstore/model"
	"github.com/spf13/viper"
	"github.com/galaxy-solar/starstore/i18n"
)

type Enterprise struct {
	model.Base
}

func (enterprise Enterprise) GetMessage() *viper.Viper {
	return i18n.I18NViper.Sub("message.auth.enterprise")
}

func (enterprise *Enterprise) GetEntity() interface{} {
	return enterprise
}

type EnterpriseMeta struct {
	model.Meta
}
