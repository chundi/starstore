package controller

import (
	"github.com/galaxy-solar/starstore/model/earth"
	"github.com/gin-gonic/gin"
)

func SpaceGet(g *gin.Context) {
	space := earth.Space{}
	var spaceList = []earth.Space{}
	BaseGet(g, DBWithContext(g), &space, &spaceList)
}

func SpacePost(g *gin.Context) {
	space := earth.Space{}
	BasePost(g, DB(), &space)
}

func SpaceDetailGet(g *gin.Context) {
	space := earth.Space{}
	BaseDetailGet(g, DBWithContext(g), &space)
}

func SpaceDetailPut(g *gin.Context) {
	space := earth.Space{}
	BaseDetailPut(g, DBWithContext(g), &space)
}

func SpaceDetailDelete(g *gin.Context) {
	space := earth.Space{}
	BaseDetailDelete(g, DBWithContext(g), &space)
}
