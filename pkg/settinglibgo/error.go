package settinglibgo

import (
	"errors"
	"fmt"
)

var (
	ErrEmptyKunci          = errors.New("Yikes! ~ settinglibgo :: key kosong")
	ErrCreateRequestFailed = func(kunci, err string) error {
		return fmt.Errorf("Yikes! ~ settinglibgo :: failed to create request for kunci=%s, error: %s", kunci, err)
	}

	ErrRequestFailed = func(kunci, err string) error {
		return fmt.Errorf("Yikes! ~ settinglibgo :: failed to hit kunci service for kunci=%s, error: %s", kunci, err)
	}

	ErrNonOKResponse = func(kunci, status string) error {
		return fmt.Errorf("Yikes! ~ settinglibgo :: failed to get variable for kunci=%s, status: %s", kunci, status)
	}

	ErrReadBodyFailed = func(kunci, err string) error {
		return fmt.Errorf("Yikes! ~ settinglibgo :: failed to read response body for kunci=%s, error: %s", kunci, err)
	}
)
