package auth

import (
	"backend/internal/config"
	"backend/internal/user"
	apperrors "backend/pkg/errors"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Service struct {
	userService *user.Service
	jwtCfg      config.JWTConfig
}

func NewService(userService *user.Service, jwtCfg config.JWTConfig) *Service {
	return &Service{
		userService: userService,
		jwtCfg:      jwtCfg,
	}
}

type Claims struct {
	UserID uuid.UUID `json:"user_id"`
	Email  string    `json:"email"`
	jwt.RegisteredClaims
}

func (s *Service) Register(req RegisterRequest) (RegisterResponse, error) {
	userResp, err := s.userService.Create(user.CreateUserRequest{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Phone:     req.Phone,
		Password:  req.Password,
	})
	if err != nil {
		return RegisterResponse{}, err
	}

	return RegisterResponse{User: userResp}, nil
}

func (s *Service) Login(req LoginRequest) (LoginResponse, error) {
	u, err := s.userService.GetByEmail(req.Email)
	if err != nil {
		var appErr *apperrors.AppError
		if errors.As(err, &appErr) && errors.Is(appErr.Err, apperrors.ErrNotFound) {
			return LoginResponse{}, apperrors.NewUnauthorized("invalid email or password")
		}
		return LoginResponse{}, err
	}

	if err := user.CheckPassword(u.Password, req.Password); err != nil {
		return LoginResponse{}, apperrors.NewUnauthorized("invalid email or password")
	}

	token, err := s.generateToken(u.ID, u.Email)
	if err != nil {
		return LoginResponse{}, err
	}

	return LoginResponse{
		Token: token,
		User:  user.ToResponse(*u),
	}, nil
}

func (s *Service) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, apperrors.NewUnauthorized("invalid token")
		}
		return []byte(s.jwtCfg.Secret), nil
	})
	if err != nil {
		return nil, apperrors.NewUnauthorized("invalid or expired token")
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, apperrors.NewUnauthorized("invalid token")
	}
	return claims, nil
}

func (s *Service) generateToken(userID uuid.UUID, email string) (string, error) {
	now := time.Now()
	claims := Claims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Duration(s.jwtCfg.ExpirationHours) * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(s.jwtCfg.Secret))
	if err != nil {
		return "", err
	}
	return signed, nil
}
