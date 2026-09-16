package device

type CreateDeviceRequest struct {
	ServerID   string `json:"server_id"   binding:"required,uuid"`
	Name       string `json:"name"        binding:"required,min=1,max=100"`
	PublicKey  string `json:"public_key"  binding:"required,max=64"`
	AssignedIP string `json:"assigned_ip" binding:"required,max=45"`
}

type UpdateDeviceRequest struct {
	ServerID   *string `json:"server_id"   binding:"omitempty,uuid"`
	Name       *string `json:"name"        binding:"omitempty,min=1,max=100"`
	PublicKey  *string `json:"public_key"  binding:"omitempty,max=64"`
	AssignedIP *string `json:"assigned_ip" binding:"omitempty,max=45"`
}

type DeviceResponse struct {
	ID            string  `json:"id"`
	UserID        string  `json:"user_id"`
	ServerID      string  `json:"server_id"`
	Name          string  `json:"name"`
	PublicKey     string  `json:"public_key"`
	AssignedIP    string  `json:"assigned_ip"`
	LastHandshake *string `json:"last_handshake,omitempty"`
	RevokedAt     *string `json:"revoked_at,omitempty"`
	CreatedAt     string  `json:"created_at"`
	UpdatedAt     string  `json:"updated_at"`
}
