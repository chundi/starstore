package auth

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/galaxy-solar/starstore/conf"
)

var (
	AUTH_HMAC_SECRET []byte = []byte{101, 109, 104, 108, 98, 71, 108, 122, 97, 71, 108, 105, 97, 87, 70, 118, 90, 50, 85, 115, 99, 51, 82, 104, 99, 110, 78, 48, 98, 51, 74, 119, 99, 109, 57, 113, 90, 87, 78, 48, 76, 103, 111, 61}
)

type CapabilityClaims struct {
	AuthType string `json:"auth_type"`

	Capabilities []Capability `json:"capabilities,omitempty"`
	Id           string       `json:"id,omitempty"`
	OwnerId      string       `json:"owner_id,omitempty"`
	Type         string       `json:"type,omitempty"`

	jwt.StandardClaims
}

func (claim CapabilityClaims) IsEnterprise() bool {
	return claim.AuthType == conf.ENTITY_ENTERPRISE
}

func (claim CapabilityClaims) IsEmployee() bool {
	return claim.AuthType == conf.ENTITY_EMPLOYEE
}
