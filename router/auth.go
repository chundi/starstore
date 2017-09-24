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

	group.GET("/employee", controller.EmployeeGet)
	group.POST("/employee", controller.EmployeePost)
	group.GET("/employee/:id", controller.EmployeeDetailGet)
	group.PUT("/employee/:id", controller.EmployeeDetailPut)
	group.DELETE("/employee/:id", controller.EmployeeDetailDelete)
}
