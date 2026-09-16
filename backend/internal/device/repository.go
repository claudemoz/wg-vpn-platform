package device

import (
	apperrors "backend/pkg/errors"
	"errors"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(device *Device) error {
	if err := r.db.Create(device).Error; err != nil {
		if isUniqueViolation(err) {
			return apperrors.NewConflict("public_key or assigned_ip already in use")
		}
		return err
	}
	return nil
}

func (r *Repository) FindAll() ([]Device, error) {
	var devices []Device
	if err := r.db.Order("created_at ASC").Find(&devices).Error; err != nil {
		return nil, err
	}
	return devices, nil
}

func (r *Repository) FindByID(id uuid.UUID) (*Device, error) {
	var device Device
	if err := r.db.First(&device, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.NewNotFound("device not found")
		}
		return nil, err
	}
	return &device, nil
}

func (r *Repository) Update(device *Device) error {
	if err := r.db.Save(device).Error; err != nil {
		if isUniqueViolation(err) {
			return apperrors.NewConflict("public_key or assigned_ip already in use")
		}
		return err
	}
	return nil
}

func (r *Repository) Delete(id uuid.UUID) error {
	result := r.db.Delete(&Device{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return apperrors.NewNotFound("device not found")
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
