package service

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"bench-hub/internal/model"
	"bench-hub/internal/repository"
)

type AuthService struct {
	users      repository.UserRepository
	secret     []byte
	accessTTL  time.Duration
	refreshTTL time.Duration
	issuer     string
}

func NewAuthService(users repository.UserRepository, secret string, accessTTL, refreshTTL time.Duration, issuer string) *AuthService {
	return &AuthService{
		users:      users,
		secret:     []byte(secret),
		accessTTL:  accessTTL,
		refreshTTL: refreshTTL,
		issuer:     issuer,
	}
}

func (s *AuthService) Login(ctx context.Context, username, password string) (string, string, *model.User, error) {
	user, err := s.users.GetByUsername(ctx, username)
	if err != nil {
		if err == repository.ErrNotFound {
			return "", "", nil, ErrInvalidCredentials
		}
		return "", "", nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", "", nil, ErrInvalidCredentials
	}

	access, err := s.newToken(user.ID, "access", s.accessTTL)
	if err != nil {
		return "", "", nil, err
	}
	refresh, err := s.newToken(user.ID, "refresh", s.refreshTTL)
	if err != nil {
		return "", "", nil, err
	}

	return access, refresh, user, nil
}

func (s *AuthService) Refresh(ctx context.Context, refreshToken string) (string, string, error) {
	userID, err := s.validateToken(refreshToken, "refresh")
	if err != nil {
		return "", "", err
	}

	access, err := s.newToken(userID, "access", s.accessTTL)
	if err != nil {
		return "", "", err
	}
	refresh, err := s.newToken(userID, "refresh", s.refreshTTL)
	if err != nil {
		return "", "", err
	}

	return access, refresh, nil
}

func (s *AuthService) ValidateAccess(token string) (string, error) {
	return s.validateToken(token, "access")
}

func (s *AuthService) newToken(userID, tokenType string, ttl time.Duration) (string, error) {
	claims := jwt.RegisteredClaims{
		Subject:   userID,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(ttl)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		Issuer:    s.issuer,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": claims.Subject,
		"exp": claims.ExpiresAt.Unix(),
		"iat": claims.IssuedAt.Unix(),
		"iss": claims.Issuer,
		"typ": tokenType,
	})

	return token.SignedString(s.secret)
}

func (s *AuthService) validateToken(tokenString, tokenType string) (string, error) {
	parsed, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, jwt.ErrSignatureInvalid
		}
		return s.secret, nil
	})
	if err != nil || !parsed.Valid {
		return "", ErrUnauthorized
	}

	claims, ok := parsed.Claims.(jwt.MapClaims)
	if !ok {
		return "", ErrUnauthorized
	}

	typ, _ := claims["typ"].(string)
	if typ != tokenType {
		return "", ErrUnauthorized
	}

	sub, _ := claims["sub"].(string)
	if sub == "" {
		return "", ErrUnauthorized
	}

	return sub, nil
}
