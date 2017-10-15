package router

import (
	"github.com/galaxy-solar/starstore/controller"
	"github.com/gin-gonic/gin"
)

func UtilRoute(engine *gin.Engine) {
	engine.POST("rfiddecoder", controller.RfidDecoder)
	engine.POST("rfidencoder", controller.RfidEncoder)
}
