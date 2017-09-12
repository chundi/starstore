package router

import (
	"github.com/gin-gonic/gin"
	"github.com/galaxy-solar/starstore/controller"
)

func EarthRoute(group *gin.RouterGroup) {
	group.GET("/space", controller.SpaceGet)
	group.POST("/space", controller.SpacePost)
	//
	//group.GET("/space/:spaceId", controller.SpaceDetailGet)
	//group.PUT("/space/:spaceId", controller.SpaceDetailPut)
	//group.DELETE("/space/:spaceId", controller.SpaceDetailDelete)
}

