package swiftx

import (
	"testing"
)

func TestPingService_Pong(t *testing.T) {
	n := 100
	res, err := client.Services.Ping.Pong(ctx, n)
	if err != nil {
		t.Fatalf("client.Services.Ping.Pong() error: %v", err)
	}
	if res != n {
		t.Errorf("expected %d, got %d", n, res)
	}
}
