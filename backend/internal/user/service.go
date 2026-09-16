package user

import (
	apperrors "backend/pkg/errors"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetAll() ([]UserResponse, error) {
	users, err := s.repo.FindAll()
	if err != nil {
		return nil, err
	}
	return ToResponseList(users), nil
}

func (s *Service) GetByID(id uuid.UUID) (UserResponse, error) {
	user, err := s.repo.FindByID(id)
	if err != nil {
		return UserResponse{}, err
	}
	return ToResponse(*user), nil
}

func (s *Service) GetByEmail(email string) (*User, error) {
	return s.repo.FindByEmail(email)
}

func (s *Service) Update(id uuid.UUID, req UpdateUserRequest) (UserResponse, error) {
	user, err := s.repo.FindByID(id)
	if err != nil {
		return UserResponse{}, err
	}

	if req.FirstName != nil {
		user.FirstName = *req.FirstName
	}
	if req.LastName != nil {
		user.LastName = *req.LastName
	}
	if req.Email != nil {
		user.Email = *req.Email
	}
	if req.Phone != nil {
		user.Phone = *req.Phone
	}
	if req.Password != nil {
		hashed, err := hashPassword(*req.Password)
		if err != nil {
			return UserResponse{}, err
		}
		user.Password = hashed
	}

	if err := s.repo.Update(user); err != nil {
		return UserResponse{}, err
	}
	return ToResponse(*user), nil
}

func (s *Service) Delete(id uuid.UUID) error {
	return s.repo.Delete(id)
}

func (s *Service) Create(req CreateUserRequest) (UserResponse, error) {
	hashed, err := hashPassword(req.Password)
	if err != nil {
		return UserResponse{}, err
	}

	user := &User{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Phone:     req.Phone,
		Password:  hashed,
	}

	if err := s.repo.Create(user); err != nil {
		return UserResponse{}, err
	}
	return ToResponse(*user), nil
}

func hashPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", apperrors.NewInvalidInput("failed to hash password")
	}
	return string(hashed), nil
}

func CheckPassword(hashed, plain string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(plain))
}
