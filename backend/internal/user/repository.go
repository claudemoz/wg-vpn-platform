package user

import (
	apperrors "backend/pkg/errors"
	"errors"
	"strings"

	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(user *User) error {
	if err := r.db.Create(user).Error; err != nil {
		if isUniqueViolation(err) {
			return apperrors.NewConflict("email or phone already in use")
		}
		return err
	}
	return nil
}

func (r *Repository) FindAll() ([]User, error) {
	var users []User
	if err := r.db.Order("id ASC").Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (r *Repository) FindByID(id uint) (*User, error) {
	var user User
	if err := r.db.First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.NewNotFound("user not found")
		}
		return nil, err
	}
	return &user, nil
}

func (r *Repository) FindByEmail(email string) (*User, error) {
	var user User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.NewNotFound("user not found")
		}
		return nil, err
	}
	return &user, nil
}

func (r *Repository) Update(user *User) error {
	if err := r.db.Save(user).Error; err != nil {
		if isUniqueViolation(err) {
			return apperrors.NewConflict("email or phone already in use")
		}
		return err
	}
	return nil
}

func (r *Repository) Delete(id uint) error {
	result := r.db.Delete(&User{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return apperrors.NewNotFound("user not found")
	}
	return nil
}

func isUniqueViolation(err error) bool {
	if err == nil {
		return false
	}
	msg := err.Error()
	return errors.Is(err, gorm.ErrDuplicatedKey) ||
		strings.Contains(msg, "duplicate key") ||
		strings.Contains(msg, "UNIQUE constraint")
}
