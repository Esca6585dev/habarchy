package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	"github.com/Esca6585dev/habarchy/backend/internal/domain"
)

// Claims are the JWT claims of an admin access token.
type Claims struct {
	jwt.RegisteredClaims
	Email string `json:"email"`
}

// UserID returns the subject as a UUID.
func (c Claims) UserID() (uuid.UUID, error) { return uuid.Parse(c.Subject) }

// Tokens is the pair returned by login and refresh.
type Tokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"` // seconds
}

func (s *Service) signAccess(userID uuid.UUID, email string, now time.Time) (string, error) {
	claims := Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    s.cfg.JWTIssuer,
			Subject:   userID.String(),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(s.cfg.AccessTokenTTL)),
			ID:        uuid.NewString(),
		},
		Email: email,
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(s.cfg.JWTSecret))
}

// ParseAccess validates an access token and returns its claims.
func (s *Service) ParseAccess(token string) (*Claims, error) {
	claims := &Claims{}
	_, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (any, error) {
		if t.Method != jwt.SigningMethodHS256 {
			return nil, fmt.Errorf("unexpected signing method %v", t.Header["alg"])
		}
		return []byte(s.cfg.JWTSecret), nil
	},
		jwt.WithIssuer(s.cfg.JWTIssuer),
		jwt.WithExpirationRequired(),
		jwt.WithTimeFunc(s.now),
	)
	if err != nil {
		return nil, domain.ErrUnauthorized.WithMessage("invalid or expired token")
	}
	return claims, nil
}
