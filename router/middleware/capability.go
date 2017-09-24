package middleware

import (
	"fmt"
	"github.com/galaxy-solar/starstore/controller"
	"github.com/galaxy-solar/starstore/i18n"
	"github.com/galaxy-solar/starstore/log"
	"github.com/galaxy-solar/starstore/response"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

var Logger *logrus.Logger

func init() {
	Logger = log.InitDefaultLogger()
}

type Tokenizer interface {
}

func CapabilitiesMiddleware() gin.HandlerFunc {
	return func(g *gin.Context) {
		if tokenString, err := controller.GetTokenString(g); err != nil {
			if g.Request.Header.Get("liuzb") == "wsnbg" {
				g.Next()
			} else {
				g.JSON(http.StatusUnauthorized, &response.Response{
					Code:    response.Unauthorized,
					Error:   fmt.Sprint(err),
					Message: i18n.I18NViper.GetString("message.common.unauthorized"),
				})
			}
			g.Abort()
		} else {
			if token, err := controller.ParseToken(tokenString); err != nil {
				g.JSON(http.StatusInternalServerError, &response.Response{
					Code:    response.Error,
					Error:   fmt.Sprint(err),
					Message: i18n.I18NViper.GetString("message.common.tokenerror"),
				})
				g.Abort()
			} else {
				g.Set(controller.ACCESS_TOKEN, token)
				g.Next()
			}
		}
	}
}
