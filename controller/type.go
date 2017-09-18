package controller

import (
	"github.com/galaxy-solar/starstore/model/solar"
	"github.com/gin-gonic/gin"
)

func TypeGet(g *gin.Context) {

}

func TypePost(g *gin.Context) {
	typeIndicator := solar.TypeIndicator{}
	IndicatorBasePost(g, DB(), &typeIndicator)
}

func TypeDetailGet(g *gin.Context) {

}

func TypeDetailPut(g *gin.Context) {

}

func TypeDetailDelete(g *gin.Context) {

}
