package router

import (
	"github.com/gin-gonic/gin"
	"github.com/galaxy-solar/starstore/controller"
)

func AuthRoute(group *gin.RouterGroup) {
	group.GET("/enterprise", controller.EnterpriseGet)
	group.POST("/enterprise", controller.EnterprisePost)
}
