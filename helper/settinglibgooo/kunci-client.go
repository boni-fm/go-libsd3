package settinglibgooo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/boni-fm/go-libsd3/config"
)

/*
	TODO:
*/

var prefix = "mujiyono"
var Baseurl = "localhost"
var MAXRETRY = 3

type SettingLibClient struct {
	httpClient *http.Client
	key        string
}

type Params struct {
	Key string `json:"key"`
}

func NewSettingLibClient(kunci string) *SettingLibClient {
	return &SettingLibClient{
		httpClient: &http.Client{
			Timeout: config.TIME_ONE_MINUTE,
		},
		key: kunci,
	}
}

func (kc *SettingLibClient) GetVariable(key string) (string, error) {
	kunciIpEnv := os.Getenv("KUNCI_IP_DOMAIN")
	if kunciIpEnv != "" {
		Baseurl = kunciIpEnv
	}

	url := "http://" + Baseurl
	if kc.key != "" {
		url += "/" + kc.key
	}

	url += "/GetVariabel"
	bodyReq := prefix + key
	params := Params{Key: bodyReq}
	bodyByte, _ := json.Marshal(params)

	req, errReq := http.NewRequest("POST", url, bytes.NewBuffer(bodyByte))
	if errReq != nil {
		return "", fmt.Errorf("failed to create request to kunci service %s : %v", key, errReq)
	}

	req.Header.Set("Content-Type", "application/json")
	resp, errResp := kc.httpClient.Do(req)
	if errResp != nil {
		return "", fmt.Errorf("failed to hit kunci service %s : %v", key, errResp)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to get variable: %s", resp.Status)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %v", err)
	}

	result := string(bodyBytes)
	if strings.Contains(result, "Timeout") {
		result = strings.Split(result, ";")[0]
	}

	return result, nil
}
