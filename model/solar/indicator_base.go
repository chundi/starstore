package solar

import (
	"github.com/galaxy-solar/starstore/model"
	"github.com/galaxy-solar/starstore/model/feature"
	"time"
)

type IndicateBaser interface {
	model.Baser
	model.Omiter
	model.HandlerImplementer
}

type IndicatorBase struct {
	Id       string  `sql:"type:uuid; not null; primary key" gorm:"primary_key"`
	ParentId *string `sql:"type:uuid; default:'00000000-0000-0000-0000-000000000000'" json:"parent_id,omitempty"`

	Version int

	Token       string `gorm:"not null; unique"`
	Name        string `gorm:"not null; unique"`
	Title       string
	Description string

	ContentExplanation feature.JSONB `sql:"type:jsonb;not null"`

	CreatedDate   *time.Time `json:"created_date,omitempty"`
	UpdatedDate   *time.Time `json:"updated_date,omitempty"`
	DeletedDate   *time.Time `json:"deleted_date,omitempty"`
	PublishedDate *time.Time `json:"published_date,omitempty"`

	model.HandlerImplement
}

func (indicator *IndicatorBase) GetBase() interface{} {
	return indicator
}

func (indicator IndicatorBase) GetQueryOmittedFields() []string {
	return []string{}
}

func (indicator IndicatorBase) GetOrmOmittedFields() []string {
	return []string{}
}
