# Change Log

## [v0.0.7, v0.0.8] - 23.07.2025
### Added
* Extended token claim data - added "Fields" management
* Changed JWT-manger service-component:
  * Changed `GenerateJWT` - new signature - `GenerateJWT(expiredAt time.Time, values ...Field) (string, error)`
  * Renamed `GetTokenData` to `GetTokenClaimsData` function
    * Changed function signature `GetTokenClaimsData(accessToken string) (TokenClaimValues, error)`
  * Added `ValidateToken` method
  * Added `DecodeToken` method - `DecodeToken(accessToken string) (*jwt.Token, TokenClaimValues, error)`
* Added new receiver-function `ScanByKey` method to JWT-token claim struct
  * Add new receiver-function to `TokenClaimValues` struct - `ScanByKey(key string, target any) error`
### Changed
* Updated README.md - changed example of usage
* Updated example apps:
  * InstallmentCreate - [installment_creator/main.go](./cmd/installment_creator/main.go)
  * JwtCreator - [jwt_creator/main.go](./cmd/jwt_creator/main.go)
* Fixed linter issues
* Moved copyright management to golangci-lint - `goheader` linter

## [v0.0.6] - 27.04.2025
### Fix
* Bump golang-jwt/jwt/v4 version to - v4.5.2
  * Fixed dependency on jwt-go library version which allows excessive memory allocation during header parsing

## [v0.0.5] - 09.03.2025
### Changed
* Added support last version of
  * [lib-tinyerrors](https://github.com/crypto-bundle/bc-wallet-common-lib-tinyerrors)
  * [lib-errors](https://github.com/crypto-bundle/bc-wallet-common-lib-errors)
* Added support of Go 1.23
* Added linter and fixed up all linter issues
* Updated License
  * Copyright - new year 2025
  * MIT -> MIT NON-AI
  * Added License banner to *.go files

## [v0.0.4] - 18.05.2024
### Changed
* Claim builder now internal part of jwt package
* Changd GenerateJWT method of jwt-service:
  * New signature - `GenerateJWT(expiredAt time.Time, values map[string]string) (string, error)`

## [v0.0.2, v0.0.3] - 13.05.2024
### Changed
* Added JWT typical config struct
* Re-worked of custom claim - added Claim Builder service
* Re-worked JWT-service:
  * Removed logger dependency
  * Added GenerateJWT receiver method
  * Added GetTokenData receiver method

## [initial - v0.0.1] - 01.05.2023
### Info
* lib-jwt move to another repository - https://github.com/crypto-bundle/bc-wallet-common-lib-jwt