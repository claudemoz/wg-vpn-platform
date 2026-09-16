package device

import "time"

func ToResponse(d Device) DeviceResponse {
	return DeviceResponse{
		ID:            d.ID.String(),
		UserID:        d.UserID.String(),
		ServerID:      d.ServerID.String(),
		Name:          d.Name,
		PublicKey:     d.PublicKey,
		AssignedIP:    d.AssignedIP,
		LastHandshake: formatTimePtr(d.LastHandshake),
		RevokedAt:     formatTimePtr(d.RevokedAt),
		CreatedAt:     d.CreatedAt.Format(time.RFC3339),
		UpdatedAt:     d.UpdatedAt.Format(time.RFC3339),
	}
}

func ToResponseList(devices []Device) []DeviceResponse {
	responses := make([]DeviceResponse, len(devices))
	for i, d := range devices {
		responses[i] = ToResponse(d)
	}
	return responses
}

func formatTimePtr(t *time.Time) *string {
	if t == nil {
		return nil
	}
	s := t.Format(time.RFC3339)
	return &s
}
