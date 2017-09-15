package earth

import (
	"github.com/galaxy-solar/starstore/i18n"
	"github.com/galaxy-solar/starstore/model"
	"github.com/spf13/viper"
)

type Space struct {
	*model.Base
}

func (space Space) GetMessage() *viper.Viper {
	return i18n.I18NViper.Sub("message.auth.space")
}

func (space *Space) GetEntity() interface{} {
	return space
}

type SpaceMeta struct {
	model.Meta
}
