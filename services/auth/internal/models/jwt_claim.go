package models

import "github.com/golang-jwt/jwt/v5"

type JwtClaim struct {
	jwt.RegisteredClaims
	Data map[string]any
}
