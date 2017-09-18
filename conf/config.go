package conf

import (
	"fmt"
	"github.com/spf13/viper"
)

type AppConfiguration struct {
	AppName          string
	Port             int
	RunMode          string
	RegenerateTables bool

	Db struct {
		Adapter  string
		Conn     string
		User     string
		Password string
		Server   string
		Database string
	}

	Log struct {
		Output string
		Level  string
		Format string
	}

	Api struct {
		Version string
	}
}

var AppConfig AppConfiguration

func DefaultConfigurator(fileName string) *viper.Viper {
	configurator := viper.New()
	configurator.AddConfigPath("./conf")
	configurator.SetConfigName(fileName)
	if err := configurator.ReadInConfig(); err != nil {
		panic(fmt.Sprintf("loading configuration %s error: %s", fileName, err))
	}
	return configurator
}

func IsDevelopMode() bool {
	return !(AppConfig.RunMode == "production" || AppConfig.RunMode == "product" || AppConfig.RunMode == "release")
}

func init() {
	cfgViper := DefaultConfigurator("configuration")
	cfgViper.SetDefault("RegenerateTables", false)
	err := cfgViper.Unmarshal(&AppConfig)
	if err != nil {
		panic(fmt.Sprintf("configuration unmarshal error: %s", err))
	}

	var cfgFileName string
	var viper *viper.Viper
	var appCfg AppConfiguration

	switch AppConfig.RunMode {
	case "production", "release", "product":
		cfgFileName = "production"
	case "development", "dev", "develop":
		fallthrough
	default:
		cfgFileName = "development"
	}

	viper = DefaultConfigurator(cfgFileName)
	if err := viper.Unmarshal(&appCfg); err != nil {
		panic(fmt.Sprintf("db configuration unmarshal error: %s", err))
	}
	AppConfig.RegenerateTables = appCfg.RegenerateTables
	AppConfig.Db = appCfg.Db
	AppConfig.Log = appCfg.Log

	if !IsDevelopMode() {
		AppConfig.RegenerateTables = false
	}
}
