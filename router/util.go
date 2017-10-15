package router

import (
	"github.com/galaxy-solar/starstore/controller"
	"github.com/gin-gonic/gin"
)

func UtilRoute(engine *gin.Engine) {
	engine.GET("rfiddecoder/:rfid", controller.RfidDecoder)
	engine.GET("rfidencoder/:code", controller.RfidEncoder)
}
