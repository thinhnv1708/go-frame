package security

import (
	"identify/internal/exception"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JwtProviderImp struct {
	accessTokenSecret  []byte
	refreshTokenSecret []byte
	accessTokenTtlSec  int
	refreshTokenTtlSec int
}

type JwtAccessTokenClaims struct {
	UserID string `json:"userId"`
	jwt.RegisteredClaims
}

type JwtRefreshTokenClaims struct {
	UserID string `json:"userId"`
	jwt.RegisteredClaims
}

func NewJwtProviderImp(accessTokenSecret, refreshTokenSecret string, accessTokenTtlSec, refreshTokenTtlSec int) JwtProvider {
	return &JwtProviderImp{
		accessTokenSecret:  []byte(accessTokenSecret),
		refreshTokenSecret: []byte(refreshTokenSecret),
		accessTokenTtlSec:  accessTokenTtlSec,
		refreshTokenTtlSec: refreshTokenTtlSec,
	}
}

func (p *JwtProviderImp) GenerateAccessToken(userID string) (string, error) {
	claims := JwtAccessTokenClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(p.accessTokenTtlSec) * time.Second)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(p.accessTokenSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (p *JwtProviderImp) GenerateRefreshToken(userID string) (string, error) {
	claims := JwtRefreshTokenClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(p.refreshTokenTtlSec) * time.Second)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(p.refreshTokenSecret)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (p *JwtProviderImp) ValidateAccessToken(tokenString string) (AccessTokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JwtAccessTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return p.accessTokenSecret, nil
	})

	if err != nil {
		return AccessTokenClaims{}, exception.INVALID_ACCESS_TOKEN
	}

	claims, ok := token.Claims.(*JwtAccessTokenClaims)
	if !ok || !token.Valid {
		return AccessTokenClaims{}, exception.INVALID_ACCESS_TOKEN
	}

	var iat int64
	if claims.IssuedAt != nil {
		iat = claims.IssuedAt.Unix()
	}

	return AccessTokenClaims{UserID: claims.UserID, Iat: int(iat)}, nil
}

func (p *JwtProviderImp) ValidateRefreshToken(tokenString string) (RefreshTokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JwtRefreshTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return p.refreshTokenSecret, nil
	})

	if err != nil {
		return RefreshTokenClaims{}, exception.INVALID_REFRESH_TOKEN
	}

	claims, ok := token.Claims.(*JwtRefreshTokenClaims)
	if !ok || !token.Valid {
		return RefreshTokenClaims{}, exception.INVALID_REFRESH_TOKEN
	}

	var iat int64
	if claims.IssuedAt != nil {
		iat = claims.IssuedAt.Unix()
	}

	return RefreshTokenClaims{UserID: claims.UserID, Iat: int(iat)}, nil
}
