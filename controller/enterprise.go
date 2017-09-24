package controller

import (
	"errors"
	"github.com/galaxy-solar/starstore/i18n"
	"github.com/galaxy-solar/starstore/model"
	"github.com/galaxy-solar/starstore/model/auth"
	"github.com/galaxy-solar/starstore/response"
	"github.com/galaxy-solar/starstore/util"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"net/http"
)

func EnterpriseGet(g *gin.Context) {
	enterprise := auth.Enterprise{}
	var enterpriseList = []auth.Enterprise{}
	BaseGet(g, DBWithContext(g), &enterprise, &enterpriseList)
}

func EnterprisePost(g *gin.Context) {
	enterprise := auth.Enterprise{}
	enterprise.AddHandler(model.POSITION_POST_TRANSACTION_END, func(g *gin.Context, db *gorm.DB) error {
		authorization := auth.EnterpriseAuthorization{}
		if body, err := ioutil.ReadAll(g.Request.Body); err != nil {
			Logger.Error(err)
			g.JSON(http.StatusInternalServerError, &response.Response{
				Code:    response.Error,
				Message: i18n.I18NViper.GetString("message.auth.enterprise.jsonerror"),
			})
			return errors.New(i18n.I18NViper.GetString("message.auth.enterprise.jsonerror"))
		} else {
			authJson := gjson.GetBytes(body, "authorization")
			gjson.Unmarshal([]byte(authJson.Raw), &authorization.AuthBase)

			var validationMessage string
			validatedErrorFound := true
			if !util.REG_EMAIL.Match([]byte(authorization.AuthBase.Email)) {
				validationMessage = i18n.I18NViper.GetString("message.auth.validation.email")
			} else if !util.REG_MOBILE.Match([]byte(authorization.AuthBase.Mobile)) {
				validationMessage = i18n.I18NViper.GetString("message.auth.validation.mobile")
			} else if authorization.AuthBase.Mobile == authorization.AuthBase.Username {
				validationMessage = i18n.I18NViper.GetString("message.auth.validation.usernamesameasmobile")
			} else if authorization.AuthBase.Email == authorization.AuthBase.Username {
				validationMessage = i18n.I18NViper.GetString("message.auth.validation.usernamesameasemail")
			} else {
				validatedErrorFound = false
			}
			if validatedErrorFound {
				g.JSON(http.StatusInternalServerError, &response.Response{
					Code:    response.Error,
					Message: validationMessage,
				})
				return errors.New(validationMessage)
			} else {
				authorization.Id = enterprise.Id
				Logger.Debug("enterprise creating authorization: ", authorization)
			}
		}
		if authErr := db.Create(&authorization).Error; authErr != nil {
			db.Rollback()
			Logger.Error(authErr)
			g.JSON(http.StatusInternalServerError, &response.Response{
				Code:    response.Error,
				Message: i18n.I18NViper.GetString("message.auth.enterprise.authuncreated"),
			})
			return errors.New(i18n.I18NViper.GetString("message.auth.enterprise.authuncreated"))
		}
		return nil
	})
	BasePost(g, DB(), &enterprise)
}

func EnterpriseDetailGet(g *gin.Context) {
	enterprise := auth.Enterprise{}
	enterprise.AddHandler(model.POSITION_DETAIL_GET_AFTER, func(g *gin.Context, db *gorm.DB) error {
		id := g.Param("id")
		if db.Where("id = ?", id).First(&enterprise.Authorization).RecordNotFound() {
			g.JSON(http.StatusNotFound, &response.Response{
				Code:    response.NotFound,
				Message: BaseMessage(&enterprise, "resourcenotfound"),
				Data:    nil,
			})
			return errors.New("authorization get failed")
		} else {
			enterprise.Authorization.Password = ""
		}
		return nil
	})
	BaseDetailGet(g, DBWithContext(g), &enterprise)
}

func EnterpriseDetailPut(g *gin.Context) {
	enterprise := auth.Enterprise{}
	BaseDetailPut(g, DBWithContext(g), &enterprise)
}

func EnterpriseDetailDelete(g *gin.Context) {
	enterprise := auth.Enterprise{}
	BaseDetailDelete(g, DBWithContext(g), &enterprise)
}
