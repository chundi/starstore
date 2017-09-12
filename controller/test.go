package controller

import (
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"github.com/tidwall/gjson"
)

func TempGet(g *gin.Context) {
	a := struct {
		Id string 	`form:"id"`
		Name string	`form:"name"`
	} {
	}
	BindWithTeeReader(g, &a)
	//g.BindJSON(&a)
	Logger.Info("TempGet: ", a.Id)
	body, _ := ioutil.ReadAll(g.Request.Body)
	id := gjson.GetBytes(body, "id")
	Logger.Info("get id again: ", id)
	g.String(200, "temp get, a: %v, id: %s", a, id)
}