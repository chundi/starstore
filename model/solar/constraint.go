package solar

import "github.com/galaxy-solar/starstore/model/feature"

type Constraint struct {
	ConstraintId feature.Uuid `sql:"type:uuid;not null";gorm:"primary_key"`
}
