package model

import (
	"github.com/jinzhu/gorm"
)

func GetOneById(db *gorm.DB, id string, entity Baser) bool {
	r := db.Where("id = ?", id).First(entity.GetEntity())
	if r.Error != nil {
		flog := Logger.WithField("id", id).
			WithField("entity", entity)
		if r.RecordNotFound() {
			flog.Error("Record Not Found!")
		} else {
			flog.Error(r.Error.Error())
		}
		return false
	}
	return true
}
