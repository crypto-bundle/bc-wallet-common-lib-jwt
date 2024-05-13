package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

//go:generate easyjson types.go

var (
	ErrDataPartAlreadyExists = errors.New("data already exist in claim")
)

// tokenClaim for store data map of uuid's
// easyjson:json
type tokenClaim struct {
	register  jwt.RegisteredClaims
	ValuesMap map[string]string `json:"values_map"`
}

func (c *tokenClaim) Valid() error {
	return c.register.Valid()
}

func (c *tokenClaim) GetAllData() map[string]string {
	return c.ValuesMap
}

func (c *tokenClaim) AddData(dataLabel string, dataValue string) error {
	_, isExists := c.ValuesMap[dataLabel]
	if isExists {
		return ErrDataPartAlreadyExists
	}

	c.ValuesMap[dataLabel] = dataValue

	return nil
}

func NewTokenClaimBuilder(expiredAt time.Time) *tokenClaim {
	return &tokenClaim{
		register: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiredAt),
		},
		ValuesMap: make(map[string]string, 1),
	}
}
