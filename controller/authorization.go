package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/galaxy-solar/starstore/model/auth"
	"net/http"
	"github.com/galaxy-solar/starstore/response"
	"github.com/galaxy-solar/starstore/i18n"
)

type LoginInput struct {
	account string
	password string
}

func AuthLogin(g *gin.Context) {
	authType := g.Param("authType")

	authorization := LoginInput{}
	g.BindJSON(&authorization)

	switch authType {
	case "enterprise":
		var authorEnterprise *auth.EnterpriseAuthorization
		modeledDB := DB().Model(&auth.EnterpriseAuthorization{})
		modeledDB.Where(&auth.EnterpriseAuthorization{AuthBase: auth.AuthBase{Username: authorization.account, Password: auth.Password(authorization.password)}}).
			First(authorEnterprise)
		if authorEnterprise == nil {
			modeledDB.Where(&auth.EnterpriseAuthorization{AuthBase: auth.AuthBase{Email: authorization.account, Password: auth.Password(authorization.password)}}).
				First(authorEnterprise)
		}
		if authorEnterprise == nil {
			modeledDB.Where(&auth.EnterpriseAuthorization{AuthBase: auth.AuthBase{Mobile: authorization.account, Password: auth.Password(authorization.password)}}).
				First(authorEnterprise)
		}
		if authorEnterprise == nil {
			g.JSON(http.StatusNotFound, &response.Response{
				response.OK,
				i18n.I18NViper.GetString("message.auth.enterprise.notfound"),
				nil,
			})
		}
		g.JSON(http.StatusOK, &response.Response{
			response.OK,
			i18n.I18NViper.GetString("message.auth.login.ok"),
			&authorEnterprise,
		})
	case "user":
		var authorUser *auth.UserAuthorization
		modeledDB := DB().Model(&auth.UserAuthorization{})
		modeledDB.Where(&auth.UserAuthorization{AuthBase: auth.AuthBase{Username: authorization.account, Password: auth.Password(authorization.password)}}).
			First(authorUser)
		if authorUser == nil {
			modeledDB.Where(&auth.UserAuthorization{AuthBase: auth.AuthBase{Email: authorization.account, Password: auth.Password(authorization.password)}}).
				First(authorUser)
		}
		if authorUser == nil {
			modeledDB.Where(&auth.UserAuthorization{AuthBase: auth.AuthBase{Mobile: authorization.account, Password: auth.Password(authorization.password)}}).
				First(authorUser)
		}
		if authorUser == nil {
			g.JSON(http.StatusNotFound, &response.Response{
				response.OK,
				i18n.I18NViper.GetString("message.auth.enterprise.notfound"),
				nil,
			})
		}
		g.JSON(http.StatusOK, &response.Response{
			response.OK,
			i18n.I18NViper.GetString("message.auth.login.ok"),
			&authorUser,
		})
	default:
		g.JSON(http.StatusNotFound, &response.Response{
			response.NotFound,
			i18n.I18NViper.GetString("message.auth.login.unsupported"),
			nil,
		})
	}
}