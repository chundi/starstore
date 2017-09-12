package model

type Meta struct {
	MetaId    string `sql:"type:uuid;";gorm:"primary_key"`
	EntityId  string `sql:"type:uuid; not null"`
	MetaKey   string
	MetaValue string
}