package kunci

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

/*
	TODO:
	- buat GetVariable baca info constring dari kunci, miripin dengan settinglibb
*/

var prefix = "mujiyono"
var Baseurl = "localhost"

type KunciClient struct {
	httpClient *http.Client
	Kunci      string
}

type Params struct {
	Key string `json:"key"`
}

func NewKunciClient(kunci string) *KunciClient {
	log.Say("buat kunci client")
	return &KunciClient{
		httpClient: &http.Client{},
		Kunci:      kunci,
	}
}

func (kc *KunciClient) GetVariable(key string, pathKunci string) (string, error) {
	kunciIpEnv := os.Getenv("KUNCI_IP_DOMAIN")
	if kunciIpEnv != "" {
		Baseurl = kunciIpEnv
	}

	url := "http://" + Baseurl
	if pathKunci != "" {
		url += "/" + pathKunci
	}

	url += "/GetVariabel"
	bodyReq := prefix + key
	params := Params{Key: bodyReq}
	bodyByte, _ := json.Marshal(params)

	req, errReq := http.NewRequest("POST", url, bytes.NewBuffer(bodyByte))
	if errReq != nil {
		return "", errReq
	}

	req.Header.Set("Content-Type", "application/json")
	resp, errResp := kc.httpClient.Do(req)
	if errResp != nil {
		return "", errResp
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to get variable: %s", resp.Status)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	result := string(bodyBytes)
	if strings.Contains(result, "Timeout") {
		result = strings.Split(result, ";")[0]
	}

	return result, nil
}
