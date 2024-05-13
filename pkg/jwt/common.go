package jwt

type tokenClaimBuilderService interface {
	Valid() error
}
