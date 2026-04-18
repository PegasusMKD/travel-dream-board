package auth

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/PegasusMKD/travel-dream-board/internal/db"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/oauth2"
)

type Service interface {
	ProcessGoogleCallback(ctx context.Context, code string) (*db.User, string, error)
	ValidateToken(tokenString string) (string, error)
}

type authServiceImpl struct {
	oauthConfig *oauth2.Config
	repo        Repository
	jwtSecret   []byte
}

func NewService(repo Repository, oauthConfig *oauth2.Config, jwtSecret string) Service {
	return &authServiceImpl{
		oauthConfig: oauthConfig,
		repo:        repo,
		jwtSecret:   []byte(jwtSecret),
	}
}

func (s *authServiceImpl) ProcessGoogleCallback(ctx context.Context, code string) (*db.User, string, error) {
	// Exchange code for token
	token, err := s.oauthConfig.Exchange(ctx, code)
	if err != nil {
		return nil, "", err
	}

	// Make a request to Google API using token
	client := s.oauthConfig.Client(ctx, token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()

	var gUser GoogleUser
	if err := json.NewDecoder(resp.Body).Decode(&gUser); err != nil {
		return nil, "", err
	}

	// Upsert User
	avatar := gUser.Picture
	var avatarPtr *string
	if avatar != "" {
		avatarPtr = &avatar
	}

	params := db.UpsertUserParams{
		Email:     &gUser.Email,
		Name:      gUser.Name,
		AvatarUrl: avatarPtr,
	}

	user, err := s.repo.UpsertUser(ctx, params)
	if err != nil {
		return nil, "", err
	}

	// Generate JWT
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":   user.Uuid.String(),
		"email": user.Email,
		"exp":   time.Now().Add(time.Hour * 24 * 7).Unix(), // 7 days
	})

	tokenString, err := jwtToken.SignedString(s.jwtSecret)
	if err != nil {
		return nil, "", err
	}

	return user, tokenString, nil
}

func (s *authServiceImpl) ValidateToken(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return s.jwtSecret, nil
	})

	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		sub, ok := claims["sub"].(string)
		if !ok {
			return "", errors.New("sub claim missing or not string")
		}
		return sub, nil
	}

	return "", errors.New("invalid token")
}
