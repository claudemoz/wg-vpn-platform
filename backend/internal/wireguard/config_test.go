package wireguard

import (
	"strings"
	"testing"
)

func TestRenderClientConfig(t *testing.T) {
	config := RenderClientConfig(ClientConfigInput{
		PrivateKey:      "client-private-key",
		AssignedIP:      "10.8.0.5",
		DNS:             []string{"1.1.1.1"},
		ServerPublicKey: "server-public-key",
		ServerEndpoint:  "vpn.example.com:51820",
		AllowedIPs:      []string{"0.0.0.0/0"},
		Keepalive:       25,
	})

	for _, expected := range []string{
		"PrivateKey = client-private-key",
		"Address = 10.8.0.5/32",
		"DNS = 1.1.1.1",
		"PublicKey = server-public-key",
		"Endpoint = vpn.example.com:51820",
		"AllowedIPs = 0.0.0.0/0",
		"PersistentKeepalive = 25",
	} {
		if !strings.Contains(config, expected) {
			t.Fatalf("expected config to contain %q, got:\n%s", expected, config)
		}
	}
}
