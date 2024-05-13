package jwt

type JWTConfig struct {
	Key            string `envconfig:"JWT_SECRET_KEY" secret:"true"`
	ExpirationTime string `envconfig:"JWT_DEFAULT_TTL" default:"24h"`
}
