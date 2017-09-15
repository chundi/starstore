package controller

import (
	"github.com/galaxy-solar/starstore/model/earth"
	"github.com/gin-gonic/gin"
)

func DeviceGet(g *gin.Context) {
	device := earth.Device{}
	var deviceList = []earth.Device{}
	BaseGet(g, DB(), &device, &deviceList)
}

func DevicePost(g *gin.Context) {
	device := earth.Device{}
	BasePost(g, DB(), &device)
}

func DeviceDetailGet(g *gin.Context) {
	device := earth.Device{}
	BaseDetailGet(g, DB(), &device)
}

func DeviceDetailPut(g *gin.Context) {
	device := earth.Device{}
	BaseDetailPut(g, DB(), &device)
}

func DeviceDetailDelete(g *gin.Context) {
	device := earth.Device{}
	BaseDetailDelete(g, DB(), &device)
}
