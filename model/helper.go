package model

import (
	"github.com/jinzhu/gorm"
)

func GetOneById(db *gorm.DB, id string, entity Baser) bool {
	r := db.Where("id = ?", id).First(entity.GetEntity())
	if r.Error != nil {
		if r.RecordNotFound() {
			logger.WithField("id", id).
				WithField("entity", entity).
				Error("Record Not Found!")
		} else {
			logger.WithField("id", id).
				WithField("entity", entity).
				Error(r.Error.Error())
		}
		return false
	}
	return true
}
