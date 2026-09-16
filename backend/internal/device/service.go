package device

import (
	apperrors "backend/pkg/errors"

	"github.com/google/uuid"
)

type ServerChecker interface {
	Exists(id uuid.UUID) error
}

type Service struct {
	repo          *Repository
	serverChecker ServerChecker
}

func NewService(repo *Repository, serverChecker ServerChecker) *Service {
	return &Service{repo: repo, serverChecker: serverChecker}
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

func (s *Service) Create(userID uuid.UUID, req CreateDeviceRequest) (DeviceResponse, error) {
	serverID, err := uuid.Parse(req.ServerID)
	if err != nil {
		return DeviceResponse{}, apperrors.NewInvalidInput("invalid server_id")
	}

	if err := s.serverChecker.Exists(serverID); err != nil {
		return DeviceResponse{}, err
	}

	device := &Device{
		UserID:     userID,
		ServerID:   serverID,
		Name:       req.Name,
		PublicKey:  req.PublicKey,
		AssignedIP: req.AssignedIP,
	}

	if err := s.repo.Create(device); err != nil {
		return DeviceResponse{}, err
	}
	return ToResponse(*device), nil
}

func (s *Service) Update(id uuid.UUID, req UpdateDeviceRequest) (DeviceResponse, error) {
	device, err := s.repo.FindByID(id)
	if err != nil {
		return DeviceResponse{}, err
	}

	if req.ServerID != nil {
		serverID, err := uuid.Parse(*req.ServerID)
		if err != nil {
			return DeviceResponse{}, apperrors.NewInvalidInput("invalid server_id")
		}
		if err := s.serverChecker.Exists(serverID); err != nil {
			return DeviceResponse{}, err
		}
		device.ServerID = serverID
	}
	if req.Name != nil {
		device.Name = *req.Name
	}
	if req.PublicKey != nil {
		device.PublicKey = *req.PublicKey
	}
	if req.AssignedIP != nil {
		device.AssignedIP = *req.AssignedIP
	}

	if err := s.repo.Update(device); err != nil {
		return DeviceResponse{}, err
	}
	return ToResponse(*device), nil
}

func (s *Service) Delete(id uuid.UUID) error {
	return s.repo.Delete(id)
}
