package i18n

import (
	"github.com/spf13/viper"
	"fmt"
)

var I18NViper *viper.Viper

func init() {
	I18NViper = viper.New()
	I18NViper.AddConfigPath("./i18n")
	I18NViper.SetConfigName("zh_Hans") //Todo: change different locale
	if err := I18NViper.ReadInConfig(); err != nil {
		panic(fmt.Sprintf("configuration error: %s", err))
	}
}