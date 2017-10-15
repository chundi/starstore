package main

import (
	"github.com/galaxy-solar/starstore/model"
	"github.com/galaxy-solar/starstore/model/auth"
	"github.com/galaxy-solar/starstore/model/earth"

	"fmt"

	"github.com/fvbock/endless"
	"github.com/galaxy-solar/starstore/conf"
	"github.com/galaxy-solar/starstore/message"
	"github.com/galaxy-solar/starstore/router"
	"github.com/galaxy-solar/starstore/router/middleware"
	"github.com/gin-gonic/gin"
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

	// serve websocket
	rootRouter.GET("/ws", func(c *gin.Context) {
		message.WsHandler(c.Writer, c.Request)
	})

	rootRouter.Use(middleware.HttpLogger())
	rootRouter.Use(gin.Logger())
	rootRouter.Use(gin.Recovery())

	router.LoginRoute(rootRouter)
	router.UtilRoute(rootRouter)

	currentVersion := rootRouter.Group(fmt.Sprintf("/%s", conf.AppConfig.Api.Version))
	currentVersion.Use(middleware.CapabilitiesMiddleware())
	router.AuthRoute(currentVersion.Group("/auth"))
	router.EarthRoute(currentVersion.Group("/earth"))
	router.SolarRoute(currentVersion.Group("/solar"))
}

func init() {
	initTables(&auth.Enterprise{}, &auth.EnterpriseMeta{}, &auth.EnterpriseAuthorization{},
		&auth.Employee{}, &auth.EmployeeMeta{}, &auth.EmployeeAuthorization{})
	initTables(&earth.Space{}, &earth.SpaceMeta{},
		&earth.Device{}, &earth.DeviceMeta{})
	initRouter()
}

func main() {
	logrus.Info("starting server ")
	endless.ListenAndServe(fmt.Sprintf(":%v", conf.AppConfig.Port), rootRouter)
	//rootRouter.Run(fmt.Sprintf(":%v", conf.AppConfig.Port))
}
