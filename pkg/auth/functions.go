package helper

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/boni-fm/go-libsd3/pkg/settinglibgo"
	"github.com/golang-jwt/jwt/v5"
)

var (
	cachedToken    string
	cachedTokenExp int64
	tokenMutex     sync.Mutex
)

// =====  PANGGIL API UNTUK DAPATKAN TOKEN DARI AWS =====
func GetTokenFromAuthAPIAWS() (string, error) {
	ApiAuthUrl := os.Getenv("API_URL")

	req, err := http.NewRequest(
		http.MethodGet,
		ApiAuthUrl,
		nil,
	)
	if err != nil {
		return "", err
	}

	encryptedUsername := os.Getenv("API_USERNAME")
	encryptedPassword := os.Getenv("API_PASSWORD")

	req.SetBasicAuth(encryptedUsername, encryptedPassword)

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("login-auth gagal, status: %d", resp.StatusCode)
	}

	var token string
	if err := json.NewDecoder(resp.Body).Decode(&token); err != nil {
		return "", err
	}

	if token == "" {
		return "", errors.New("token kosong dari auth api")
	}

	return token, nil
}

func FetchTokenAuthFromSettingLib(KodeDc string) (string, error) {
	settingLibClient := settinglibgo.NewSettingLibClient(KodeDc)
	ApiAuthUrl, err := settingLibClient.GetVariable("BaseUrlCloud")
	if err != nil {
		return "", fmt.Errorf("[BaseUrlCloud] base url pada kunci kosong, cek kembali SettingWeb.xml pada kunci")
	}

	ApiAuthUrl = ApiAuthUrl + "/apicloudauth/login-auth"
	req, err := http.NewRequest(
		http.MethodGet,
		ApiAuthUrl,
		nil,
	)
	if err != nil {
		return "", err
	}

	encryptedUsername := os.Getenv("API_USERNAME")
	encryptedPassword := os.Getenv("API_PASSWORD")
	if encryptedUsername == "" || encryptedPassword == "" {
		return "", errors.New("API_USERNAME atau API_PASSWORD tidak ditemukan di environment variables")
	}

	req.SetBasicAuth(encryptedUsername, encryptedPassword)

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("login-auth gagal, status: %d", resp.StatusCode)
	}

	var token string
	if err := json.NewDecoder(resp.Body).Decode(&token); err != nil {
		return "", err
	}

	if token == "" {
		return "", errors.New("token kosong dari auth api")
	}

	return token, nil
}

func FetchTokenAuthFromUrl(ApiAuthUrl string) (string, error) {
	if ApiAuthUrl == "" {
		return "", errors.New("ApiAuthUrl kosong, pastikan parameter url auth api tidak kosong")
	}

	req, err := http.NewRequest(
		http.MethodGet,
		ApiAuthUrl,
		nil,
	)
	if err != nil {
		return "", err
	}

	encryptedUsername := os.Getenv("API_USERNAME")
	encryptedPassword := os.Getenv("API_PASSWORD")
	if encryptedUsername == "" || encryptedPassword == "" {
		return "", errors.New("API_USERNAME atau API_PASSWORD tidak ditemukan di environment variables")
	}

	req.SetBasicAuth(encryptedUsername, encryptedPassword)

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("login-auth gagal, status: %d", resp.StatusCode)
	}

	var token string
	if err := json.NewDecoder(resp.Body).Decode(&token); err != nil {
		return "", err
	}

	if token == "" {
		return "", errors.New("token kosong dari auth api")
	}

	return token, nil
}

// ===== CEK TOKEN EXPIRED & AMBIL TOKEN YANG VALID (BLM EXP) =====
func GetValidToken(ctx context.Context) (string, error) {
	tokenMutex.Lock()
	defer tokenMutex.Unlock()

	// kalau token masih ada & belum expired → langsung pakai
	if cachedToken != "" && time.Now().Unix() < cachedTokenExp {
		return cachedToken, nil
	}

	// selain itu → ambil token baru
	newToken, err := GetTokenFromAuthAPIAWS()
	if err != nil {
		return "", err
	}

	expired, exp, err := IsTokenExpired(newToken)
	if err != nil {
		return "", err
	}
	if expired {
		return "", errors.New("token baru langsung expired")
	}

	cachedToken = newToken
	cachedTokenExp = exp

	return cachedToken, nil
}

func IsTokenExpired(tokenString string) (bool, int64, error) {
	parser := jwt.NewParser(jwt.WithoutClaimsValidation())

	token, _, err := parser.ParseUnverified(tokenString, jwt.MapClaims{})
	if err != nil {
		return true, 0, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return true, 0, errors.New("claims bukan MapClaims")
	}

	expVal, ok := claims["exp"]
	if !ok {
		return true, 0, errors.New("claim exp tidak ditemukan")
	}

	exp, ok := expVal.(float64)
	if !ok {
		return true, 0, errors.New("exp bukan float64")
	}

	return time.Now().Unix() > int64(exp), int64(exp), nil
}

// ===== AMBIL SECRET JWT DARI AWS SECRETS MANAGER & VALIDASI TOKEN AWS =====
func ValidateTokenAWS(ctx context.Context, tokenString string) (bool, error) {
	if tokenString == "" {
		return false, errors.New("token kosong")
	}

	jwtSecret, err := GetJWTSecretFromAWS(ctx)
	if err != nil {
		return false, err
	}

	token, err := jwt.Parse(
		tokenString,
		func(token *jwt.Token) (interface{}, error) {

			if token.Method.Alg() != jwt.SigningMethodHS512.Alg() {
				return nil, fmt.Errorf("unexpected signing method: %s", token.Method.Alg())
			}

			return []byte(jwtSecret), nil
		},
		jwt.WithExpirationRequired(),
		jwt.WithLeeway(0),
	)

	if err != nil {
		return false, err
	}

	return token.Valid, nil
}

func GetJWTSecretFromAWS(ctx context.Context) (string, error) {
	region := os.Getenv("AWS_REGION")
	secretName := os.Getenv("AWS_SECRET_NAME")

	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		return "", err
	}

	client := secretsmanager.NewFromConfig(cfg)

	out, err := client.GetSecretValue(ctx, &secretsmanager.GetSecretValueInput{
		SecretId: &secretName,
	})
	if err != nil {
		return "", err
	}

	var data map[string]string
	if err := json.Unmarshal([]byte(*out.SecretString), &data); err != nil {
		return "", err
	}

	secret := data["AppSettings__Token"]
	if secret == "" {
		return "", errors.New("AppSettings__Token tidak ditemukan di AWS")
	}

	return secret, nil
}
