package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/galaxy-solar/starstore/model/auth"
	"net/http"
	"github.com/galaxy-solar/starstore/response"
	"github.com/galaxy-solar/starstore/i18n"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"github.com/galaxy-solar/starstore/model"
	"github.com/jinzhu/gorm"
	"errors"
)

func EnterpriseGet(g *gin.Context) {
	enterprise := auth.Enterprise{}
	var enterpriseList = []auth.Enterprise{}
	BaseGet(g, DB(), &enterprise, &enterpriseList)
}

func EnterprisePost(g *gin.Context) {
	enterprise := auth.Enterprise{}
	enterprise.AddHandler(model.POSITION_POST_TRANSACTION_END, func (g *gin.Context, db *gorm.DB) error {
		authorization := auth.EnterpriseAuthorization{}
		if body, err := ioutil.ReadAll(g.Request.Body); err != nil {
			Logger.Error(err)
			g.JSON(http.StatusInternalServerError, &response.Response{
				Code: response.Error,
				Message: i18n.I18NViper.GetString("message.auth.enterprise.jsonerror"),
			})
		} else {
			authJson := gjson.GetBytes(body, "authorization")
			gjson.Unmarshal([]byte(authJson.Raw), &authorization.AuthBase)
			authorization.EnterpriseId = enterprise.Id
			Logger.Debug("enterprise creating authorization: ", authorization)
		}
		if authErr := db.Create(&authorization).Error; authErr != nil {
			db.Rollback()
			Logger.Error(authErr)
			g.JSON(http.StatusInternalServerError, &response.Response{
				Code: response.Error,
				Message: i18n.I18NViper.GetString("message.auth.enterprise.authuncreated"),
			})
			return errors.New("authorization create failed")
		}
		return nil
	})
	BasePost(g, DB(), &enterprise)
}

func EnterpriseDetailGet(g *gin.Context) {
	enterprise := auth.Enterprise{}
	BaseDetailGet(g, DB(), &enterprise)
}

func EnterpriseDetailPut(g *gin.Context) {
	enterprise := auth.Enterprise{}
	BaseDetailPut(g, DB(), &enterprise)
}

func EnterpriseDetailDelete(g *gin.Context) {
	enterprise := auth.Enterprise{}
	BaseDetailDelete(g, DB(), &enterprise)
}
