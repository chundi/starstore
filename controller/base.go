package controller

import (
	"github.com/jinzhu/gorm"
	"github.com/galaxy-solar/starstore/model"
	"github.com/galaxy-solar/starstore/log"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"net/http"
	"github.com/galaxy-solar/starstore/response"
	"github.com/galaxy-solar/starstore/i18n"
	"fmt"
	"github.com/galaxy-solar/starstore/util"
	"io"
	"io/ioutil"
	"bytes"
	"encoding/json"
)

var Logger *logrus.Logger

func init() {
	Logger = log.InitDefaultLogger()
}

func DB() *gorm.DB {
	return model.DB.New()
}

func BaseMessage(baser model.Baser, key string) string {
	if baser.GetMessage().IsSet(key) {
		return baser.GetMessage().GetString(key)
	} else {
		return i18n.I18NViper.GetString(fmt.Sprintf("message.base.%s", key))
	}
}

func BindWithTeeReader(g *gin.Context, obj interface{}) error {
	b := bytes.NewBuffer(make([]byte, 0))
	reader := io.TeeReader(g.Request.Body, b)
	err := json.NewDecoder(reader).Decode(obj);
	g.Request.Body = ioutil.NopCloser(b)
	return err
}

func BasePost(g *gin.Context, db *gorm.DB, baser model.Baser) {
	BindWithTeeReader(g, baser.GetBase())
	if err := baser.ExecuteHandlers(g, db, model.POSITION_POST_BEFORE_CREATE); err != nil {
		return
	}
	tx := db.Begin()
	if err := baser.ExecuteHandlers(g, db, model.POSITION_POST_TRANSACTION_START); err != nil {
		return
	}

	if err := tx.Create(baser.GetEntity()).Error; err != nil {
		tx.Rollback()

		Logger.Error("base pose creating error: ", err)
		g.JSON(http.StatusInternalServerError, &response.Response{
			Code: response.Error,
			Message: BaseMessage(baser, "uncreated"),
		})
	} else {
		if err := baser.ExecuteHandlers(g, tx, model.POSITION_POST_TRANSACTION_END); err != nil {
			return
		}
		tx.Commit()
		if err := baser.ExecuteHandlers(g, tx, model.POSITION_POST_AFTER_CREATE); err != nil {
			return
		}

		Logger.Info(fmt.Sprintf("created %s id: %s", util.ModelType(baser), baser.GetId()))
		g.JSON(http.StatusCreated, &response.Response{
			Code: response.OK,
			Message: BaseMessage(baser, "created"),
			Data: baser,
		})
	}
}
