package main

import (
	"errors"
	"flag"
	"log"
	"time"

	jwttool "github.com/crypto-bundle/bc-wallet-common-lib-jwt/pkg/jwt"

	_ "github.com/mailru/easyjson/gen"
)

var (
	ErrWrongDateFormat = errors.New("wrong date format")
)

const (
	MerchantUUIDLabel = "merchant_uuid"
	ExpiredAtLabel    = "expired_at"
)

func main() {
	var (
		key, uuid, expiration string
	)

	flag.StringVar(&key, "key", "", "secret key")
	flag.StringVar(&uuid, "uuid", "", "merchant uuid")
	flag.StringVar(&expiration, "expiration", "", "expiration date with format: '2006-01-02'")
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

	claimBuilder := jwttool.NewTokenClaimBuilder(mExpTime)
	err := claimBuilder.AddData(MerchantUUIDLabel, uuid)
	if err != nil {
		log.Fatalf("%s", err)
	}

	err = claimBuilder.AddData(ExpiredAtLabel, mExpTime.String())
	if err != nil {
		log.Fatalf("%s", err)
	}

	jwtSvc := jwttool.NewJWTService(key)
	token, err := jwtSvc.GenerateJWT(claimBuilder)
	if err != nil {
		log.Fatalf("cant make JWT token. Error: %v ", err.Error())
	}

	log.Println("Token: ", token)

	data, err := jwtSvc.GetTokenData(token)
	if err != nil {
		log.Fatalf("unable to get data from JWT token. Error: %v ", err.Error())
	}

	log.Println("Origin data from token: ", data[MerchantUUIDLabel], data[ExpiredAtLabel])
}
