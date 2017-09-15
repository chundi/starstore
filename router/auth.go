package router

import (
	"github.com/galaxy-solar/starstore/controller"
	"github.com/gin-gonic/gin"
)

func AuthRoute(group *gin.RouterGroup) {
	group.GET("/enterprise", controller.EnterpriseGet)
	group.POST("/enterprise", controller.EnterprisePost)
	group.GET("/enterprise/:id", controller.EnterpriseDetailGet)
	group.PUT("/enterprise/:id", controller.EnterpriseDetailPut)
	group.DELETE("/enterprise/:id", controller.EnterpriseDetailDelete)

	group.GET("/user", controller.UserGet)
	group.POST("/user", controller.UserPost)
	group.GET("/user/:id", controller.UserDetailGet)
	group.PUT("/user/:id", controller.UserDetailPut)
	group.DELETE("/user/:id", controller.UserDetailDelete)
}
