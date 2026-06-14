package security

type PasswordEncoder interface {
	Encode(rawPassword string) (string, error)
	Verify(hashedPassword, rawPassword string) bool
}

type AccessTokenClaims struct {
	UserID string
	Iat    int
}

type RefreshTokenClaims struct {
	UserID string
	Iat    int
}

type JwtProvider interface {
	GenerateAccessToken(userID string) (string, error)
	GenerateRefreshToken(userID string) (string, error)
	ValidateAccessToken(token string) (AccessTokenClaims, error)
	ValidateRefreshToken(token string) (RefreshTokenClaims, error)
}
