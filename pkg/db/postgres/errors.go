package postgres

import "errors"

// koleksi custom error untuk db postgres
var (
	ErrConnClose          = errors.New("Yikes! ~ postgres :: connection is closed")
	ErrEmptyConnString    = errors.New("Yikes! ~ postgres :: empty connection string")
	ErrConfigNotFound     = errors.New("Yikes! ~ postgres :: config not found")
	ErrConnectionNotFound = errors.New("Yikes! ~ postgres :: connection not found")
	ErrConfigExists       = errors.New("Yikes! ~ postgres :: config already exists")
)
