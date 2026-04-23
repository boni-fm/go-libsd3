package postgres

import "errors"

// ErrConnClose adalah error yang dikembalikan ketika koneksi database sudah ditutup.
var ErrConnClose = errors.New("Yikes! ~ postgres :: connection is closed")

// ErrEmptyConnString adalah error yang dikembalikan ketika connection string kosong.
var ErrEmptyConnString = errors.New("Yikes! ~ postgres :: empty connection string")

// ErrConfigNotFound adalah error yang dikembalikan ketika konfigurasi tidak ditemukan.
var ErrConfigNotFound = errors.New("Yikes! ~ postgres :: config not found")

// ErrConnectionNotFound adalah error yang dikembalikan ketika koneksi tidak ditemukan.
var ErrConnectionNotFound = errors.New("Yikes! ~ postgres :: connection not found")

// ErrConfigExists adalah error yang dikembalikan ketika konfigurasi dengan kunci yang sama sudah ada.
var ErrConfigExists = errors.New("Yikes! ~ postgres :: config already exists")
