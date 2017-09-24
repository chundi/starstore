package controller

import (
	"github.com/galaxy-solar/starstore/model/auth"
	"github.com/galaxy-solar/starstore/response"
	"github.com/gin-gonic/gin"
	"net/http"

	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/galaxy-solar/starstore/conf"
	"github.com/galaxy-solar/starstore/i18n"
	"github.com/galaxy-solar/starstore/model"
	"github.com/galaxy-solar/starstore/util"
	"strings"
	"time"
)

const (
	ACCESS_TOKEN string = "jwt"
)

type LoginInput struct {
	Account  string
	Password string
}

func AuthLogin(g *gin.Context) {
	entity := g.Param("entity")
	input := LoginInput{}
	g.BindJSON(&input)

	var baser model.Baser
	var authorization interface{}

	switch entity {
	case conf.ENTITY_ENTERPRISE:
		baser = &auth.Enterprise{}
		authorization = &auth.EnterpriseAuthorization{}
		baser.GetEntity().(*auth.Enterprise).Authorization = *(authorization.(*auth.EnterpriseAuthorization))
	case conf.ENTITY_EMPLOYEE:
		baser = &auth.Employee{}
		authorization = &auth.EmployeeAuthorization{}
		baser.GetEntity().(*auth.Employee).Authorization = *(authorization.(*auth.EmployeeAuthorization))
	default:
		g.JSON(http.StatusBadGateway, &response.Response{
			Code:    response.Error,
			Message: i18n.I18NViper.GetString("message.auth.login.unsupported"),
			Data:    nil,
		})
		return
	}
	switch {
	case util.REG_MOBILE.Match([]byte(input.Account)):
		DB().Model(authorization).Where(auth.AuthBase{
			Mobile:   input.Account,
			Password: auth.Password(input.Password)}).First(authorization)
		fallthrough
	case util.REG_EMAIL.Match([]byte(input.Account)):
		if authorization.(auth.Authorizer).GetId() == "" {
			DB().Model(authorization).Where(auth.AuthBase{
				Email:    input.Account,
				Password: auth.Password(input.Password)}).First(authorization)
		}
		fallthrough
	default:
		if authorization.(auth.Authorizer).GetId() == "" {
			DB().Model(authorization).Where(auth.AuthBase{
				Username: input.Account,
				Password: auth.Password(input.Password)}).First(authorization)
		}
	}
	if authorization.(auth.Authorizer).GetId() == "" {
		g.JSON(http.StatusNotFound, &response.Response{
			Code:    response.Error,
			Message: i18n.I18NViper.GetString("message.auth.login.notfound"),
			Data:    nil,
		})
	} else {
		if DB().Model(baser.GetEntity()).
			Where(&model.Base{Id: authorization.(auth.Authorizer).GetId()}).
			First(baser.GetEntity()).RecordNotFound() {
			g.JSON(http.StatusNotFound, &response.Response{
				Code:    response.Error,
				Message: i18n.I18NViper.GetString("message.auth.login.notfound"),
				Data:    nil,
			})
		} else {
			switch entity {
			case conf.ENTITY_ENTERPRISE:
				baser.GetEntity().(*auth.Enterprise).Authorization = *(authorization.(*auth.EnterpriseAuthorization))
				baser.GetEntity().(*auth.Enterprise).Authorization.Password = ""
				if token, err := GenerateToken(baser.GetEntity().(*auth.Enterprise), jwt.StandardClaims{
					ExpiresAt: time.Now().Add(time.Duration(conf.AppConfig.Api.AuthTokenExpiration) * time.Hour).Unix(),
				}); err != nil {
					g.JSON(http.StatusInternalServerError, &response.Response{
						Code:    response.Error,
						Error:   fmt.Sprint(err),
						Message: i18n.I18NViper.GetString("message.common.internalerror"),
						Data:    nil,
					})
					return
				} else {
					baser.GetEntity().(*auth.Enterprise).Authorization.Jwt = token
				}
			case conf.ENTITY_EMPLOYEE:
				baser.GetEntity().(*auth.Employee).Authorization = *(authorization.(*auth.EmployeeAuthorization))
				baser.GetEntity().(*auth.Employee).Authorization.Password = ""
				if token, err := GenerateToken(baser.GetEntity().(*auth.Employee), jwt.StandardClaims{
					ExpiresAt: time.Now().Add(time.Duration(conf.AppConfig.Api.AuthTokenExpiration) * time.Hour).Unix(),
				}); err != nil {
					g.JSON(http.StatusInternalServerError, &response.Response{
						Code:    response.Error,
						Error:   fmt.Sprint(err),
						Message: i18n.I18NViper.GetString("message.common.internalerror"),
						Data:    nil,
					})
					return
				} else {
					baser.GetEntity().(*auth.Employee).Authorization.Jwt = token
				}
			}
			g.JSON(http.StatusOK, &response.Response{
				Code:    response.OK,
				Message: i18n.I18NViper.GetString("message.auth.login.ok"),
				Data:    baser,
			})
		}
	}
}

func GenerateToken(author auth.Author, s jwt.StandardClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, auth.CapabilityClaims{
		AuthType:       author.GetAuthType(),
		Capabilities:   []auth.Capability{},
		Id:             author.GetBase().(*model.Base).Id,
		OwnerId:        author.GetBase().(*model.Base).OwnerId,
		Type:           author.GetBase().(*model.Base).Type,
		StandardClaims: s,
	})
	if accessToken, err := token.SignedString(auth.AUTH_HMAC_SECRET); err != nil {
		return "", err
	} else {
		return accessToken, nil
	}
}

func GetTokenString(g *gin.Context) (tokenString string, err error) {
	tokenStringWithPrefix := g.Request.Header.Get("Authorization")
	tokenString = strings.TrimPrefix(tokenStringWithPrefix, "Bearer ")
	if tokenString == "" {
		tokenString = g.Request.URL.Query().Get("access_token")
	}
	if tokenString == "" {
		err = fmt.Errorf("no token found.")
	}
	return tokenString, err
}

func ParseToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.ParseWithClaims(tokenString, &auth.CapabilityClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return auth.AUTH_HMAC_SECRET, nil
	})
	if token != nil {
		if claims, ok := token.Claims.(*auth.CapabilityClaims); ok && token.Valid {
			Logger.Info("parsed ", claims.Id, claims)
		} else {
			Logger.Info("parse error: ", err)
		}
	}
	return token, err
}
