package usecase

import (
	"errors"
	"fmt"
	"time"

	"backend_keping/internal/domain"
	"backend_keping/internal/repository"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrEmailAlreadyExists = errors.New("email already registered")
	ErrInvalidCredentials = errors.New("invalid email or password")
)

type AuthUsecase interface {
	Register(req domain.RegisterRequest) (*domain.AuthResponse, error)
	Login(req domain.LoginRequest) (*domain.AuthResponse, error)
}

type authUsecase struct {
	userRepo          repository.UserRepository
	supabaseJWTSecret string
}

func NewAuthUsecase(userRepo repository.UserRepository, supabaseJWTSecret string) AuthUsecase {
	return &authUsecase{
		userRepo:          userRepo,
		supabaseJWTSecret: supabaseJWTSecret,
	}
}

func (u *authUsecase) Register(req domain.RegisterRequest) (*domain.AuthResponse, error) {
	existingUser, _ := u.userRepo.FindByEmail(req.Email)
	if existingUser != nil {
		return nil, ErrEmailAlreadyExists
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	newUser := &domain.User{
		ID:           uuid.New(),
		Email:        req.Email,
		Name:         req.Name,
		PasswordHash: string(hashedPassword),
		CreatedAt:    time.Now(),
	}

	if err := u.userRepo.Create(newUser); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	token, err := u.generateJWT(newUser.ID, newUser.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return &domain.AuthResponse{
		Token: token,
		User:  *newUser,
	}, nil
}

func (u *authUsecase) Login(req domain.LoginRequest) (*domain.AuthResponse, error) {
	user, err := u.userRepo.FindByEmail(req.Email)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	token, err := u.generateJWT(user.ID, user.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return &domain.AuthResponse{
		Token: token,
		User:  *user,
	}, nil
}

func (u *authUsecase) generateJWT(userID uuid.UUID, email string) (string, error) {
	claims := jwt.MapClaims{
		"sub":   userID.String(),
		"email": email,
		"role":  "authenticated",
		"exp":   time.Now().Add(72 * time.Hour).Unix(),
		"iat":   time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(u.supabaseJWTSecret))
}
