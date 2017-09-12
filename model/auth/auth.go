package auth

import (
	"database/sql/driver"
	"time"
)

type Password string

func (u *Password) Scan(value interface{}) error {
	*u = Password(value.(string));
	return nil
}

func (u Password) Value() (driver.Value, error) {
	return string(u), nil
}

type AuthOAuth2 struct {
	// OAuth2
	Oauth2Uid      string
	Oauth2Provider string
	Oauth2Token    string
	Oauth2Refresh  string
	Oauth2Expiry   time.Time
}

type AuthConfirm struct {
	// Confirm
	ConfirmToken string
	Confirmed    bool
}

type AuthLock struct {
	// Lock
	AttemptNumber int64
	AttemptTime   time.Time
	Locked        time.Time
}

type AuthRecover struct {
	// Recover
	RecoverToken       string
	RecoverTokenExpiry time.Time
}


type AuthBase struct {
	Username	string	`gorm:"not null;unique" json:"username"`
	Email 		string	`gorm:"not null;unique" json:"email"`
	Mobile  	string	`gorm:"not null;unique" json:"mobile"`
	Password	Password	`binding:"required" json:"password"`

	AuthConfirm
	AuthRecover
	AuthLock
}

type UserAuthorization struct {
	UserId string `sql:"type:uuid;not null"`

	AuthBase
}

type EnterpriseAuthorization struct {
	EnterpriseId string `sql:"type:uuid; not null"`

	AuthBase
	AuthOAuth2
}