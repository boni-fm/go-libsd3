package fileutil

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/jlaffaye/ftp"
)

// FTPConfig menyimpan parameter koneksi untuk server FTP.
type FTPConfig struct {
	Host     string
	Port     int           // default 21 jika kosong
	Username string
	Password string
	Timeout  time.Duration // default 30 detik jika kosong
}

// FTPClient membungkus koneksi FTP yang aktif.
// Tidak aman untuk penggunaan bersamaan (concurrent) — buat satu per goroutine atau gunakan penguncian eksternal.
type FTPClient struct {
	conn *ftp.ServerConn
	cfg  FTPConfig
}

// NewFTPClient menghubungi server FTP yang dijelaskan oleh cfg, melakukan login, dan mengembalikan FTPClient.
// Pemanggil harus memanggil Close() setelah selesai.
// Port default 21, Timeout default 30 detik.
func NewFTPClient(cfg FTPConfig) (*FTPClient, error) {
	if cfg.Port == 0 {
		cfg.Port = 21
	}
	if cfg.Timeout == 0 {
		cfg.Timeout = 30 * time.Second
	}

	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	conn, err := ftp.Dial(addr, ftp.DialWithTimeout(cfg.Timeout))
	if err != nil {
		return nil, fmt.Errorf("fileutil/ftp: dial %q: %w", addr, err)
	}

	if err := conn.Login(cfg.Username, cfg.Password); err != nil {
		_ = conn.Quit()
		return nil, fmt.Errorf("fileutil/ftp: login: %w", err)
	}

	return &FTPClient{conn: conn, cfg: cfg}, nil
}

// Close mengakhiri sesi FTP. Selalu panggil defer client.Close().
func (c *FTPClient) Close() error {
	if err := c.conn.Quit(); err != nil {
		return fmt.Errorf("fileutil/ftp: quit: %w", err)
	}
	return nil
}

// Upload mengirimkan localPath ke remotePath pada server FTP secara streaming tanpa memuat file ke memori.
func (c *FTPClient) Upload(localPath, remotePath string) error {
	f, err := os.Open(localPath)
	if err != nil {
		return fmt.Errorf("fileutil/ftp: open local %q: %w", localPath, err)
	}
	defer f.Close()

	if err := c.conn.Stor(remotePath, f); err != nil {
		return fmt.Errorf("fileutil/ftp: stor %q: %w", remotePath, err)
	}
	return nil
}

// Download mengunduh remotePath dari server FTP ke localPath secara streaming tanpa memuat file ke memori.
// Direktori induk localPath dibuat jika diperlukan.
func (c *FTPClient) Download(remotePath, localPath string) error {
	if err := os.MkdirAll(aferoDir(localPath), 0755); err != nil {
		return fmt.Errorf("fileutil/ftp: mkdir parent of %q: %w", localPath, err)
	}

	resp, err := c.conn.Retr(remotePath)
	if err != nil {
		return fmt.Errorf("fileutil/ftp: retr %q: %w", remotePath, err)
	}
	defer resp.Close()

	f, err := os.Create(localPath)
	if err != nil {
		return fmt.Errorf("fileutil/ftp: create local %q: %w", localPath, err)
	}
	defer f.Close()

	if _, err := io.Copy(f, resp); err != nil {
		return fmt.Errorf("fileutil/ftp: download %q: %w", remotePath, err)
	}
	return nil
}

// Delete menghapus file di remotePath dari server.
func (c *FTPClient) Delete(remotePath string) error {
	if err := c.conn.Delete(remotePath); err != nil {
		return fmt.Errorf("fileutil/ftp: delete %q: %w", remotePath, err)
	}
	return nil
}

// List mengembalikan nama semua entri dalam direktori remote di remotePath.
func (c *FTPClient) List(remotePath string) ([]string, error) {
	entries, err := c.conn.List(remotePath)
	if err != nil {
		return nil, fmt.Errorf("fileutil/ftp: list %q: %w", remotePath, err)
	}
	names := make([]string, 0, len(entries))
	for _, e := range entries {
		names = append(names, e.Name)
	}
	return names, nil
}

// MakeDir membuat direktori pada server FTP.
func (c *FTPClient) MakeDir(remotePath string) error {
	if err := c.conn.MakeDir(remotePath); err != nil {
		return fmt.Errorf("fileutil/ftp: mkdir %q: %w", remotePath, err)
	}
	return nil
}
