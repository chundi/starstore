package router

import (
	"github.com/galaxy-solar/starstore/controller"
	"github.com/gin-gonic/gin"
)

func SolarRoute(group *gin.RouterGroup) {
	group.GET("/type", controller.TypeGet)
	group.POST("/type", controller.TypePost)
	group.GET("/type/:id", controller.TypeDetailGet)
	group.PUT("/type/:id", controller.TypeDetailPut)
	group.DELETE("/type/:id", controller.TypeDetailDelete)
}
