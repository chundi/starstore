package solar

import "github.com/galaxy-solar/starstore/model/feature"

type Widget struct {
	WidgetId feature.Uuid	`sql:"type:uuid; not null";gorm:"primary_key"`
	
}
