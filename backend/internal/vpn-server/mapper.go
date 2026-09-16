package server

import "time"

func ToResponse(s Server) ServerResponse {
	return ServerResponse{
		ID:             s.ID,
		Region:         s.Region,
		Name:           s.Name,
		PublicEndpoint: s.PublicEndpoint,
		WGPublicKey:    s.WGPublicKey,
		GRPCEndpoint:   s.GRPCEndpoint,
		MaxPeers:       s.MaxPeers,
		CurrentPeers:   s.CurrentPeers,
		IsActive:       s.IsActive,
		CreatedAt:      s.CreatedAt.Format(time.RFC3339),
		UpdatedAt:      s.UpdatedAt.Format(time.RFC3339),
	}
}

func ToResponseList(servers []Server) []ServerResponse {
	responses := make([]ServerResponse, len(servers))
	for i, s := range servers {
		responses[i] = ToResponse(s)
	}
	return responses
}
