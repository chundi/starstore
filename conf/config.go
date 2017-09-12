package conf

import (
	"github.com/spf13/viper"
	"fmt"
)

type AppConfiguration struct {
	AppName string
	Port int
	RunMode string
	RegenerateTables bool

	Db struct {
		Adapter string
		Conn string
		User string
		Password string
		Server string
		Database string
	}

	Log struct {
		Output string
	}

	Api struct {
		Version string
	}
}

var AppConfig AppConfiguration

func DefaultConfigor(fileName string) *viper.Viper {
	configor := viper.New()
	configor.AddConfigPath("./conf")
	configor.SetConfigName(fileName)
	if err := configor.ReadInConfig(); err != nil {
		panic(fmt.Sprintf("loading configuration %s error: %s", fileName, err))
	}
	return configor
}

func IsDevelopMode() bool {
	return !(AppConfig.RunMode == "production" || AppConfig.RunMode == "product" || AppConfig.RunMode == "release")
}

func init() {
	ConfigViper := DefaultConfigor("configuration")
	ConfigViper.SetDefault("RegenerateTables", false)
	err := ConfigViper.Unmarshal(&AppConfig)
	if err != nil {
		panic(fmt.Sprintf("configuration unmarshal error: %s", err))
	}
	if !IsDevelopMode() {
		AppConfig.RegenerateTables = false
	}
	var dbConfigFileName string
	var dbViper *viper.Viper
	var dbConfigor AppConfiguration
	switch AppConfig.RunMode {
	case "production", "release", "product":
		dbConfigFileName = "production"
	case "development", "dev", "develop":
		fallthrough
	default :
		dbConfigFileName = "development"
	}
	dbViper = DefaultConfigor(dbConfigFileName)
	if err := dbViper.Unmarshal(&dbConfigor); err != nil {
		panic(fmt.Sprintf("db configuration unmarshal error: %s", err))
	}
	AppConfig.Db = dbConfigor.Db
}