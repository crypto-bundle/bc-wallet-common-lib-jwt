# bc-wallet-common-lib-jwt

## Description

Library for management JWT entities - small wrapper for **golang-jwt** library - [github.com/golang-jwt/jwt/v4](github.com/golang-jwt/jwt/v4)

Supports:
* Generate JWT
* Decrypt JWT

Also contains binary application for generate JWT token by CLI command

### Usage example
You can see and build generator preview in [cmd/jwt_creator](cmd/jwt_creator/main.go)

```go
import (
    "errors"
    "flag"
    "log"
    "time"
    
    jwttool "gitlab.heronodes.io/bc-platform/bc-wallet-common-lib-jwt/pkg/jwt"
)

const (
    DataInfoLabel = "data_label_uuid" 
	  DataInfoUUID = "930de0c3-fd5f-4fa4-bcaa-ced05a68eade"
	  ExpiredAtLabel = "expired_at"
)


func main() {
    // set expiration time
    mExpTime := time.Now().Add(jwttool.LongLiveTokenDuration)

    jwtSrv := jwttool.NewJWTManger(mockerrors.NewMockErrFormatter(), key)

    token, err := jwtSvc.GenerateJWT(mExpTime,
		    jwttool.String(DataInfoLabel, DataInfoUUID),
        jwttool.Time(ExpiredAtLabel, mExpTime),
		)
    if err != nil {
        log.Fatalf("unable to  make JWT token. Error: %v ", err.Error())
    }

    data, err := jwtSvc.GetTokenClaimsData(token)
    if err != nil {
        log.Fatalf("unable to get data from JWT token. Error: %v ", err.Error())
    }
    
    var expitedAt time.Time
    err = data.ScanByKey(ExpiredAtLabel, &expitedAt)
    if err != nil {
        log.Fatalf("unable to scan data from JWT token. Error: %v ", err.Error())
    }

    var dataInfoUUID string
    err = data.ScanByKey(DataInfoLabel, &dataInfoUUID)
    if err != nil {
        log.Fatalf("unable to scan data from JWT token. Error: %v ", err.Error())
    }
	
    log.Println("Token: ", token)
    log.Println("ExpiredAt: ", expitedAt)
	  log.Println("DataInfoUUID: ", dataInfoUUID)
}
```

## Contributors

* Author and maintainer - [@gudron (Alex V Kotelnikov)](https://github.com/gudron)

## Licence

**bc-wallet-common-lib-jwt** is licensed under the [MIT NON-AI](./LICENSE) License.