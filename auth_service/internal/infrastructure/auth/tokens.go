package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type TokenType string

const (
	TokenTypeAccess TokenType = "access"
	TokeTypeRefresh TokenType = "refresh"
)

type TokenManager struct {
	TokenSecret         string
	AccessTokenExpTime  time.Duration
	RefreshTokenExpTime time.Duration
}

func (t TokenManager) MakeAccessToken(userID uuid.UUID) (string, error) {
	return t.makeJWT(
		TokenTypeAccess,
		userID,
		t.RefreshTokenExpTime,
		t.TokenSecret,
	)
}

func (t TokenManager) MakeRefreshToken(userID uuid.UUID) (string, error) {
	return t.makeJWT(
		TokeTypeRefresh,
		userID,
		t.RefreshTokenExpTime,
		t.TokenSecret,
	)
}

func (t TokenManager) ValidateAccessToken(
	tokenString string,
) (uuid.UUID, error) {
	return t.validateJWT(tokenString, TokenTypeAccess)
}

func (t TokenManager) ValidateRefreshToken(
	tokenString string,
) (uuid.UUID, error) {
	return t.validateJWT(tokenString, TokeTypeRefresh)
}

func (t TokenManager) validateJWT(
	tokenString string,
	tokenType TokenType,
) (uuid.UUID, error) {
	claimsStruct := jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(
		tokenString,
		&claimsStruct,
		func(token *jwt.Token) (interface{}, error) { return []byte(t.TokenSecret), nil },
	)
	if err != nil {
		return uuid.Nil, err
	}

	userIDString, err := token.Claims.GetSubject()
	if err != nil {
		return uuid.Nil, err
	}

	issuer, err := token.Claims.GetIssuer()
	if err != nil {
		return uuid.Nil, err
	}
	if issuer != string(tokenType) {
		return uuid.Nil, errors.New("invalid issuer")
	}

	id, err := uuid.Parse(userIDString)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid user ID: %w", err)
	}
	return id, nil
}

func (t TokenManager) makeJWT(
	tokenType TokenType,
	userID uuid.UUID,
	expiresIn time.Duration,
	tokenSecret string,
) (string, error) {
	signingKey := []byte(tokenSecret)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    string(tokenType),
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(expiresIn)),
		Subject:   userID.String(),
	})
	return token.SignedString(signingKey)
}
