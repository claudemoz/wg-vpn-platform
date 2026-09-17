package wireguard

import "testing"

func TestPoolAllocateSkipsGatewayAndUsedIPs(t *testing.T) {
	pool, err := NewPool("10.8.0.0/30")
	if err != nil {
		t.Fatalf("NewPool: %v", err)
	}

	ip, err := pool.Allocate([]string{})
	if err != nil {
		t.Fatalf("Allocate: %v", err)
	}

	if ip != "10.8.0.2" {
		t.Fatalf("expected 10.8.0.2, got %s", ip)
	}

	next, err := pool.Allocate([]string{ip})
	if err != nil {
		t.Fatalf("Allocate second: %v", err)
	}

	if next != "10.8.0.3" {
		t.Fatalf("expected 10.8.0.3, got %s", next)
	}

	if _, err := pool.Allocate([]string{ip, next}); err == nil {
		t.Fatal("expected pool exhaustion error")
	}
}

func TestPoolContains(t *testing.T) {
	pool, err := NewPool("10.8.0.0/24")
	if err != nil {
		t.Fatalf("NewPool: %v", err)
	}

	if !pool.Contains("10.8.0.42") {
		t.Fatal("expected ip to be inside subnet")
	}

	if pool.Contains("192.168.0.1") {
		t.Fatal("expected ip to be outside subnet")
	}
}
