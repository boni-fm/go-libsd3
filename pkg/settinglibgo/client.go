package settinglibgo

/*

	Http Client untuk manggil kunci
	Note:
	- manggil POST request /GetVariabel dari apiservice kunci
	- IP apikunci didapatkan dari env KUNCI_IP_DOMAIN -> docker-hub-nginx-1
	  > jika env kosong, maka default localhost
	- body nya make prefix "mujiyono"
	- ada 3 kali proses retry jika panggilan pertama gagal, dengan linear backoff

*/

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/boni-fm/go-libsd3/pkg/config/constant"
)

var (
	PREFIX_KUNCI = "kunci"
	PREFIX       = "mujiyono"
	BASEURL      = "localhost"
	MAXRETRY     = 3
)

// Struct untuk http client nya
type SettingLibClient struct {
	httpClient *http.Client
	key        string // ini key constr nya, cth. "IPPostgres"
}

// Struct untuk body request kunci
type Params struct {
	Key string `json:"key"`
}

func NewSettingLibClient(kunci string) *SettingLibClient {
	trimmed := strings.TrimSpace(kunci)
	key := trimmed
	if !strings.Contains(strings.ToLower(trimmed), PREFIX_KUNCI) {
		key = PREFIX_KUNCI + trimmed
	}

	return &SettingLibClient{
		// No client-level timeout; per-request context timeout is used instead,
		// so each retry attempt gets its own independent deadline.
		httpClient: &http.Client{},
		key:        key,
	}
}

func (kc *SettingLibClient) GetVariable(key string) (string, error) {
	// Read env on every call so runtime changes are picked up,
	// but never write to the global to avoid race conditions.
	baseURL := BASEURL
	if env := os.Getenv("KUNCI_IP_DOMAIN"); env != "" {
		baseURL = env
	}

	url := "http://" + baseURL
	if kc.key != "" {
		url += "/" + kc.key
	}
	url += "/GetVariabel"

	bodyByte, err := json.Marshal(Params{Key: PREFIX + key})
	if err != nil {
		return "", fmt.Errorf("failed to marshal request for %s: %v", key, err)
	}

	var lastErr error
	for attempt := 0; attempt < MAXRETRY; attempt++ {
		if attempt > 0 {
			time.Sleep(time.Duration(attempt) * time.Second)
		}

		ctx, cancel := context.WithTimeout(context.Background(), constant.TIME_ONE_MINUTE)
		result, err := kc.doRequest(ctx, url, bodyByte, key)
		cancel()

		if err == nil {
			return result, nil
		}
		lastErr = err
	}

	return "", fmt.Errorf("all %d attempts failed for key %s: %w", MAXRETRY, key, lastErr)
}

func (kc *SettingLibClient) doRequest(ctx context.Context, url string, body []byte, key string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		return "", fmt.Errorf("failed to create request for %s: %v", key, err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := kc.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to call kunci service for %s: %v", key, err)
	}
	defer func() {
		// Drain body to allow TCP connection reuse before closing.
		io.Copy(io.Discard, resp.Body) //nolint:errcheck
		resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("kunci service returned %s for key %s", resp.Status, key)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response for %s: %v", key, err)
	}

	result := string(bodyBytes)
	// Treat a "Timeout" body as an error so the retry loop can handle it.
	if strings.Contains(result, "Timeout") {
		return "", fmt.Errorf("kunci service timeout response for key %s", key)
	}

	return result, nil
}
