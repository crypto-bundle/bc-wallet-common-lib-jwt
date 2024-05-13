package jwt

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

//go:generate easyjson types.go

var (
	ErrDataPartAlreadyExists = errors.New("data already exist in claim")
)

// tokenClaim for store data map of uuid's
// easyjson:json
type tokenClaim struct {
	register jwt.RegisteredClaims
	UUIDMap  map[string]uuid.UUID `json:"uuid_map"`
}

func (c *tokenClaim) Valid() error {
	return c.register.Valid()
}

func (c *tokenClaim) GetAllData() map[string]uuid.UUID {
	return c.UUIDMap
}

func (c *tokenClaim) AddData(dataLabel string, dataUUID string) error {
	_, isExists := c.UUIDMap[dataLabel]
	if isExists {
		return ErrDataPartAlreadyExists
	}

	dataUUIDRaw, err := uuid.Parse(dataUUID)
	if err != nil {
		return fmt.Errorf("%s, %w", "wrong merchant uuid format", err)
	}

	c.UUIDMap[dataLabel] = dataUUIDRaw

	return nil
}

func NewTokenClaimBuilder(expiredAt time.Time) *tokenClaim {
	return &tokenClaim{
		register: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiredAt),
		},
		UUIDMap: make(map[string]uuid.UUID, 1),
	}
}
