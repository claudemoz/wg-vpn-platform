package server

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (Server) TableName() string {
	return "server"
}

type Server struct {
	ID             uuid.UUID      `gorm:"column:id;type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Region         string         `gorm:"column:region;size:10;not null;index"                    json:"region"`
	Name           string         `gorm:"column:name;size:100;not null"                           json:"name"`
	PublicEndpoint string         `gorm:"column:public_endpoint;size:255;not null"                json:"public_endpoint"`
	WGPublicKey    string         `gorm:"column:wg_public_key;size:64;not null;uniqueIndex"       json:"wg_public_key"`
	GRPCEndpoint   string         `gorm:"column:grpc_endpoint;size:255;not null"                  json:"grpc_endpoint"`
	Subnet         string         `gorm:"column:subnet;size:43;not null"                          json:"subnet"`
	MaxPeers       int            `gorm:"column:max_peers;not null;default:100"                   json:"max_peers"`
	CurrentPeers   int            `gorm:"column:current_peers;not null;default:0"                 json:"current_peers"`
	IsActive       bool           `gorm:"column:is_active;not null;default:true;index"            json:"is_active"`
	CreatedAt      time.Time      `gorm:"column:created_at"                                       json:"created_at"`
	UpdatedAt      time.Time      `gorm:"column:updated_at"                                       json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"column:deleted_at;index"                                 json:"-"`
}

func (s *Server) BeforeCreate(_ *gorm.DB) error {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	return nil
}
