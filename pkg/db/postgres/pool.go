package postgres

import (
	"context"
	"fmt"
	"sync"
)

// struct model untuk kumpulan koneksi db nya
type ConnectionPool struct {
	connections map[string]*Database
	configs     map[string]Config
	mu          sync.RWMutex
}

// singleton lifetime untuk kolam renang db nya
var (
	instance *ConnectionPool
	once     sync.Once
)

// inisiasi/mulai instance kolam renang nya
func GetConnectionPool() *ConnectionPool {
	once.Do(func() {
		instance = &ConnectionPool{
			connections: make(map[string]*Database),
			configs:     make(map[string]Config),
		}
	})
	return instance
}

// RegisterConfig registers a database configuration with a unique key
func (cp *ConnectionPool) RegisterConfig(kodeDc string, cfg Config) error {
	cp.mu.Lock()
	defer cp.mu.Unlock()

	if _, exists := cp.configs[kodeDc]; exists {
		return fmt.Errorf("configuration with key '%s' already exists", kodeDc)
	}

	cp.configs[kodeDc] = cfg
	return nil
}

// GetConnection returns a connection from the pool by key
// If the connection doesn't exist, it creates a new one
func (cp *ConnectionPool) GetConnection(ctx context.Context, kodeDc string) (*Database, error) {
	cp.mu.RLock()

	// Check if connection already exists
	if conn, exists := cp.connections[kodeDc]; exists {
		cp.mu.RUnlock()
		return conn, nil
	}

	// Get the configuration for this key
	cfg, configExists := cp.configs[kodeDc]
	cp.mu.RUnlock()

	if !configExists {
		return nil, fmt.Errorf("no configuration found for key '%s'", kodeDc)
	}

	// Create new connection
	db, err := NewDatabase(ctx, &cfg)
	if err != nil {
		return nil, fmt.Errorf("error creating database connection for key '%s': %w", kodeDc, err)
	}

	// Store connection in pool
	cp.mu.Lock()
	cp.connections[kodeDc] = db
	cp.mu.Unlock()

	return db, nil
}

// HasConnection checks if a connection exists in the pool
func (cp *ConnectionPool) HasConnection(kodeDc string) bool {
	cp.mu.RLock()
	defer cp.mu.RUnlock()

	_, exists := cp.connections[kodeDc]
	return exists
}

// CloseConnection closes a specific connection from the pool
func (cp *ConnectionPool) CloseConnection(kodeDc string) error {
	cp.mu.Lock()
	defer cp.mu.Unlock()

	conn, exists := cp.connections[kodeDc]
	if !exists {
		return fmt.Errorf("no connection found for key '%s'", kodeDc)
	}

	err := conn.Close()
	delete(cp.connections, kodeDc)
	return err
}

// CloseAll closes all connections in the pool
func (cp *ConnectionPool) CloseAll() error {
	cp.mu.Lock()
	defer cp.mu.Unlock()

	var errors []error
	for key, conn := range cp.connections {
		if err := conn.Close(); err != nil {
			errors = append(errors, fmt.Errorf("error closing connection '%s': %w", key, err))
		}
	}

	cp.connections = make(map[string]*Database)

	if len(errors) > 0 {
		return fmt.Errorf("multiple errors occurred while closing connections: %v", errors)
	}

	return nil
}

// GetAllKeys returns all registered connection keys
func (cp *ConnectionPool) GetAllKodeDcPoolKey() []string {
	cp.mu.RLock()
	defer cp.mu.RUnlock()

	keys := make([]string, 0, len(cp.connections))
	for key := range cp.connections {
		keys = append(keys, key)
	}
	return keys
}

// GetConnectionStats returns statistics about a connection
func (cp *ConnectionPool) GetConnectionStats(ctx context.Context, kodeDc string) (map[string]interface{}, error) {
	conn, err := cp.GetConnection(ctx, kodeDc)
	if err != nil {
		return nil, err
	}

	stats := conn.GetStats()

	return map[string]interface{}{
		"total_connections":    stats.TotalConns(),
		"idle_connections":     stats.IdleConns(),
		"acquired_connections": stats.AcquiredConns(),
		"construction_slots":   stats.ConstructingConns(),
		"max_size":             stats.MaxConns(),
	}, nil
}

// ReinitializeConnection closes and recreates a connection
func (cp *ConnectionPool) ReinitializeConnection(ctx context.Context, kodeDc string) error {
	cp.mu.Lock()
	defer cp.mu.Unlock()

	// Close existing connection
	if conn, exists := cp.connections[kodeDc]; exists {
		if err := conn.Close(); err != nil {
			return fmt.Errorf("error closing existing connection '%s': %w", kodeDc, err)
		}
		delete(cp.connections, kodeDc)
	}

	// Get configuration
	cfg, configExists := cp.configs[kodeDc]
	if !configExists {
		return fmt.Errorf("no configuration found for key '%s'", kodeDc)
	}

	// Create new connection
	db, err := NewDatabase(ctx, &cfg)
	if err != nil {
		return fmt.Errorf("error creating new database connection for key '%s': %w", kodeDc, err)
	}

	cp.connections[kodeDc] = db
	return nil
}

// UpdateConfig updates the configuration for a key
func (cp *ConnectionPool) UpdateConfig(kodeDc string, cfg Config) {
	cp.mu.Lock()
	defer cp.mu.Unlock()

	cp.configs[kodeDc] = cfg
}

// HealthCheck checks the health of a connection
func (cp *ConnectionPool) HealthCheck(ctx context.Context, kodeDc string) error {
	conn, err := cp.GetConnection(ctx, kodeDc)
	if err != nil {
		return err
	}

	return conn.Ping(ctx)
}

// HealthCheckAll checks health of all connections
func (cp *ConnectionPool) HealthCheckAll(ctx context.Context) map[string]error {
	cp.mu.RLock()
	keys := make([]string, 0, len(cp.connections))
	for kodeDc := range cp.connections {
		keys = append(keys, kodeDc)
	}
	cp.mu.RUnlock()

	results := make(map[string]error)
	for _, kodeDc := range keys {
		results[kodeDc] = cp.HealthCheck(ctx, kodeDc)
	}

	return results
}
