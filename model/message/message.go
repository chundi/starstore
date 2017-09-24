package message

import (
	"github.com/galaxy-solar/starstore/model/feature"
	"time"
	"github.com/galaxy-solar/starstore/model"
	"github.com/jinzhu/gorm"
)

type Message struct {
	Entity      string        `json:"entity"`
	EntityId    string        `sql:"type:uuid; not null; primary key" gorm:"primary_key" json:"entity_id"`
	Type        string        `json:"type,omitempty"`
	CreatedDate time.Time     `sql:"DEFAULT:current_timestamp"`
	Content     feature.JSONB `sql:"type:jsonb" json:"content,omitempty"`
}

func AddMessage(m *Message) bool {
	return model.DB.NewRecord(m)
}

func UpdateMessage(d string, m *Message) *gorm.DB {
	sql := ``
	return model.DB.Exec(sql)
}
