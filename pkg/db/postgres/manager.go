package postgres

import (
	"context"
	"fmt"
)

// ConnectionManager is a convenience wrapper for managing connections
type ConnectionManager struct {
	pool *ConnectionPool
}

// NewConnectionManager creates a new connection manager
func NewConnectionManager() *ConnectionManager {
	return &ConnectionManager{
		pool: GetConnectionPool(),
	}
}

// InitializeConnections initializes multiple database connections
func (cm *ConnectionManager) InitializeConnections(ctx context.Context, configs map[string]Config) error {
	for kodeDc, cfg := range configs {
		if err := cm.pool.RegisterConfig(kodeDc, cfg); err != nil {
			return fmt.Errorf("error registering config for kode dc '%s': %w", kodeDc, err)
		}

		// Pre-connect to verify configurations
		if _, err := cm.pool.GetConnection(ctx, kodeDc); err != nil {
			return fmt.Errorf("error initializing connection for kode dc '%s': %w", kodeDc, err)
		}
	}
	return nil
}

// GetDB returns a database connection by key
func (cm *ConnectionManager) GetDB(ctx context.Context, kodeDc string) (*Database, error) {
	return cm.pool.GetConnection(ctx, kodeDc)
}

// Close closes a specific connection
func (cm *ConnectionManager) Close(kodeDc string) error {
	return cm.pool.CloseConnection(kodeDc)
}

// CloseAllConnections closes all connections
func (cm *ConnectionManager) CloseAllConnections() error {
	return cm.pool.CloseAll()
}

// PrintPoolStatus prints the status of all connections
func (cm *ConnectionManager) PrintPoolStatus(ctx context.Context) {
	pool := cm.pool
	pool.mu.RLock()
	keys := make([]string, 0, len(pool.connections))
	for kodeDc := range pool.connections {
		keys = append(keys, kodeDc)
	}
	pool.mu.RUnlock()

	fmt.Println("=== Connection Pool Status ===")
	fmt.Printf("Total connections: %d\n", len(keys))

	for _, kodeDc := range keys {
		stats, err := cm.pool.GetConnectionStats(ctx, kodeDc)
		if err == nil {
			fmt.Printf("  - %s: %v\n", kodeDc, stats)
		} else {
			fmt.Printf("  - %s: Error - %v\n", kodeDc, err)
		}
	}
}

// HealthCheck checks the health of all connections
func (cm *ConnectionManager) HealthCheck(ctx context.Context) map[string]error {
	return cm.pool.HealthCheckAll(ctx)
}

// PrintHealthStatus prints health check results
func (cm *ConnectionManager) PrintHealthStatus(ctx context.Context) {
	results := cm.HealthCheck(ctx)
	fmt.Println("=== Health Check Status ===")
	for kodeDc, err := range results {
		if err == nil {
			fmt.Printf("  ✓ %s: Healthy\n", kodeDc)
		} else {
			fmt.Printf("  ✗ %s: Unhealthy - %v\n", kodeDc, err)
		}
	}
}
