package solar

import "github.com/galaxy-solar/starstore/model/feature"

type TypeIndicator struct {
	TypeId feature.Uuid `sql:"type:uuid; not null";gorm:"primary_key"`

	EntityType string	`gorm:"unique_index:idx_type_indicator_entity"`
	EntityId   feature.Uuid `sql:"type:uuid; not null";gorm:"unique_index:idx_type_indicator_entity"`
	TypeToken  string	`gorm:"unique_index:idx_type_indicator_entity"`

	ContentExplanation feature.JSONB	`sql:"type:jsonb;not null"`
}
