package test

import (
	"testing"

	"github.com/boni-fm/go-libsd3/helper/kunci"
)

func TestKunciClient_GetVariable_WithPathKunci(t *testing.T) {
	client := kunci.NewKunciClient("kuncig009sim")

	// Optionally set KUNCI_IP_DOMAIN env if needed
	// os.Setenv("KUNCI_IP_DOMAIN", "localhost")

	val, err := client.GetVariable("UserPostgres", "kuncig009sim")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	t.Logf("Value: %s", val)
}
