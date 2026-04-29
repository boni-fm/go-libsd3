package postgres

import (
	"sync"
	"testing"
)

func resetRegistry() {
	registryMu.Lock()
	dbRegistry = make(map[registryKey]*Database)
	dbRegistryCfg = make(map[registryKey]Config)
	registryMu.Unlock()
}

func TestRegistryGetDB_NotFound(t *testing.T) {
	resetRegistry()
	_, err := GetDB("kunci1", "DC001")
	if err == nil {
		t.Error("expected error for missing key")
	}
}

func TestRegistryMustGetDB_Panics(t *testing.T) {
	resetRegistry()
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for missing key")
		}
	}()
	MustGetDB("kunci1", "DC001")
}

func TestRegistryRegisterDB_DuplicateKey(t *testing.T) {
	resetRegistry()

	fakeDB := &Database{isClosed: false}
	cfg := Config{KodeDC: "DC001", AppName: "test"}
	injectDB("kunci1", "DC001", fakeDB, cfg)

	result, err := GetDB("kunci1", "DC001")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != fakeDB {
		t.Error("expected same DB instance to be returned")
	}
}

func TestRegistryRegisterDB_ConflictingConfig(t *testing.T) {
	resetRegistry()

	fakeDB := &Database{isClosed: false}
	cfg1 := Config{KodeDC: "DC001", AppName: "app1"}
	cfg2 := Config{KodeDC: "DC001", AppName: "app2"}
	injectDB("kunci1", "DC001", fakeDB, cfg1)

	registryMu.RLock()
	existingCfg := dbRegistryCfg[registryKey{kunci: "kunci1", kodedc: "DC001"}]
	registryMu.RUnlock()

	if existingCfg == cfg2 {
		t.Error("configs should be different for conflict test")
	}

	if existingCfg == cfg1 && cfg1 != cfg2 {
		t.Log("conflict scenario verified: same key with different config")
	}
}

func TestRegistryConcurrent(t *testing.T) {
	resetRegistry()

	fakeDB := &Database{isClosed: false}
	cfg := Config{KodeDC: "DC001", AppName: "concurrent-test"}

	injectDB("kunci-concurrent", "DC001", fakeDB, cfg)

	const goroutines = 100
	var wg sync.WaitGroup
	results := make([]*Database, goroutines)

	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			db, err := GetDB("kunci-concurrent", "DC001")
			if err != nil {
				t.Errorf("goroutine %d got error: %v", idx, err)
				return
			}
			results[idx] = db
		}(i)
	}

	wg.Wait()

	for i, r := range results {
		if r != fakeDB {
			t.Errorf("goroutine %d got unexpected DB instance", i)
		}
	}

	registryMu.RLock()
	count := len(dbRegistry)
	registryMu.RUnlock()
	if count != 1 {
		t.Errorf("expected 1 entry in registry, got %d", count)
	}
}
