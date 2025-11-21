package yaml

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/spf13/viper"
)

var (
	ErrorNullConfigValue = errors.New("config.yaml tidak memiliki value data")
	ErrorNullKeyValue    = errors.New("value dari key tidak ditemukan dalam config.yaml")

	// Cache for config files (thread-safe)
	configCache = struct {
		sync.RWMutex
		data map[string]map[string]interface{}
	}{
		data: make(map[string]map[string]interface{}),
	}
)

// LoadViperInstance creates a new isolated Viper instance per config file
func loadViperInstance(path string) (*viper.Viper, error) {
	v := viper.New() // ✅ Create new instance instead of using global
	v.SetConfigFile(path)
	v.SetConfigType("yaml")

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config from %s: %w", path, err)
	}

	return v, nil
}

// ReadConfigDynamic reads entire config file into map
func ReadConfigDynamic(path string) (map[string]interface{}, error) {
	// Check cache first
	configCache.RLock()
	if cfg, ok := configCache.data[path]; ok {
		configCache.RUnlock()
		return cfg, nil
	}
	configCache.RUnlock()

	// Load fresh instance
	v, err := loadViperInstance(path)
	if err != nil {
		return nil, err
	}

	cfg := v.AllSettings()
	if len(cfg) == 0 {
		return nil, ErrorNullConfigValue
	}

	// Cache the result
	configCache.Lock()
	configCache.data[path] = cfg
	configCache.Unlock()

	return cfg, nil
}

// ReadConfigDynamicWithKey reads a specific key from config file
func ReadConfigDynamicWithKey(path string, key string) (interface{}, error) {
	// Try cache first
	configCache.RLock()
	if cfg, ok := configCache.data[path]; ok {
		configCache.RUnlock()
		if val, ok := cfg[key]; ok {
			return val, nil
		}
		return nil, ErrorNullKeyValue
	}
	configCache.RUnlock()

	// Load fresh instance
	v, err := loadViperInstance(path)
	if err != nil {
		return nil, err
	}

	if !v.IsSet(key) {
		return nil, ErrorNullKeyValue
	}

	val := v.Get(key)

	// Cache the entire config for future calls
	cfg := v.AllSettings()
	configCache.Lock()
	configCache.data[path] = cfg
	configCache.Unlock()

	return val, nil
}

// GetKunciConfigFilepath returns path to config.yaml
// Searches in multiple locations for flexibility
func GetKunciConfigFilepath() (string, error) {
	searchPaths := []string{
		// Priority 1: Same directory as executable
		getExecutableDir(),
		// Priority 2: Current working directory
		".",
		// Priority 3: Environment variable
		os.Getenv("CONFIG_PATH"),
	}

	for _, path := range searchPaths {
		if path == "" {
			continue
		}

		configPath := filepath.Join(path, "config.yaml")
		if _, err := os.Stat(configPath); err == nil {
			abs, _ := filepath.Abs(configPath)
			fmt.Printf("[DEBUG] Found config.yaml at: %s\n", abs)
			return abs, nil
		}
	}

	return "", fmt.Errorf("config.yaml not found in search paths")
}

// getExecutableDir returns directory of running binary
func getExecutableDir() string {
	ex, err := os.Executable()
	if err != nil {
		return ""
	}
	return filepath.Dir(ex)
}

// ClearCache clears the config cache (useful for testing)
func ClearCache() {
	configCache.Lock()
	configCache.data = make(map[string]map[string]interface{})
	configCache.Unlock()
}
