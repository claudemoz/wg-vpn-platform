package wireguard

import (
	"context"
	"testing"
)

func TestGenerateKeyPairIsValidAndDerivable(t *testing.T) {
	kp, err := GenerateKeyPair()
	if err != nil {
		t.Fatalf("GenerateKeyPair: %v", err)
	}

	if err := ValidateKey(kp.PrivateKey); err != nil {
		t.Fatalf("private key invalid: %v", err)
	}
	if err := ValidateKey(kp.PublicKey); err != nil {
		t.Fatalf("public key invalid: %v", err)
	}

	derived, err := PublicKeyFromPrivate(kp.PrivateKey)
	if err != nil {
		t.Fatalf("PublicKeyFromPrivate: %v", err)
	}
	if derived != kp.PublicKey {
		t.Fatalf("derived public key %q != generated %q", derived, kp.PublicKey)
	}
}

func TestValidateKeyRejectsGarbage(t *testing.T) {
	for _, bad := range []string{"", "not-base64!", "c2hvcnQ="} {
		if err := ValidateKey(bad); err == nil {
			t.Fatalf("expected %q to be rejected", bad)
		}
	}
}

func TestMemoryClientPeerLifecycle(t *testing.T) {
	ctx := context.Background()
	client := NewMemoryClient("wg0")

	kp, _ := GenerateKeyPair()
	peer := NewDevicePeer(kp.PublicKey, "10.8.0.5")

	if err := client.AddPeer(ctx, "wg0", peer); err != nil {
		t.Fatalf("AddPeer: %v", err)
	}

	peers, err := client.ListPeers(ctx, "wg0")
	if err != nil {
		t.Fatalf("ListPeers: %v", err)
	}
	if len(peers) != 1 || peers[0].AllowedIPs[0] != "10.8.0.5/32" {
		t.Fatalf("unexpected peers: %+v", peers)
	}

	if err := client.RemovePeer(ctx, "wg0", kp.PublicKey); err != nil {
		t.Fatalf("RemovePeer: %v", err)
	}
	if keys := client.PublicKeys("wg0"); len(keys) != 0 {
		t.Fatalf("expected no peers, got %v", keys)
	}

	if _, err := client.ListPeers(ctx, "wg1"); err == nil {
		t.Fatal("expected unknown interface error")
	}
}
