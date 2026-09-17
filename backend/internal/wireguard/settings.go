package wireguard

import "fmt"

// Settings groups the platform-wide WireGuard parameters used when
// provisioning devices and rendering client configurations.
type Settings struct {
	Pool       *Pool
	DNS        []string
	AllowedIPs []string
	Keepalive  int
}

// NewSettings builds Settings from raw configuration values.
// dns and allowedIPs are comma-separated lists.
func NewSettings(subnet, dns, allowedIPs string, keepalive int) (*Settings, error) {
	pool, err := NewPool(subnet)
	if err != nil {
		return nil, fmt.Errorf("wireguard subnet: %w", err)
	}

	return &Settings{
		Pool:       pool,
		DNS:        SplitList(dns),
		AllowedIPs: SplitList(allowedIPs),
		Keepalive:  keepalive,
	}, nil
}

// ClientConfig renders a client configuration for a device on the given server.
func (s *Settings) ClientConfig(privateKey, assignedIP, serverPublicKey, serverEndpoint string) string {
	return RenderClientConfig(ClientConfigInput{
		PrivateKey:      privateKey,
		AssignedIP:      assignedIP,
		DNS:             s.DNS,
		ServerPublicKey: serverPublicKey,
		ServerEndpoint:  serverEndpoint,
		AllowedIPs:      s.AllowedIPs,
		Keepalive:       s.Keepalive,
	})
}
