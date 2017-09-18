package auth

import (
	"database/sql/driver"
	"time"
)

type Password string

func (u *Password) Scan(value interface{}) error {
	*u = Password(value.(string))
	return nil
}

func (u Password) Value() (driver.Value, error) {
	return string(u), nil
}

type AuthOAuth2 struct {
	// OAuth2
	Oauth2Uid      string     `json:"oauth2uid,omitempty"`
	Oauth2Provider string     `json:"oauth2provider,omitempty"`
	Oauth2Token    string     `json:"oauth2token,omitempty"`
	Oauth2Refresh  string     `json:"oauth2refresh,omitempty"`
	Oauth2Expiry   *time.Time `json:"oauth2expiry,omitempty"`
}

type AuthConfirm struct {
	// Confirm
	ConfirmToken string `json:"confirm_token,omitempty"`
	Confirmed    bool   `json:"confirmed,omitempty"`
}

type AuthLock struct {
	// Lock
	AttemptNumber int64      `json:"attempt_number,omitempty"`
	AttemptTime   *time.Time `json:"attempt_time,omitempty"`
	Locked        *time.Time `json:"locked,omitempty"`
}

type AuthRecover struct {
	// Recover
	RecoverToken       string     `json:"recover_token,omitempty"`
	RecoverTokenExpiry *time.Time `json:"recover_token_expiry,omitempty"`
}

type AuthBase struct {
	Username string   `gorm:"not null;unique" json:"username,omitempty"`
	Email    string   `gorm:"not null;unique" json:"email,omitempty"`
	Mobile   string   `gorm:"not null;unique" json:"mobile,omitempty"`
	Password Password `binding:"required" json:"password,omitempty"`

	AuthConfirm
	AuthRecover
	AuthLock
}

type UserAuthorization struct {
	Id string `sql:"type:uuid;not null" json:"id,omitempty"`

	AuthBase
}

type EnterpriseAuthorization struct {
	Id string `sql:"type:uuid; not null" json:"id,omitempty"`

	AuthBase
	AuthOAuth2
}
