package device

import (
	vpnserver "backend/internal/vpn-server"
	"backend/internal/wireguard"
	apperrors "backend/pkg/errors"

	"github.com/google/uuid"
)

// ServerProvider exposes the VPN servers a device can be attached to.
type ServerProvider interface {
	GetByID(id uuid.UUID) (vpnserver.ServerResponse, error)
}

type Service struct {
	repo    *Repository
	servers ServerProvider
	wg      *wireguard.Settings
}

func NewService(repo *Repository, servers ServerProvider, wg *wireguard.Settings) *Service {
	return &Service{repo: repo, servers: servers, wg: wg}
}

func (s *Service) GetAll() ([]DeviceResponse, error) {
	devices, err := s.repo.FindAll()
	if err != nil {
		return nil, err
	}
	return ToResponseList(devices), nil
}

func (s *Service) GetByID(id uuid.UUID) (DeviceResponse, error) {
	device, err := s.repo.FindByID(id)
	if err != nil {
		return DeviceResponse{}, err
	}
	return ToResponse(*device), nil
}

func (s *Service) Create(userID uuid.UUID, req CreateDeviceRequest) (CreateDeviceResponse, error) {
	server, err := s.resolveServer(req.ServerID)
	if err != nil {
		return CreateDeviceResponse{}, err
	}

	publicKey, privateKey := req.PublicKey, ""
	if publicKey == "" {
		kp, err := wireguard.GenerateKeyPair()
		if err != nil {
			return CreateDeviceResponse{}, err
		}
		publicKey, privateKey = kp.PublicKey, kp.PrivateKey
	} else if err := wireguard.ValidateKey(publicKey); err != nil {
		return CreateDeviceResponse{}, apperrors.NewInvalidInput("invalid public_key")
	}

	assignedIP, err := s.pickIP(server, req.AssignedIP)
	if err != nil {
		return CreateDeviceResponse{}, err
	}

	device := &Device{
		UserID:     userID,
		ServerID:   server.id,
		Name:       req.Name,
		PublicKey:  publicKey,
		AssignedIP: assignedIP,
	}

	if err := s.repo.Create(device); err != nil {
		return CreateDeviceResponse{}, err
	}

	return CreateDeviceResponse{
		DeviceResponse: ToResponse(*device),
		PrivateKey:     privateKey,
		Config:         s.wg.ClientConfig(privateKey, device.AssignedIP, server.WGPublicKey, server.PublicEndpoint),
	}, nil
}

// ClientConfig renders the wg-quick configuration for an existing device.
// The private key is never stored, so it is rendered as a placeholder.
func (s *Service) ClientConfig(id uuid.UUID) (string, error) {
	device, err := s.repo.FindByID(id)
	if err != nil {
		return "", err
	}

	server, err := s.servers.GetByID(device.ServerID)
	if err != nil {
		return "", err
	}

	return s.wg.ClientConfig("", device.AssignedIP, server.WGPublicKey, server.PublicEndpoint), nil
}

func (s *Service) Update(id uuid.UUID, req UpdateDeviceRequest) (DeviceResponse, error) {
	device, err := s.repo.FindByID(id)
	if err != nil {
		return DeviceResponse{}, err
	}

	serverChanged := req.ServerID != nil && *req.ServerID != device.ServerID.String()

	if req.Name != nil {
		device.Name = *req.Name
	}
	if req.PublicKey != nil {
		if err := wireguard.ValidateKey(*req.PublicKey); err != nil {
			return DeviceResponse{}, apperrors.NewInvalidInput("invalid public_key")
		}
		device.PublicKey = *req.PublicKey
	}

	// Any change touching the server or the IP is validated against the
	// target server's subnet; moving to another server re-allocates the IP
	// unless one is explicitly provided.
	if serverChanged || req.AssignedIP != nil {
		targetID := device.ServerID.String()
		if req.ServerID != nil {
			targetID = *req.ServerID
		}
		server, err := s.resolveServer(targetID)
		if err != nil {
			return DeviceResponse{}, err
		}

		requested := ""
		if req.AssignedIP != nil {
			requested = *req.AssignedIP
		} else if !serverChanged {
			requested = device.AssignedIP
		}

		ip, err := s.pickIP(server, requested)
		if err != nil {
			return DeviceResponse{}, err
		}

		device.ServerID = server.id
		device.AssignedIP = ip
	}

	if err := s.repo.Update(device); err != nil {
		return DeviceResponse{}, err
	}
	return ToResponse(*device), nil
}

func (s *Service) Delete(id uuid.UUID) error {
	return s.repo.Delete(id)
}

type resolvedServer struct {
	vpnserver.ServerResponse
	id   uuid.UUID
	pool *wireguard.Pool
}

func (s *Service) resolveServer(raw string) (resolvedServer, error) {
	id, err := uuid.Parse(raw)
	if err != nil {
		return resolvedServer{}, apperrors.NewInvalidInput("invalid server_id")
	}

	server, err := s.servers.GetByID(id)
	if err != nil {
		return resolvedServer{}, err
	}
	if !server.IsActive {
		return resolvedServer{}, apperrors.NewInvalidInput("server is not active")
	}

	pool, err := wireguard.NewPool(server.Subnet)
	if err != nil {
		return resolvedServer{}, apperrors.NewInvalidInput("server has an invalid subnet: " + server.Subnet)
	}

	return resolvedServer{ServerResponse: server, id: id, pool: pool}, nil
}

// pickIP validates the requested IP against the server subnet, or allocates
// a free one from the server pool when requested is empty.
func (s *Service) pickIP(server resolvedServer, requested string) (string, error) {
	if requested != "" {
		if !server.pool.Contains(requested) {
			return "", apperrors.NewInvalidInput("assigned_ip is outside the server subnet " + server.pool.Network())
		}
		return requested, nil
	}

	used, err := s.repo.FindAssignedIPsByServer(server.id)
	if err != nil {
		return "", err
	}

	ip, err := server.pool.Allocate(used)
	if err != nil {
		return "", apperrors.NewConflict("no available ip address in " + server.pool.Network())
	}
	return ip, nil
}
