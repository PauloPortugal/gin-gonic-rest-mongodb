package model

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

// JWTOutput Represents the response after a successful attempt.
// swagger:model
type JWTOutput struct {
	// the JWT token
	// required: true
	Token string `json:"token"`
	// Token's time-to-live (TTL)
	// required: true
	Expires time.Time `json:"expires"`
}
