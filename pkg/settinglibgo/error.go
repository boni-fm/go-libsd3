package settinglibgo

import (
	"errors"
	"fmt"
)

// ErrEmptyKunci adalah error yang dikembalikan ketika kunci kosong.
var ErrEmptyKunci = errors.New("Yikes! ~ settinglibgo :: key kosong")

// ErrCreateRequestFailed mengembalikan error saat pembuatan HTTP request gagal.
var ErrCreateRequestFailed = func(kunci, err string) error {
	return fmt.Errorf("Yikes! ~ settinglibgo :: failed to create request for kunci=%s, error: %s", kunci, err)
}

// ErrRequestFailed mengembalikan error saat HTTP request ke layanan kunci gagal.
var ErrRequestFailed = func(kunci, err string) error {
	return fmt.Errorf("Yikes! ~ settinglibgo :: failed to hit kunci service for kunci=%s, error: %s", kunci, err)
}

// ErrNonOKResponse mengembalikan error saat response HTTP bukan OK.
var ErrNonOKResponse = func(kunci, status string) error {
	return fmt.Errorf("Yikes! ~ settinglibgo :: failed to get variable for kunci=%s, status: %s", kunci, status)
}

// ErrReadBodyFailed mengembalikan error saat membaca body response gagal.
var ErrReadBodyFailed = func(kunci, err string) error {
	return fmt.Errorf("Yikes! ~ settinglibgo :: failed to read response body for kunci=%s, error: %s", kunci, err)
}
