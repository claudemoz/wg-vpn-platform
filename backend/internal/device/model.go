package device

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (Device) TableName() string {
	return "device"
}

type Device struct {
	ID            uuid.UUID      `gorm:"column:id;type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	UserID        uuid.UUID      `gorm:"column:user_id;type:uuid;not null;index"                  json:"user_id"`
	ServerID      uuid.UUID      `gorm:"column:server_id;type:uuid;not null;index;uniqueIndex:idx_device_server_assigned_ip" json:"server_id"`
	Name          string         `gorm:"column:name;size:100;not null"                                                        json:"name"`
	PublicKey     string         `gorm:"column:public_key;size:64;not null;uniqueIndex"                                       json:"public_key"`
	AssignedIP    string         `gorm:"column:assigned_ip;size:45;not null;uniqueIndex:idx_device_server_assigned_ip"        json:"assigned_ip"`
	LastHandshake *time.Time     `gorm:"column:last_handshake"                                    json:"last_handshake,omitempty"`
	RevokedAt     *time.Time     `gorm:"column:revoked_at;index"                                  json:"revoked_at,omitempty"`
	CreatedAt     time.Time      `gorm:"column:created_at"                                        json:"created_at"`
	UpdatedAt     time.Time      `gorm:"column:updated_at"                                        json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"column:deleted_at;index"                                  json:"-"`
}

func (d *Device) BeforeCreate(_ *gorm.DB) error {
	if d.ID == uuid.Nil {
		d.ID = uuid.New()
	}
	return nil
}

func (d *Device) IsRevoked() bool {
	return d.RevokedAt != nil
}
