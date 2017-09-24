package auth

import (
	"github.com/galaxy-solar/starstore/conf"
	"github.com/galaxy-solar/starstore/i18n"
	"github.com/galaxy-solar/starstore/model"
	"github.com/spf13/viper"
)

type Enterprise struct {
	model.Base

	Authorization EnterpriseAuthorization `gorm:"-" json:"authorization,omitempty"`
}

func (enterprise Enterprise) GetAuthorization() Authorizer {
	return enterprise.Authorization
}

func (enterprise Enterprise) GetAuthType() string {
	return conf.ENTITY_ENTERPRISE
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
