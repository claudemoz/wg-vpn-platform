package server

import (
	"github.com/google/uuid"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetAll() ([]ServerResponse, error) {
	servers, err := s.repo.FindAll()
	if err != nil {
		return nil, err
	}
	return ToResponseList(servers), nil
}

func (s *Service) GetByID(id uuid.UUID) (ServerResponse, error) {
	server, err := s.repo.FindByID(id)
	if err != nil {
		return ServerResponse{}, err
	}
	return ToResponse(*server), nil
}

func (s *Service) Create(req CreateServerRequest) (ServerResponse, error) {
	maxPeers := req.MaxPeers
	if maxPeers == 0 {
		maxPeers = 100
	}

	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	server := &Server{
		Region:         req.Region,
		Name:           req.Name,
		PublicEndpoint: req.PublicEndpoint,
		WGPublicKey:    req.WGPublicKey,
		GRPCEndpoint:   req.GRPCEndpoint,
		MaxPeers:       maxPeers,
		CurrentPeers:   0,
		IsActive:       isActive,
	}

	if err := s.repo.Create(server); err != nil {
		return ServerResponse{}, err
	}
	return ToResponse(*server), nil
}

func (s *Service) Update(id uuid.UUID, req UpdateServerRequest) (ServerResponse, error) {
	server, err := s.repo.FindByID(id)
	if err != nil {
		return ServerResponse{}, err
	}

	if req.Region != nil {
		server.Region = *req.Region
	}
	if req.Name != nil {
		server.Name = *req.Name
	}
	if req.PublicEndpoint != nil {
		server.PublicEndpoint = *req.PublicEndpoint
	}
	if req.WGPublicKey != nil {
		server.WGPublicKey = *req.WGPublicKey
	}
	if req.GRPCEndpoint != nil {
		server.GRPCEndpoint = *req.GRPCEndpoint
	}
	if req.MaxPeers != nil {
		server.MaxPeers = *req.MaxPeers
	}
	if req.CurrentPeers != nil {
		server.CurrentPeers = *req.CurrentPeers
	}
	if req.IsActive != nil {
		server.IsActive = *req.IsActive
	}

	if err := s.repo.Update(server); err != nil {
		return ServerResponse{}, err
	}
	return ToResponse(*server), nil
}

func (s *Service) Delete(id uuid.UUID) error {
	return s.repo.Delete(id)
}

func (s *Service) Exists(id uuid.UUID) error {
	_, err := s.repo.FindByID(id)
	return err
}
