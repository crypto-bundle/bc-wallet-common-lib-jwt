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
)


func main() {

    // set expiration time
    mExpTime := time.Now().Add(jwttool.LongLiveTokenDuration) 
    
	// create claim builder
	claimBuilder := jwttool.NewTokenClaimBuilder(mExpTime)
	
	// Add data to claim - Label, and target data
    err := claimBuilder.AddData(DataInfoLabel, DataInfoUUID)
    if err != nil {
        log.Fatalf("%s", err)
    }

    jwtSrv := jwttool.NewJWTService(key)
    token, err := jwtSrv.GenerateJWT(claimBuilder)
    if err != nil {
        log.Fatalf("unable to  make JWT token. Error: %v ", err.Error())
    }
	
    log.Println("Token: ", token)
	
}
```

## Contributors

* Author and maintainer - [@gudron (Alex V Kotelnikov)](https://github.com/gudron)

## Licence

**bc-wallet-common-lib-jwt** is licensed under the [MIT NON-AI](./LICENSE) License.