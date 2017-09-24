package router

import (
	"fmt"
	"github.com/galaxy-solar/starstore/conf"
	"github.com/galaxy-solar/starstore/controller"
	"github.com/gin-gonic/gin"
)

func LoginRoute(engine *gin.Engine) {
	engine.POST(fmt.Sprintf("/%s/login/:entity", conf.AppConfig.Api.Version), controller.AuthLogin)
	engine.POST("/tmp", controller.TempGet)
}
