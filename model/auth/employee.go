package auth

import (
	"github.com/galaxy-solar/starstore/conf"
	"github.com/galaxy-solar/starstore/i18n"
	"github.com/galaxy-solar/starstore/model"
	"github.com/spf13/viper"
)

type Employee struct {
	model.Base

	Authorization EmployeeAuthorization `gorm:"-" json:"authorization,omitempty"`
}

func (employee Employee) GetAuthorization() Authorizer {
	return employee.Authorization
}

func (employee Employee) GetAuthType() string {
	return conf.ENTITY_EMPLOYEE
}

func (employee Employee) GetMessage() *viper.Viper {
	return i18n.I18NViper.Sub("message.auth.employee")
}

func (employee *Employee) GetEntity() interface{} {
	return employee
}

type EmployeeMeta struct {
	model.Meta
}
