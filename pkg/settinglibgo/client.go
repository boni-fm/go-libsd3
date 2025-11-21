package settinglibgo

/*

	Http Client untuk manggil kunci
	Note:
	- manggil POST request /GetVariabel dari apiservice kunci
	- IP apikunci didapatkan dari env KUNCI_IP_DOMAIN -> docker-hub-nginx-1
	  > jika env kosong, maka default localhost
	- body nya make prefix "mujiyono"
	- ada 3 kali proses retry jika panggilan pertama gagal

	Err:
	- kadang masih sering kena i/o timeout
	  > solusi sementara restart service aplikasi yang implementasi

	TODO:
	- cari solusi permanen untuk masalah i/o timeout

*/

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/boni-fm/go-libsd3/pkg/config/constant"
)

var (
	PREFIX   = "mujiyono"
	BASEURL  = "localhost"
	MAXRETRY = 3
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
	return &SettingLibClient{
		httpClient: &http.Client{
			Timeout: constant.TIME_ONE_MINUTE,
		},
		key: kunci,
	}
}

func (kc *SettingLibClient) GetVariable(key string) (string, error) {
	kunciIpEnv := os.Getenv("KUNCI_IP_DOMAIN")
	if kunciIpEnv != "" {
		BASEURL = kunciIpEnv
	}

	url := "http://" + BASEURL
	if kc.key != "" {
		url += "/" + kc.key
	}

	url += "/GetVariabel"
	bodyReq := PREFIX + key
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
