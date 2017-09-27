package model

import (
	"errors"
	"fmt"
	"github.com/galaxy-solar/starstore/conf"
	"github.com/galaxy-solar/starstore/log"
	"github.com/galaxy-solar/starstore/util"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/sirupsen/logrus"
)

var (
	DB     *gorm.DB
	logCfg = conf.AppConfig.Log
	Logger *logrus.Entry
)

func init() {
	DB = InitDB()
	Logger = log.NewLogger(logCfg.Format, logCfg.Level, logCfg.Output).WithField("module", "model")
}

func InitDB() *gorm.DB {
	var err error

	if conf.AppConfig.Db.Adapter == "mysql" {
		DB, err = gorm.Open("mysql", fmt.Sprintf(conf.AppConfig.Db.Conn,
			conf.AppConfig.Db.User,
			conf.AppConfig.Db.Password,
			conf.AppConfig.Db.Server,
			conf.AppConfig.Db.Database))
	} else if conf.AppConfig.Db.Adapter == "postgres" {
		DB, err = gorm.Open("postgres", fmt.Sprintf(conf.AppConfig.Db.Conn,
			conf.AppConfig.Db.User,
			conf.AppConfig.Db.Password,
			conf.AppConfig.Db.Server,
			conf.AppConfig.Db.Database))
	} else {
		panic(errors.New("unsupported database adapter."))
	}

	if err == nil {
		DB.SingularTable(true)
	} else {
		panic(err)
	}
	return DB
}

func MigrateTable(tables ...interface{}) {
	for _, table := range tables {
		DB.AutoMigrate(table)
	}
}

func DropTable(tables ...interface{}) {
	for _, table := range tables {
		fmt.Println("droping table: " + util.ModelType(table).Name())
		DB.DropTable(table)
	}
}
