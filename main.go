package main

import (
	"github.com/galaxy-solar/starstore/model/earth"
	"github.com/galaxy-solar/starstore/model"
	"github.com/galaxy-solar/starstore/model/auth"

	"github.com/galaxy-solar/starstore/conf"
	"github.com/galaxy-solar/starstore/router"
	"github.com/fvbock/endless"
	"github.com/gin-gonic/gin"
	"github.com/galaxy-solar/starstore/router/middleware"
	"fmt"
	"github.com/sirupsen/logrus"
)

var rootRouter *gin.Engine

func initTables(tables ...interface{}) {
	if conf.AppConfig.RegenerateTables {
		model.DropTable(tables...)
	}
	model.MigrateTable(tables...)
}

func initRouter() {
	rootRouter = gin.New()
	rootRouter.Use(middleware.HttpLogger())
	rootRouter.Use(gin.Logger())
	rootRouter.Use(gin.Recovery())

	router.LoginRoute(rootRouter)

	currentVersion := rootRouter.Group(fmt.Sprintf("/%s", conf.AppConfig.Api.Version))
	currentVersion.Use(middleware.Authorized())
	router.AuthRoute(currentVersion.Group("/auth"))
	router.EarthRoute(currentVersion.Group("/earth"))
	router.SolarRoute(currentVersion.Group("/solar"))
}

func init() {
	initTables(&auth.Enterprise{}, &auth.EnterpriseMeta{}, &auth.EnterpriseAuthorization{},
		&auth.User{}, &auth.UserMeta{}, &auth.UserAuthorization{})
	initTables(&earth.Space{}, &earth.SpaceMeta{},
		&earth.Device{}, &earth.DeviceMeta{})
	initRouter()
}

func main() {
	logrus.Info("starting server ")
	endless.ListenAndServe(fmt.Sprintf(":%v", conf.AppConfig.Port), rootRouter)
	//rootRouter.Run(fmt.Sprintf(":%v", conf.AppConfig.Port))
}