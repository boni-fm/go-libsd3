//go:build !windows

package tzlocal

import (
	"errors"
	"os"
	"os/exec"
	"strings"
)

// getSystemTimezone detects the system timezone on Linux/macOS.
// Resolution order:
//  1. TZ environment variable
//  2. /etc/timezone (Debian/Ubuntu)
//  3. timedatectl (systemd-based distros)
//  4. /etc/localtime symlink resolution (RHEL/macOS)
//  5. Fallback to "Local"
func getSystemTimezone() (string, error) {
	// 1. TZ environment variable
	if tz := os.Getenv("TZ"); tz != "" {
		return tz, nil
	}

	// 2. /etc/timezone (Debian, Ubuntu)
	if data, err := os.ReadFile("/etc/timezone"); err == nil {
		tz := strings.TrimSpace(string(data))
		if tz != "" {
			return tz, nil
		}
	}

	// 3. timedatectl (systemd distros: CentOS 7+, Fedora, Arch, etc.)
	if out, err := exec.Command("timedatectl", "show", "--property=Timezone", "--value").Output(); err == nil {
		tz := strings.TrimSpace(string(out))
		if tz != "" {
			return tz, nil
		}
	}

	// 4. /etc/localtime symlink (RHEL, macOS, Alpine)
	if link, err := os.Readlink("/etc/localtime"); err == nil {
		// e.g. /usr/share/zoneinfo/Asia/Jakarta  ->  Asia/Jakarta
		const zoneInfoPath = "zoneinfo/"
		if idx := strings.Index(link, zoneInfoPath); idx != -1 {
			tz := link[idx+len(zoneInfoPath):]
			if tz != "" {
				return tz, nil
			}
		}
	}

	// 5. Fallback
	if _, err := os.Stat("/etc/localtime"); err == nil {
		return "Local", nil
	}

	return "", errors.New("unable to determine system timezone")
}
