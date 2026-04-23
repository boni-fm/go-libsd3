package postgres

import "errors"

// ErrConnClose adalah error yang dikembalikan ketika koneksi database sudah ditutup.
var ErrConnClose = errors.New("postgres: koneksi sudah ditutup")

// ErrEmptyConnString adalah error yang dikembalikan ketika connection string kosong.
var ErrEmptyConnString = errors.New("postgres: connection string kosong")

// ErrConfigNotFound adalah error yang dikembalikan ketika konfigurasi tidak ditemukan.
var ErrConfigNotFound = errors.New("postgres: konfigurasi tidak ditemukan")

// ErrConnectionNotFound adalah error yang dikembalikan ketika koneksi tidak ditemukan.
var ErrConnectionNotFound = errors.New("postgres: koneksi tidak ditemukan")

// ErrConfigExists adalah error yang dikembalikan ketika konfigurasi dengan kunci yang sama sudah ada.
var ErrConfigExists = errors.New("postgres: konfigurasi dengan kunci yang sama sudah ada")
