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

func main() {
	var (
		key, uuid, tokenUUID, expiration, domains string
	)

	flag.StringVar(&key, "key", "", "secret key")
	flag.StringVar(&uuid, "uuid", "", "installment uuid")
	flag.StringVar(&tokenUUID, "token_uuid", "", "token uuid")
	flag.StringVar(&expiration, "expiration", "", "expiration date with format: '2006-01-02'")
	flag.StringVar(&domains, "domains", "", "domains list separated by comma - domain.com,sub.domain.com")
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

	log.Println("Origin data from token: ", instUUID, expitedAt, domainsResultList)
	log.Println("Token hash: ", fmt.Sprintf("%x", sha256.Sum256([]byte(token))))
}
