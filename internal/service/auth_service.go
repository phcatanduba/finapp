package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"finapp/internal/model"
	"finapp/internal/repository"
)

var (
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrEmailTaken         = errors.New("email already in use")
	ErrInvalidToken       = errors.New("invalid or expired token")
	ErrNotFound           = errors.New("not found")
	ErrForbidden          = errors.New("forbidden")
)

type authService struct {
	users            repository.UserRepository
	jwtSecret        []byte
	accessTTLMinutes int
	refreshTTLDays   int
}

func NewAuthService(
	users repository.UserRepository,
	jwtSecret string,
	accessTTLMinutes int,
	refreshTTLDays int,
) AuthService {
	return &authService{
		users:            users,
		jwtSecret:        []byte(jwtSecret),
		accessTTLMinutes: accessTTLMinutes,
		refreshTTLDays:   refreshTTLDays,
	}
}

func (s *authService) Register(ctx context.Context, req model.RegisterRequest) (*model.AuthResponse, error) {
	if req.Email == "" || req.Password == "" {
		return nil, fmt.Errorf("email and password are required")
	}

	existing, err := s.users.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, fmt.Errorf("check email: %w", err)
	}
	if existing != nil {
		return nil, ErrEmailTaken
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), 12)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}

	user := &model.User{
		ID:           uuid.New(),
		Email:        req.Email,
		PasswordHash: string(hash),
		Name:         req.Name,
	}
	if err := s.users.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}

	return s.buildAuthResponse(ctx, user)
}

func (s *authService) Login(ctx context.Context, req model.LoginRequest) (*model.AuthResponse, error) {
	user, err := s.users.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, fmt.Errorf("find user: %w", err)
	}
	if user == nil {
		return nil, ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	return s.buildAuthResponse(ctx, user)
}

func (s *authService) RefreshToken(ctx context.Context, refreshToken string) (*model.AuthResponse, error) {
	claims, err := s.ValidateToken(refreshToken)
	if err != nil {
		return nil, ErrInvalidToken
	}
	if claims.Type != "refresh" {
		return nil, ErrInvalidToken
	}

	user, err := s.users.FindByID(ctx, claims.UserID)
	if err != nil || user == nil {
		return nil, ErrInvalidToken
	}

	return s.buildAuthResponse(ctx, user)
}

func (s *authService) ValidateToken(tokenString string) (*model.Claims, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return s.jwtSecret, nil
	}, jwt.WithValidMethods([]string{"HS256"}))

	if err != nil || !token.Valid {
		return nil, ErrInvalidToken
	}

	mapClaims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, ErrInvalidToken
	}

	userIDStr, _ := mapClaims["user_id"].(string)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, ErrInvalidToken
	}

	return &model.Claims{
		UserID: userID,
		Email:  mapClaims["email"].(string),
		Type:   mapClaims["type"].(string),
	}, nil
}

func (s *authService) buildAuthResponse(_ context.Context, user *model.User) (*model.AuthResponse, error) {
	accessToken, err := s.generateToken(user, "access", time.Duration(s.accessTTLMinutes)*time.Minute)
	if err != nil {
		return nil, err
	}
	refreshToken, err := s.generateToken(user, "refresh", time.Duration(s.refreshTTLDays)*24*time.Hour)
	if err != nil {
		return nil, err
	}
	return &model.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         *user,
	}, nil
}

func (s *authService) generateToken(user *model.User, tokenType string, ttl time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"user_id": user.ID.String(),
		"email":   user.Email,
		"type":    tokenType,
		"exp":     time.Now().Add(ttl).Unix(),
		"iat":     time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtSecret)
}
