package wireguard

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os"
	"slices"
	"strings"
	"sync"
	"time"

	"golang.zx2c4.com/wireguard/wgctrl"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

var ErrInterfaceNotFound = errors.New("wireguard interface not found")

// Client controls WireGuard interfaces (used by the agent running on VPN servers).
type Client interface {
	// Interface returns the current state of the named interface.
	Interface(ctx context.Context, iface string) (*InterfaceStatus, error)
	// ListPeers returns the peers currently configured on the interface.
	ListPeers(ctx context.Context, iface string) ([]PeerStatus, error)
	// AddPeer adds or updates a single peer.
	AddPeer(ctx context.Context, iface string, peer Peer) error
	// RemovePeer removes the peer identified by its public key.
	RemovePeer(ctx context.Context, iface string, publicKey string) error
	// SyncPeers replaces the whole peer set with the given list.
	SyncPeers(ctx context.Context, iface string, peers []Peer) error
	Close() error
}

// ---------------------------------------------------------------------------
// wgctrl implementation (talks to the kernel module or a userspace daemon)
// ---------------------------------------------------------------------------

type WGCtrlClient struct {
	client *wgctrl.Client
}

// NewWGCtrlClient opens a control client to the local WireGuard implementation.
func NewWGCtrlClient() (*WGCtrlClient, error) {
	c, err := wgctrl.New()
	if err != nil {
		return nil, fmt.Errorf("open wgctrl: %w", err)
	}
	return &WGCtrlClient{client: c}, nil
}

func (c *WGCtrlClient) Interface(_ context.Context, iface string) (*InterfaceStatus, error) {
	dev, err := c.client.Device(iface)
	if err != nil {
		return nil, wrapDeviceErr(iface, err)
	}

	status := &InterfaceStatus{
		Name:       dev.Name,
		PublicKey:  dev.PublicKey.String(),
		ListenPort: dev.ListenPort,
		Peers:      make([]PeerStatus, 0, len(dev.Peers)),
	}
	for _, p := range dev.Peers {
		status.Peers = append(status.Peers, fromWGPeer(p))
	}
	return status, nil
}

func (c *WGCtrlClient) ListPeers(ctx context.Context, iface string) ([]PeerStatus, error) {
	status, err := c.Interface(ctx, iface)
	if err != nil {
		return nil, err
	}
	return status.Peers, nil
}

func (c *WGCtrlClient) AddPeer(_ context.Context, iface string, peer Peer) error {
	cfg, err := toWGPeerConfig(peer)
	if err != nil {
		return err
	}
	cfg.ReplaceAllowedIPs = true

	if err := c.client.ConfigureDevice(iface, wgtypes.Config{Peers: []wgtypes.PeerConfig{cfg}}); err != nil {
		return wrapDeviceErr(iface, err)
	}
	return nil
}

func (c *WGCtrlClient) RemovePeer(_ context.Context, iface string, publicKey string) error {
	key, err := wgtypes.ParseKey(publicKey)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrInvalidKey, err)
	}

	cfg := wgtypes.Config{Peers: []wgtypes.PeerConfig{{PublicKey: key, Remove: true}}}
	if err := c.client.ConfigureDevice(iface, cfg); err != nil {
		return wrapDeviceErr(iface, err)
	}
	return nil
}

func (c *WGCtrlClient) SyncPeers(_ context.Context, iface string, peers []Peer) error {
	cfgs := make([]wgtypes.PeerConfig, 0, len(peers))
	for _, p := range peers {
		cfg, err := toWGPeerConfig(p)
		if err != nil {
			return err
		}
		cfgs = append(cfgs, cfg)
	}

	if err := c.client.ConfigureDevice(iface, wgtypes.Config{ReplacePeers: true, Peers: cfgs}); err != nil {
		return wrapDeviceErr(iface, err)
	}
	return nil
}

func (c *WGCtrlClient) Close() error {
	return c.client.Close()
}

func toWGPeerConfig(p Peer) (wgtypes.PeerConfig, error) {
	key, err := wgtypes.ParseKey(p.PublicKey)
	if err != nil {
		return wgtypes.PeerConfig{}, fmt.Errorf("%w: %v", ErrInvalidKey, err)
	}

	cfg := wgtypes.PeerConfig{PublicKey: key}

	if p.PresharedKey != "" {
		psk, err := wgtypes.ParseKey(p.PresharedKey)
		if err != nil {
			return wgtypes.PeerConfig{}, fmt.Errorf("preshared key: %w: %v", ErrInvalidKey, err)
		}
		cfg.PresharedKey = &psk
	}

	if p.Endpoint != "" {
		addr, err := net.ResolveUDPAddr("udp", p.Endpoint)
		if err != nil {
			return wgtypes.PeerConfig{}, fmt.Errorf("resolve endpoint %q: %w", p.Endpoint, err)
		}
		cfg.Endpoint = addr
	}

	if p.Keepalive > 0 {
		d := time.Duration(p.Keepalive) * time.Second
		cfg.PersistentKeepaliveInterval = &d
	}

	for _, cidr := range p.AllowedIPs {
		_, ipnet, err := net.ParseCIDR(hostIPWithCIDR(cidr))
		if err != nil {
			return wgtypes.PeerConfig{}, fmt.Errorf("parse allowed ip %q: %w", cidr, err)
		}
		cfg.AllowedIPs = append(cfg.AllowedIPs, *ipnet)
	}

	return cfg, nil
}

func fromWGPeer(p wgtypes.Peer) PeerStatus {
	status := PeerStatus{
		Peer: Peer{
			PublicKey: p.PublicKey.String(),
			Keepalive: int(p.PersistentKeepaliveInterval / time.Second),
		},
		LastHandshake: p.LastHandshakeTime,
		ReceiveBytes:  p.ReceiveBytes,
		TransmitBytes: p.TransmitBytes,
	}

	if p.PresharedKey != (wgtypes.Key{}) {
		status.PresharedKey = p.PresharedKey.String()
	}
	if p.Endpoint != nil {
		status.Endpoint = p.Endpoint.String()
	}
	for _, ipnet := range p.AllowedIPs {
		status.AllowedIPs = append(status.AllowedIPs, ipnet.String())
	}

	return status
}

func wrapDeviceErr(iface string, err error) error {
	if errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("%w: %s", ErrInterfaceNotFound, iface)
	}
	return fmt.Errorf("wireguard %s: %w", iface, err)
}

// ---------------------------------------------------------------------------
// In-memory implementation (development & tests, no kernel access required)
// ---------------------------------------------------------------------------

type MemoryClient struct {
	mu    sync.RWMutex
	ifces map[string]map[string]PeerStatus // iface -> publicKey -> peer
}

// NewMemoryClient returns a Client that keeps peer state in memory.
func NewMemoryClient(interfaces ...string) *MemoryClient {
	c := &MemoryClient{ifces: make(map[string]map[string]PeerStatus)}
	for _, name := range interfaces {
		c.ifces[name] = make(map[string]PeerStatus)
	}
	return c
}

func (c *MemoryClient) Interface(_ context.Context, iface string) (*InterfaceStatus, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	peers, ok := c.ifces[iface]
	if !ok {
		return nil, fmt.Errorf("%w: %s", ErrInterfaceNotFound, iface)
	}

	status := &InterfaceStatus{Name: iface, Peers: make([]PeerStatus, 0, len(peers))}
	for _, p := range peers {
		status.Peers = append(status.Peers, p)
	}
	slices.SortFunc(status.Peers, func(a, b PeerStatus) int {
		return strings.Compare(a.PublicKey, b.PublicKey)
	})
	return status, nil
}

func (c *MemoryClient) ListPeers(ctx context.Context, iface string) ([]PeerStatus, error) {
	status, err := c.Interface(ctx, iface)
	if err != nil {
		return nil, err
	}
	return status.Peers, nil
}

func (c *MemoryClient) AddPeer(_ context.Context, iface string, peer Peer) error {
	if err := ValidateKey(peer.PublicKey); err != nil {
		return err
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	peers, ok := c.ifces[iface]
	if !ok {
		return fmt.Errorf("%w: %s", ErrInterfaceNotFound, iface)
	}
	peers[peer.PublicKey] = PeerStatus{Peer: peer}
	return nil
}

func (c *MemoryClient) RemovePeer(_ context.Context, iface string, publicKey string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	peers, ok := c.ifces[iface]
	if !ok {
		return fmt.Errorf("%w: %s", ErrInterfaceNotFound, iface)
	}
	delete(peers, publicKey)
	return nil
}

func (c *MemoryClient) SyncPeers(_ context.Context, iface string, peers []Peer) error {
	for _, p := range peers {
		if err := ValidateKey(p.PublicKey); err != nil {
			return err
		}
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	if _, ok := c.ifces[iface]; !ok {
		return fmt.Errorf("%w: %s", ErrInterfaceNotFound, iface)
	}

	next := make(map[string]PeerStatus, len(peers))
	for _, p := range peers {
		next[p.PublicKey] = PeerStatus{Peer: p}
	}
	c.ifces[iface] = next
	return nil
}

func (c *MemoryClient) Close() error { return nil }

// PublicKeys is a test helper returning the sorted public keys of an interface.
func (c *MemoryClient) PublicKeys(iface string) []string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	keys := make([]string, 0, len(c.ifces[iface]))
	for k := range c.ifces[iface] {
		keys = append(keys, k)
	}
	slices.Sort(keys)
	return keys
}
