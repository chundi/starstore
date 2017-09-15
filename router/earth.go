package router

import (
	"github.com/galaxy-solar/starstore/controller"
	"github.com/gin-gonic/gin"
)

func EarthRoute(group *gin.RouterGroup) {
	group.GET("/space", controller.SpaceGet)
	group.POST("/space", controller.SpacePost)
	group.GET("/space/:id", controller.SpaceDetailGet)
	group.PUT("/space/:id", controller.SpaceDetailPut)
	group.DELETE("/space/:id", controller.SpaceDetailDelete)

	group.GET("/device", controller.DeviceGet)
	group.POST("/device", controller.DevicePost)
	group.GET("/device/:id", controller.DeviceDetailGet)
	group.PUT("/device/:id", controller.DeviceDetailPut)
	group.DELETE("/device/:id", controller.DeviceDetailDelete)
}
