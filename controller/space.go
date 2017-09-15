package controller

import (
	"github.com/galaxy-solar/starstore/model/earth"
	"github.com/gin-gonic/gin"
)

func SpaceGet(g *gin.Context) {
	space := earth.Space{}
	var spaceList = []earth.Space{}
	BaseGet(g, DB(), &space, &spaceList)
}

func SpacePost(g *gin.Context) {
	space := earth.Space{}
	BasePost(g, DB(), &space)
}

func SpaceDetailGet(g *gin.Context) {
	space := earth.Space{}
	BaseDetailGet(g, DB(), &space)
}

func SpaceDetailPut(g *gin.Context) {
	space := earth.Space{}
	BaseDetailPut(g, DB(), &space)
}

func SpaceDetailDelete(g *gin.Context) {
	space := earth.Space{}
	BaseDetailDelete(g, DB(), &space)
}
