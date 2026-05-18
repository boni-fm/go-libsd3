package retry

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestDo_SuccessOnFirstAttempt(t *testing.T) {
	cfg := DefaultConfig()
	calls := 0
	err := Do(context.Background(), cfg, func() error {
		calls++
		return nil
	})
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if calls != 1 {
		t.Errorf("expected 1 call, got %d", calls)
	}
}

func TestDo_SuccessOnRetry(t *testing.T) {
	cfg := Config{MaxAttempts: 3, BaseDelay: time.Millisecond, MaxDelay: time.Second, Multiplier: 2}
	calls := 0
	err := Do(context.Background(), cfg, func() error {
		calls++
		if calls < 3 {
			return errors.New("temporary error")
		}
		return nil
	})
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if calls != 3 {
		t.Errorf("expected 3 calls, got %d", calls)
	}
}

func TestDo_AllAttemptsFail(t *testing.T) {
	cfg := Config{MaxAttempts: 3, BaseDelay: time.Millisecond, MaxDelay: time.Second, Multiplier: 2}
	err := Do(context.Background(), cfg, func() error {
		return errors.New("persistent error")
	})
	if err == nil {
		t.Error("expected error when all attempts fail")
	}
}

func TestDo_ContextCancelled(t *testing.T) {
	cfg := Config{MaxAttempts: 5, BaseDelay: time.Second, MaxDelay: 10 * time.Second, Multiplier: 2}
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	err := Do(ctx, cfg, func() error {
		return errors.New("error")
	})
	if err == nil {
		t.Error("expected error when context is cancelled")
	}
}

func TestDoWithResult_Success(t *testing.T) {
	cfg := DefaultConfig()
	result, err := DoWithResult(context.Background(), cfg, func() (int, error) {
		return 42, nil
	})
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if result != 42 {
		t.Errorf("expected 42, got %d", result)
	}
}

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()
	if cfg.MaxAttempts != 3 {
		t.Errorf("expected MaxAttempts=3, got %d", cfg.MaxAttempts)
	}
}
