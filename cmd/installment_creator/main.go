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

package main

import (
	"crypto/sha256"
	"errors"
	"flag"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/crypto-bundle/bc-wallet-common-lib-jwt/internal/mockerrors"
	jwttool "github.com/crypto-bundle/bc-wallet-common-lib-jwt/pkg/jwt"

	_ "github.com/mailru/easyjson/gen"
)

var (
	ErrWrongDateFormat = errors.New("wrong date format")
)

const (
	installmentUUIDLabel = "installment_uuid"
	tokenUUIDLabel       = "token_uuid"
	expiredAtLabel       = "expired_at"
	domainsLabel         = "domains"
)

type installmentTokenInfo struct {
	uuid, tokenUUID string
	expiration      time.Time
	domains         []string
}

func (s *installmentTokenInfo) ScanClaims(values jwttool.TokenClaimValues) error {
	var expiredAt time.Time

	err := values.ScanByKey(expiredAtLabel, &expiredAt)
	if err != nil {
		log.Fatalf("unable to scan data from JWT token. Error: %v ", err.Error())
	}

	var tokenUUID string

	err = values.ScanByKey(tokenUUIDLabel, &tokenUUID)
	if err != nil {
		log.Fatalf("unable to scan data from JWT token. Error: %v ", err.Error())
	}

	var instUUID string

	err = values.ScanByKey(installmentUUIDLabel, &instUUID)
	if err != nil {
		log.Fatalf("unable to scan data from JWT token. Error: %v ", err.Error())
	}

	var domainsResultList []string = nil

	err = values.ScanByKey(domainsLabel, &domainsResultList)
	if err != nil {
		log.Fatalf("unable to scan data from JWT token. Error: %v ", err.Error())
	}

	s.domains = domainsResultList
	s.tokenUUID = tokenUUID
	s.uuid = instUUID
	s.expiration = expiredAt

	return nil
}

//nolint:funlen
func main() {
	var (
		key, uuid, tokenUUID, expiration, domains string
	)

	flag.StringVar(&key, "key", "", "secret key")
	flag.StringVar(&uuid, "uuid", "", "installment uuid")
	flag.StringVar(&tokenUUID, "token_uuid", "", "token uuid")
	flag.StringVar(&expiration, "expiration", "", "expiration date with format: '2006-01-02'")
	flag.StringVar(
		&domains,
		"domains",
		"",
		"domains list separated by comma - domain.com,sub.domain.com",
	)
	flag.Parse()

	if key == "" {
		log.Fatalf("missing required -%v argument/flag\n", "key")
	}

	mExpTime := time.Now().Add(jwttool.LongLiveTokenDuration)

	if expiration != "" {
		expTime, innerErr := time.Parse("2006-01-02", expiration)
		if innerErr != nil {
			log.Fatalf("%s: %s", ErrWrongDateFormat, innerErr)
		}

		mExpTime = expTime
	}

	jwtSvc := jwttool.NewJWTManger(mockerrors.NewMockErrFormatter(), key)

	domainsList := strings.Split(domains, ",")

	token, err := jwtSvc.GenerateJWT(mExpTime,
		jwttool.String(installmentUUIDLabel, uuid),
		jwttool.String(tokenUUIDLabel, tokenUUID),
		jwttool.Time(expiredAtLabel, mExpTime),
		jwttool.Strings(domainsLabel, domainsList),
	)
	if err != nil {
		log.Fatalf("cant make JWT token. Error: %v ", err.Error())
	}

	log.Println("Token: ", token)

	data, err := jwtSvc.GetTokenClaimsData(token)
	if err != nil {
		log.Fatalf("unable to get data from JWT token. Error: %v ", err.Error())
	}

	var expitedAt time.Time

	err = data.ScanByKey(expiredAtLabel, &expitedAt)
	if err != nil {
		log.Fatalf("unable to scan data from JWT token. Error: %v ", err.Error())
	}

	var instUUID string

	err = data.ScanByKey(installmentUUIDLabel, &instUUID)
	if err != nil {
		log.Fatalf("unable to scan data from JWT token. Error: %v ", err.Error())
	}

	var domainsResultList []string = nil

	err = data.ScanByKey(domainsLabel, &domainsResultList)
	if err != nil {
		log.Fatalf("unable to scan data from JWT token. Error: %v ", err.Error())
	}

	fullTokenInfo := installmentTokenInfo{}

	err = data.ScanToStruct(&fullTokenInfo)
	if err != nil {
		log.Fatalf("unable to scan data from JWT token to struct. Error: %v ", err.Error())
	}

	log.Println("Origin data from token: ", instUUID, expitedAt, domainsResultList)
	log.Println(
		"Scanned data from token: ",
		fullTokenInfo.uuid,
		fullTokenInfo.expiration,
		fullTokenInfo.domains,
	)
	log.Println("Token hash: ", fmt.Sprintf("%x", sha256.Sum256([]byte(token))))
}
