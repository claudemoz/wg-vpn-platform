package server

type CreateServerRequest struct {
	Region         string `json:"region"          binding:"required,min=2,max=10"`
	Name           string `json:"name"            binding:"required,min=1,max=100"`
	PublicEndpoint string `json:"public_endpoint" binding:"required,max=255"`
	WGPublicKey    string `json:"wg_public_key"   binding:"required,max=64"`
	GRPCEndpoint   string `json:"grpc_endpoint"   binding:"required,max=255"`
	MaxPeers       int    `json:"max_peers"       binding:"omitempty,min=1"`
	IsActive       *bool  `json:"is_active"`
}

type UpdateServerRequest struct {
	Region         *string `json:"region"          binding:"omitempty,min=2,max=10"`
	Name           *string `json:"name"            binding:"omitempty,min=1,max=100"`
	PublicEndpoint *string `json:"public_endpoint" binding:"omitempty,max=255"`
	WGPublicKey    *string `json:"wg_public_key"   binding:"omitempty,max=64"`
	GRPCEndpoint   *string `json:"grpc_endpoint"   binding:"omitempty,max=255"`
	MaxPeers       *int    `json:"max_peers"       binding:"omitempty,min=1"`
	CurrentPeers   *int    `json:"current_peers"   binding:"omitempty,min=0"`
	IsActive       *bool   `json:"is_active"`
}

type ServerResponse struct {
	ID             string `json:"id"`
	Region         string `json:"region"`
	Name           string `json:"name"`
	PublicEndpoint string `json:"public_endpoint"`
	WGPublicKey    string `json:"wg_public_key"`
	GRPCEndpoint   string `json:"grpc_endpoint"`
	MaxPeers       int    `json:"max_peers"`
	CurrentPeers   int    `json:"current_peers"`
	IsActive       bool   `json:"is_active"`
	CreatedAt      string `json:"created_at"`
	UpdatedAt      string `json:"updated_at"`
}
