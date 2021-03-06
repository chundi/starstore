package model

import (
	"github.com/galaxy-solar/starstore/i18n"
	"github.com/galaxy-solar/starstore/model/feature"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/satori/go.uuid"
	"github.com/spf13/viper"
	"time"
)

type Baser interface {
	GetEntity() interface{}
	GetBase() interface{}
}

type Omiter interface {
	GetQueryOmittedFields() []string
	GetOrmOmittedFields() []string
}

type HandlerImplementer interface {
	ExecuteHandlers(g *gin.Context, db *gorm.DB, position TemplatePosition) error
	AddHandler(position TemplatePosition, handler BaseHandlerWithDB)
}

type EntityBaser interface {
	Baser
	Omiter
	HandlerImplementer
	GetMessage() *viper.Viper
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

	POSITION_DETAIL_GET_START
	POSITION_DETAIL_GET_AFTER
	POSITION_DETAIL_PUT_START
	POSITION_DETAIL_PUT_AFTER
	POSITION_DETAIL_DELETE_START
	POSITION_DETAIL_DELETE_AFTER
)

type Base struct {
	Id          string `sql:"type:uuid; not null; primary key" gorm:"primary_key" json:"id,omitempty"`
	OwnerId     string `sql:"type:uuid; default:'00000000-0000-0000-0000-000000000000'" json:"owner_id,omitempty"`
	ParentId    string `sql:"type:uuid; default:'00000000-0000-0000-0000-000000000000'" json:"parent_id,omitempty"`
	Type        string `binding:"required" json:"type,omitempty"`
	Status      string `json:"status,omitempty"`
	Token       string `json:"token,omitempty"`
	Name        string `binding:"required" json:"name,omitempty"`
	Slug        string `json:"slug,omitempty"`
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	Info        string `json:"info,omitempty"`
	Remark      string `json:"remark,omitempty"`

	Tag     feature.JSONB `sql:"type:jsonb" json:"tag,omitempty"`
	Content feature.JSONB `sql:"type:jsonb" json:"content,omitempty"`
	Meta    interface{}   `gorm:"-" json:"meta,omitempty"`

	CreatedDate   *time.Time `json:"created_date,omitempty"`
	UpdatedDate   *time.Time `json:"updated_date,omitempty"`
	DeletedDate   *time.Time `json:"deleted_date,omitempty"`
	PublishedDate *time.Time `json:"published_date,omitempty"`

	HandlerImplement
}

type HandlerImplement struct {
	Handlers map[TemplatePosition][]BaseHandlerWithDB `gorm:"-" json:"-"`
}

func (base *Base) GetBase() interface{} {
	return base
}

func (base Base) GetQueryOmittedFields() []string {
	return []string{"id", "parent_id", "owner_id", "created_date", "updated_date", "published_date", "deleted_date"}
}

func (base Base) GetOrmOmittedFields() []string {
	return []string{"id"}
}

func (base Base) GetMessage() *viper.Viper {
	return i18n.I18NViper.Sub("message.base")
}

func (base *Base) SetCreateDate(t time.Time) {
	base.CreatedDate = &t
}

func (base *Base) SetUpdateDate(t time.Time) {
	base.UpdatedDate = &t
}

func (implementer HandlerImplement) ExecuteHandlers(g *gin.Context, db *gorm.DB, position TemplatePosition) error {
	for _, handler := range implementer.Handlers[position] {
		if err := handler(g, db); err != nil {
			return err
		}
	}
	return nil
}

func (implementer *HandlerImplement) AddHandler(position TemplatePosition, handler BaseHandlerWithDB) {
	if implementer.Handlers == nil {
		implementer.Handlers = make(map[TemplatePosition][]BaseHandlerWithDB)
	}
	if _, ok := implementer.Handlers[position]; !ok {
		implementer.Handlers[position] = []BaseHandlerWithDB{handler}
	} else {
		implementer.Handlers[position] = append(implementer.Handlers[position], handler)
	}
}

func (base *Base) BeforeCreate(scope *gorm.Scope) error {
	scope.SetColumn("Id", uuid.NewV4().String())
	return nil
}
