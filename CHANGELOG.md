# Change Log

## [initial - v0.0.1] - 01.05.2023
### Info
* lib-jwt move to another repository - https://github.com/crypto-bundle/bc-wallet-common-lib-jwt

## [v0.0.2] - 13.05.2024
### Changed
* Added JWT typical config struct
* Re-worked of custom claim - added Claim Builder service
* Re-worked JWT-service:
  * Removed logger dependency
  * Added GenerateJWT receiver method
  * Added GetTokenData receiver method