// MIT NON-AI License
//
// Copyright (c) 2022-2025 Aleksei Kotelnikov(gudron2s@gmail.com)
//
// Permission is hereby granted, free of charge, to any person obtaining a copy of the software and associated
// documentation files (the "Software"),to deal in the Software without restriction, including without limitation
// the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software,
// and to permit persons to whom the Software is furnished to do so, subject to the following conditions.
//
// The above copyright notice and this permission notice shall be included in all copies or substantial
// portions of the Software.
//
// In addition, the following restrictions apply:
//
// 1. The Software and any modifications made to it may not be used for the purpose of training or improving machine
// learning algorithms, including but not limited to artificial intelligence, natural language processing,
// or data mining.This condition applies to any derivatives, modifications, or updates based on the
// Software code. Any usage of the Software in an AI-training dataset is considered a breach of this License.
//
// 2. The Software may not be included in any dataset used for training or improving machine learning algorithms,
// including but not limited to artificial intelligence, natural language processing, or data mining.
//
// 3. Any person or organization found to be in violation of these restrictions will be subject to legal action and
// may be held liable for any damages resulting from such use.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO
// THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT,
// TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER
// DEALINGS IN THE SOFTWARE.

package jwt

import (
	"errors"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
)

const (
	LongLiveTokenDuration = time.Hour * 8760 * 50 // 50 years
)

var (
	ErrWrongTokenSigningMethod = errors.New("unexpected token signing method")
	ErrInvalidToken            = errors.New("invalid token")
	ErrInvalidTokenClaims      = errors.New("invalid token claims")
)

type TokenManager struct {
	e errorFormatterService

	secret string
}

func (s *TokenManager) GetTokenClaimsData(accessToken string) (TokenClaimValues, error) {
	_, claimData, err := s.decodeToken(accessToken)
	if err != nil {
		return nil, err
	}

	return claimData, nil
}

func (s *TokenManager) ValidateToken(accessToken string) (bool, error) {
	token, _, err := s.decodeToken(accessToken)
	if err != nil {
		return false, err
	}

	return token.Valid, nil
}

func (s *TokenManager) DecodeToken(accessToken string) (*jwt.Token, TokenClaimValues, error) {
	return s.decodeToken(accessToken)
}

func (s *TokenManager) decodeToken(accessToken string) (*jwt.Token, TokenClaimValues, error) {
	token, err := jwt.ParseWithClaims(accessToken, &tokenClaim{
		e: nil,
		register: jwt.RegisteredClaims{
			Issuer:    "",
			Subject:   "",
			Audience:  nil,
			ExpiresAt: nil,
			NotBefore: nil,
			IssuedAt:  nil,
			ID:        "",
		},
		ValuesMap: nil,
	}, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, s.e.ErrorOnly(ErrWrongTokenSigningMethod)
		}

		return []byte(s.secret), nil
	})
	if err != nil {
		return nil, nil, s.e.ErrorOnly(err)
	}

	claim, ok := token.Claims.(*tokenClaim)
	if !ok {
		return nil, nil, s.e.ErrorOnly(ErrInvalidTokenClaims)
	}

	return token, claim.GetAllData(), nil
}

func (s *TokenManager) GenerateJWT(expiredAt time.Time, values ...Field) (string, error) {
	claimBuilder := newTokenClaimBuilder(s.e, expiredAt)

	err := claimBuilder.AddValues(values...)
	if err != nil {
		return "", s.e.Error(err, "unable to add values to token claim")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claimBuilder)

	signed, err := token.SignedString([]byte(s.secret))
	if err != nil {
		return "", s.e.ErrorOnly(err, "unable to get signed string")
	}

	return signed, nil
}

func NewJWTManger(errFmtSvc errorFormatterService,
	secret string,
) *TokenManager {
	return &TokenManager{
		e:      errFmtSvc,
		secret: secret,
	}
}
