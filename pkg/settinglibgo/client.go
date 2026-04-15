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
		httpClient: &http.Client{
			Timeout: constant.TIME_FIVE_MINUTES,
		},
		key: key,
	}
}

func (kc *SettingLibClient) GetVariable(key string) (string, error) {
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

		req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(bodyByte))
		if err != nil {
			// Request creation failure is not retryable.
			return "", fmt.Errorf("failed to create request for %s: %v", key, err)
		}
		req.Header.Set("Content-Type", "application/json")

		resp, err := kc.httpClient.Do(req)
		if err != nil {
			lastErr = fmt.Errorf("failed to call kunci service for %s: %v", key, err)
			continue
		}

		if resp.StatusCode != http.StatusOK {
			io.Copy(io.Discard, resp.Body) // drain to allow connection reuse
			resp.Body.Close()
			lastErr = fmt.Errorf("kunci service returned %s for key %s", resp.Status, key)
			continue
		}

		bodyBytes, readErr := io.ReadAll(resp.Body)
		resp.Body.Close()

		if readErr != nil {
			lastErr = fmt.Errorf("failed to read response for %s: %v", key, readErr)
			continue
		}

		result := string(bodyBytes)
		if strings.Contains(result, "Timeout") {
			lastErr = fmt.Errorf("kunci service timeout response for key %s", key)
			continue
		}

		return result, nil
	}

	return "", fmt.Errorf("all %d attempts failed for key %s: %w", MAXRETRY, key, lastErr)
}
