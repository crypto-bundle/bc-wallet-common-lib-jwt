/*
 *
 *
 * MIT NON-AI License
 *
 * Copyright (c) 2022-2025 Aleksei Kotelnikov(gudron2s@gmail.com)
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy of the software and associated documentation files (the "Software"),
 * to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense,
 * and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions.
 *
 * The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
 *
 * In addition, the following restrictions apply:
 *
 * 1. The Software and any modifications made to it may not be used for the purpose of training or improving machine learning algorithms,
 * including but not limited to artificial intelligence, natural language processing, or data mining. This condition applies to any derivatives,
 * modifications, or updates based on the Software code. Any usage of the Software in an AI-training dataset is considered a breach of this License.
 *
 * 2. The Software may not be included in any dataset used for training or improving machine learning algorithms,
 * including but not limited to artificial intelligence, natural language processing, or data mining.
 *
 * 3. Any person or organization found to be in violation of these restrictions will be subject to legal action and may be held liable
 * for any damages resulting from such use.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM,
 * DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE
 * OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
 *
 */

package jwt

import (
	"errors"
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
	e errorFormatterService

	secret string
}

func (s *service) GetTokenData(accessToken string) (map[string]string, error) {
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
		return nil, s.e.ErrorOnly(err)
	}

	if !token.Valid {
		return nil, s.e.ErrorOnly(ErrInvalidToken)
	}

	claim, ok := token.Claims.(*tokenClaim)
	if !ok {
		return nil, s.e.ErrorOnly(ErrInvalidTokenClaims)
	}

	return claim.GetAllData(), nil
}

func (s *service) GenerateJWT(expiredAt time.Time, values map[string]string) (string, error) {
	claimBuilder := newTokenClaimBuilder(s.e, expiredAt)
	claimBuilder.SetValues(values)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claimBuilder)

	signed, err := token.SignedString([]byte(s.secret))
	if err != nil {
		return "", s.e.ErrorOnly(err)
	}

	return signed, nil
}

func NewJWTService(errFmtSvc errorFormatterService,
	secret string,
) *service {
	return &service{
		e:      errFmtSvc,
		secret: secret,
	}
}
