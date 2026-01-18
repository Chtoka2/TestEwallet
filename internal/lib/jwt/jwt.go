package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var (
	ErrExpiredToken = errors.New("Token was expired")
	ErrInvalidToken = errors.New("Invalid token")
)

type Claims struct{
	UserID uuid.UUID
	jwt.RegisteredClaims
}

type Service struct{
	secret []byte
}

type JWTGeneratorInterface interface{
	Generate(userID uuid.UUID, ttl time.Duration) (string, error)
}

func NewJWTService(secret string) *Service{
	return &Service{secret: []byte(secret)}
}

func (s *Service) Generate(userID uuid.UUID, ttl time.Duration) (string, error){
	claims := Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(ttl)),
			IssuedAt: jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.secret)
}
func (s *Service) Validate(tokenStr string) (uuid.UUID, error){
	token, err := jwt.ParseWithClaims(
		tokenStr,
		&Claims{},
		func(t *jwt.Token) (interface{}, error) {
			return s.secret, nil
		},
	)
	if err != nil{
		if errors.Is(err, jwt.ErrTokenExpired) {
			return uuid.Nil, ErrExpiredToken
		}
		return uuid.Nil, ErrInvalidToken
	}
	if claims, ok := token.Claims.(*Claims); ok && token.Valid{
		return claims.UserID, nil
	}
	return uuid.Nil, ErrInvalidToken
}