package controller

import (
	"errors"
	"github.com/galaxy-solar/starstore/i18n"
	"github.com/galaxy-solar/starstore/model"
	"github.com/galaxy-solar/starstore/model/auth"
	"github.com/galaxy-solar/starstore/response"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"net/http"
)

func UserGet(g *gin.Context) {
	user := auth.User{}
	var userList = []auth.User{}
	BaseGet(g, DB(), &user, &userList)
}

func UserPost(g *gin.Context) {
	user := auth.User{}
	user.AddHandler(model.POSITION_POST_TRANSACTION_END, func(g *gin.Context, db *gorm.DB) error {
		authorization := auth.UserAuthorization{}
		if body, err := ioutil.ReadAll(g.Request.Body); err != nil {
			Logger.Error(err)
			g.JSON(http.StatusInternalServerError, &response.Response{
				Code:    response.Error,
				Message: i18n.I18NViper.GetString("message.auth.user.jsonerror"),
			})
		} else {
			authJson := gjson.GetBytes(body, "authorization")
			gjson.Unmarshal([]byte(authJson.Raw), &authorization.AuthBase)
			authorization.Id = user.Id
			Logger.Debug("user creating authorization: ", authorization)
		}
		if authErr := db.Create(&authorization).Error; authErr != nil {
			db.Rollback()
			Logger.Error(authErr)
			g.JSON(http.StatusInternalServerError, &response.Response{
				Code:    response.Error,
				Message: i18n.I18NViper.GetString("message.auth.user.authuncreated"),
			})
			return errors.New("authorization create failed")
		}
		return nil
	})
	BasePost(g, DB(), &user)
}

func UserDetailGet(g *gin.Context) {
	user := auth.User{}
	user.AddHandler(model.POSITION_DETAIL_GET_AFTER, func(g *gin.Context, db *gorm.DB) error {
		id := g.Param("id")
		if db.Where("id = ?", id).First(&user.Authorization).RecordNotFound() {
			g.JSON(http.StatusNotFound, &response.Response{
				Code:    response.NotFound,
				Message: BaseMessage(&user, "resourcenotfound"),
				Data:    nil,
			})
			return errors.New("authorization get failed")
		} else {
			user.Authorization.Password = ""
		}
		return nil
	})
	BaseDetailGet(g, DB(), &user)
}

func UserDetailPut(g *gin.Context) {
	user := auth.User{}
	BaseDetailPut(g, DB(), &user)
}

func UserDetailDelete(g *gin.Context) {
	user := auth.User{}
	BaseDetailDelete(g, DB(), &user)
}
