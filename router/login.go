package router

import (
	"github.com/gin-gonic/gin"
	"fmt"
	"github.com/galaxy-solar/starstore/controller"
	"github.com/galaxy-solar/starstore/conf"
)

func LoginRoute(engine *gin.Engine) {
	engine.POST(fmt.Sprintf("/%s/login/:authType", conf.AppConfig.Api.Version), controller.AuthLogin)
	engine.POST("/tmp", controller.TempGet)
}
