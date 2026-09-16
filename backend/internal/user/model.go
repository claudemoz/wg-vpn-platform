package user

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (User) TableName() string {
	return "user"
}

type User struct {
	ID        uuid.UUID `gorm:"column:id;type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	FirstName string    `gorm:"column:first_name;size:80;not null"                       json:"first_name"`
	LastName  string    `gorm:"column:last_name;size:80;not null"                        json:"last_name"`
	Email     string    `gorm:"column:email;size:160;uniqueIndex;not null"               json:"email"`
	Phone     string    `gorm:"column:phone;size:20;uniqueIndex;default:null"            json:"phone,omitempty"`
	Password  string    `gorm:"column:password;size:255;not null"                        json:"-"`
	CreatedAt time.Time `gorm:"column:created_at"                                        json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"                                        json:"updated_at"`
}

func (u *User) BeforeCreate(_ *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return nil
}
