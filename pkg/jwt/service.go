package jwt

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

const (
	LongLiveTokenDuration = time.Hour * 8760 * 50 // 50 years
)

var (
	ErrWrongTokenSigningMethod = errors.New("unexpected token signing method")
	ErrInvalidToken            = errors.New("invalid token")
	ErrInvalidTokenClaims      = errors.New("invalid token claims")
)

type service struct {
	secret string
}

func (s *service) GetTokenData(accessToken string) (map[string]string, error) {
	token, err := jwt.ParseWithClaims(accessToken,
		&tokenClaim{},
		func(token *jwt.Token) (interface{}, error) {
			_, ok := token.Method.(*jwt.SigningMethodHMAC)
			if !ok {
				return nil, fmt.Errorf("%s, %w", "unsupported sign method", ErrWrongTokenSigningMethod)
			}

			return []byte(s.secret), nil
		},
	)
	if err != nil {
		return nil, ErrInvalidToken
	}

	if !token.Valid {
		return nil, ErrInvalidToken
	}

	claim, ok := token.Claims.(*tokenClaim)
	if !ok {
		return nil, ErrInvalidTokenClaims
	}

	return claim.GetAllData(), nil
}

func (s *service) GenerateJWT(claim tokenClaimBuilderService) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)

	return token.SignedString([]byte(s.secret))
}

func NewJWTService(secret string) (s *service) {
	s = &service{
		secret: secret,
	}
	return
}
