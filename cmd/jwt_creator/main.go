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
	merchantUUIDLabel = "merchant_uuid"
	expiredAtLabel    = "expired_at"
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

	jwtSvc := jwttool.NewJWTService(key)
	token, err := jwtSvc.GenerateJWT(mExpTime, map[string]string{
		merchantUUIDLabel: uuid,
		expiredAtLabel:    mExpTime.Format(time.DateTime),
	})
	if err != nil {
		log.Fatalf("cant make JWT token. Error: %v ", err.Error())
	}

	log.Println("Token: ", token)

	data, err := jwtSvc.GetTokenData(token)
	if err != nil {
		log.Fatalf("unable to get data from JWT token. Error: %v ", err.Error())
	}

	log.Println("Origin data from token: ", data[merchantUUIDLabel], data[expiredAtLabel])
}
