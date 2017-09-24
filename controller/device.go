package controller

import (
	"github.com/galaxy-solar/starstore/model/earth"
	"github.com/gin-gonic/gin"
)

func DeviceGet(g *gin.Context) {
	device := earth.Device{}
	var deviceList = []earth.Device{}
	BaseGet(g, DBWithContext(g), &device, &deviceList)
}

func DevicePost(g *gin.Context) {
	device := earth.Device{}
	if claim, ok := GetCapabilityClaims(g); ok && claim.IsEnterprise() {
		device.OwnerId = claim.Id
	}
	BasePost(g, DB(), &device)
}

func DeviceDetailGet(g *gin.Context) {
	device := earth.Device{}
	BaseDetailGet(g, DBWithContext(g), &device)
}

func DeviceDetailPut(g *gin.Context) {
	device := earth.Device{}
	BaseDetailPut(g, DBWithContext(g), &device)
}

func DeviceDetailDelete(g *gin.Context) {
	device := earth.Device{}
	BaseDetailDelete(g, DBWithContext(g), &device)
}
