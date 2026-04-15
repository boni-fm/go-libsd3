package tzlocal

import (
	"fmt"
	"sync"
	"time"
)

var (
	once     sync.Once
	initErr  error
	location *time.Location
)

// Init detects and globally sets the system timezone.
// It is safe to call multiple times; detection runs only once.
// Call this once at program startup (e.g., in main or an init() function).
func Init() error {
	once.Do(func() {
		tzName, err := getSystemTimezone()
		if err != nil {
			initErr = fmt.Errorf("tzlocal: failed to detect system timezone: %w", err)
			return
		}

		loc, err := time.LoadLocation(tzName)
		if err != nil {
			// Fallback: try loading "Local" directly
			loc, err = time.LoadLocation("Local")
			if err != nil {
				initErr = fmt.Errorf("tzlocal: failed to load timezone %q: %w", tzName, err)
				return
			}
		}

		// Set globally — affects ALL time.Now() calls in the entire program
		time.Local = loc
		location = loc
	})

	return initErr
}

// MustInit is like Init but panics on error.
// Useful in main() where a missing timezone is unrecoverable.
func MustInit() {
	if err := Init(); err != nil {
		panic(err)
	}
}

// Get returns the currently loaded *time.Location.
// Returns nil if Init has not been called yet.
func Get() *time.Location {
	return location
}

// Name returns the timezone name string (e.g., "Asia/Jakarta").
// Returns an empty string if Init has not been called yet.
func Name() string {
	if location == nil {
		return ""
	}
	return location.String()
}

// Now returns time.Now() in the globally registered timezone.
// Equivalent to time.Now().In(tzlocal.Get()) but safe even before Init().
func Now() time.Time {
	if location != nil {
		return time.Now().In(location)
	}
	return time.Now()
}

// Reset allows re-detection (useful in tests or hot-reload scenarios).
func Reset() {
	once = sync.Once{}
	location = nil
	initErr = nil
}
