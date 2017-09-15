package controller

import (
	"github.com/galaxy-solar/starstore/log"
	"github.com/galaxy-solar/starstore/model"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"

	"bytes"
	"encoding/json"
	"fmt"
	"github.com/galaxy-solar/starstore/conf"
	"github.com/galaxy-solar/starstore/i18n"
	"github.com/galaxy-solar/starstore/response"
	"github.com/galaxy-solar/starstore/util"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

var Logger *logrus.Logger

func init() {
	Logger = log.InitDefaultLogger()
}

func DB() *gorm.DB {
	db := model.DB.New()
	if conf.IsDevelopMode() {
		db.LogMode(true)
	}
	return db
}

func BaseMessage(baser model.Baser, key string) string {
	if baser.GetMessage().IsSet(key) {
		return baser.GetMessage().GetString(key)
	} else {
		return i18n.I18NViper.GetString(fmt.Sprintf("message.base.%s", key))
	}
}

func BindRequestBodyWithTeeReader(g *gin.Context, obj interface{}) error {
	b := bytes.NewBuffer(make([]byte, 0))
	reader := io.TeeReader(g.Request.Body, b)
	err := json.NewDecoder(reader).Decode(obj)
	g.Request.Body = ioutil.NopCloser(b)
	return err
}

func ParseFilteredQuery(g *gin.Context, db *gorm.DB, baser model.Baser) *gorm.DB {
	if escapedQuery, err := url.QueryUnescape(g.Request.URL.RawQuery); err != nil {
		g.JSON(http.StatusBadRequest, &response.Response{
			Code:    response.Error,
			Message: BaseMessage(baser, "error"),
			Data:    nil,
		})
	} else {
		db = util.ParseQuery(db, escapedQuery)
	}
	return db
}

func BaseGet(g *gin.Context, db *gorm.DB, baser model.Baser, baserList interface{}) {
	db = ParseFilteredQuery(g, db, baser)
	if err := baser.ExecuteHandlers(g, db, model.POSITION_GET_BEFORE_LIST); err != nil {
		return
	}
	if errs := db.Model(baser.GetEntity()).Scopes(BaseAvailable).Find(baserList).GetErrors(); len(errs) > 0 {
		Logger.Info(fmt.Sprintf("BaseGet error! %v", errs))
		g.JSON(http.StatusNotFound, &response.Response{
			Code:    response.NotFound,
			Error:   fmt.Sprint(errs),
			Message: BaseMessage(baser, "listerror"),
			Data:    nil,
		})
		return
	} else {
		if err := baser.ExecuteHandlers(g, db, model.POSITION_GET_AFTER_LIST); err != nil {
			return
		}
		Logger.Info(fmt.Sprintf("BaseGet ok!"))
		g.JSON(http.StatusOK, &response.Response{
			Code:    response.OK,
			Message: BaseMessage(baser, "ok"),
			Data:    baserList,
		})
	}
}

func BasePost(g *gin.Context, db *gorm.DB, baser model.Baser) {
	BindRequestBodyWithTeeReader(g, baser.GetEntity())
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
			Code:    response.Error,
			Message: BaseMessage(baser, "uncreated"),
			Data:    nil,
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
			Code:    response.OK,
			Message: BaseMessage(baser, "created"),
			Data:    baser,
		})
	}
}

func BaseDetailGet(g *gin.Context, db *gorm.DB, baser model.Baser) {
	db = ParseFilteredQuery(g, db, baser)
	id := g.Param("id")
	if err := baser.ExecuteHandlers(g, db, model.POSITION_DETAIL_GET_START); err != nil {
		return
	}
	if db.Where("id = ?", id).First(baser.GetEntity()).RecordNotFound() {
		g.JSON(http.StatusNotFound, &response.Response{
			Code:    response.NotFound,
			Message: BaseMessage(baser, "resourcenotfound"),
			Data:    nil,
		})
	} else {
		if err := baser.ExecuteHandlers(g, db, model.POSITION_DETAIL_GET_AFTER); err != nil {
			return
		}
		g.JSON(http.StatusOK, &response.Response{
			Code:    response.OK,
			Message: BaseMessage(baser, "ok"),
			Data:    baser,
		})
	}
}

func BaseDetailPut(g *gin.Context, db *gorm.DB, baser model.Baser) {
	omittedFields := []string{"id", "parent_id", "owner_id"}
	db = ParseFilteredQuery(g, db, baser)
	id := g.Param("id")
	if db.Where("id = ?", id).First(baser.GetEntity()).RecordNotFound() {
		if err := baser.ExecuteHandlers(g, db, model.POSITION_DETAIL_PUT_START); err != nil {
			return
		}
		g.JSON(http.StatusNotFound, &response.Response{
			Code:    response.NotFound,
			Message: BaseMessage(baser, "resourcenotfound"),
			Data:    nil,
		})
	} else {
		var givenJson map[string]interface{}
		BindRequestBodyWithTeeReader(g, &givenJson)
		for _, field := range omittedFields {
			if _, ok := givenJson[field]; ok {
				g.JSON(http.StatusBadGateway, &response.Response{
					Code:    response.Error,
					Message: BaseMessage(baser, "omitted"),
					Data:    nil,
				})
				return
			}
		}
		BindRequestBodyWithTeeReader(g, baser.GetEntity())

		if err := db.Model(baser.GetEntity()).Omit(omittedFields...).Updates(baser.GetEntity()).Error; err != nil {
			g.JSON(http.StatusBadGateway, &response.Response{
				Code:    response.Error,
				Error:   fmt.Sprint(err),
				Message: BaseMessage(baser, "updateerror"),
				Data:    nil,
			})
		} else {
			if err := baser.ExecuteHandlers(g, db, model.POSITION_DETAIL_PUT_AFTER); err != nil {
				return
			}
			g.JSON(http.StatusOK, &response.Response{
				Code:    response.OK,
				Message: BaseMessage(baser, "ok"),
				Data:    baser,
			})
		}
	}
}

func BaseDetailDelete(g *gin.Context, db *gorm.DB, baser model.Baser) {
	id := g.Param("id")
	if err := baser.ExecuteHandlers(g, db, model.POSITION_DETAIL_GET_START); err != nil {
		return
	}
	if db.Where("id = ?", id).First(baser.GetEntity()).RecordNotFound() {
		g.JSON(http.StatusNotFound, &response.Response{
			Code:    response.NotFound,
			Error:   fmt.Sprint(db.GetErrors()),
			Message: BaseMessage(baser, "resourcenotfound"),
			Data:    nil,
		})
	} else {
		if err := db.Model(baser.GetEntity()).Omit("id", "parent_id", "owner_id").Update("DeletedDate", time.Now()).Error; err != nil {
			g.JSON(http.StatusBadGateway, &response.Response{
				Code:    response.Error,
				Error:   fmt.Sprint(err),
				Message: BaseMessage(baser, "deleteerror"),
				Data:    nil,
			})
		} else {
			if err := baser.ExecuteHandlers(g, db, model.POSITION_DETAIL_GET_AFTER); err != nil {
				return
			}
			g.JSON(http.StatusOK, &response.Response{
				Code:    response.OK,
				Message: BaseMessage(baser, "deleteok"),
				Data:    nil,
			})
		}
	}
}
