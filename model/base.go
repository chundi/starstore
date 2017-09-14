package model

import (
	"time"
	"github.com/galaxy-solar/starstore/model/feature"
	"github.com/jinzhu/gorm"
	"github.com/satori/go.uuid"
	"github.com/spf13/viper"
	"github.com/gin-gonic/gin"
	"github.com/galaxy-solar/starstore/i18n"
)

type Baser interface {
	GetBase() *Base
	GetEntity() interface{}
	GetId() string
	GetMessage() *viper.Viper
	ExecuteHandlers(*gin.Context, *gorm.DB, TemplatePosition) error
}

type BaseHandlerWithDB func(g *gin.Context, db *gorm.DB) error
type TemplatePosition int

const (
	POSITION_POST_BEFORE_CREATE TemplatePosition = iota
	POSITION_POST_TRANSACTION_START
	POSITION_POST_TRANSACTION_END
	POSITION_POST_AFTER_CREATE

	POSITION_GET_BEFORE_LIST
	POSITION_GET_AFTER_LIST
)

type Base struct {
	Id          string `sql:"type:uuid; not null" gorm:"primary_key" json:"id,omitempty"`
	OwnerId     *string `sql:"type:uuid; default:'00000000-0000-0000-0000-000000000000'" json:"owner_id,omitempty"`
	ParentId    *string `sql:"type:uuid; default:'00000000-0000-0000-0000-000000000000'" json:"parent_id,omitempty"`
	Type        string	`binding:"required" json:"type,omitempty"`
	Status      string	`json:"status,omitempty"`
	Name        string	`binding:"required" json:"name,omitempty"`
	Slug        string	`json:"slug,omitempty"`
	Title       string	`json:"title,omitempty"`
	Description string	`json:"description,omitempty"`

	Content feature.JSONB `sql:"type:jsonb" json:"content,omitempty"`

	CreatedDate   *time.Time	`json:"created_date,omitempty"`
	UpdatedDate   *time.Time `json:"updated_date,omitempty"`
	DeletedDate   *time.Time	`json:"deleted_date,omitempty"`
	PublishedDate *time.Time	`json:"published_date,omitempty"`

	Handlers map[TemplatePosition] []BaseHandlerWithDB `gorm:"-" json:"-"`
}

func (base Base) GetId() string {
	return base.Id
}

func (base Base) GetMessage() *viper.Viper {
	return i18n.I18NViper.Sub("message.base")
}

func (base Base) ExecuteHandlers(g *gin.Context, db *gorm.DB, position TemplatePosition) error {
	for _, handler := range base.Handlers[position] {
		if err := handler(g, db); err != nil {
			return err
		}
	}
	return nil
}

func (base *Base) GetBase() *Base {
	return base
}

func (base *Base) AddHandler(position TemplatePosition, handler BaseHandlerWithDB) {
	if base.Handlers == nil {
		base.Handlers = make(map[TemplatePosition] []BaseHandlerWithDB)
	}
	if _, ok := base.Handlers[position]; !ok {
		base.Handlers[position] = []BaseHandlerWithDB{handler}
	} else {
		base.Handlers[position] = append(base.Handlers[position], handler)
	}
}

func (base *Base) BeforeCreate(scope *gorm.Scope) error {
	scope.SetColumn("Id", uuid.NewV4().String())
	return nil
}
