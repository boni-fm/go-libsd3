package example

// ======================================================
// Contoh pemakaian package auth (pkg/auth)
// ======================================================
//
// Package auth nyediain fungsi untuk autentikasi berbasis JWT dan AWS.
//
// Fitur yang di-demo:
//   - GetTokenFromAuthAPIAWS — dapetin JWT token dari API auth yang pakai AWS credentials
//   - FetchTokenAuthFromSettingLib — dapetin token dengan URL dari kunci service
//   - FetchTokenAuthFromUrl — dapetin token dari URL auth yang diberikan langsung
//   - GetValidToken — dapetin token yang masih valid (dengan caching, thread-safe)
//   - IsTokenExpired — cek apakah JWT token sudah expired
//   - ValidateTokenAWS — validasi tanda tangan JWT token dengan secret dari AWS
//   - GetJWTSecretFromAWS — ambil JWT secret dari AWS Secrets Manager
//
// Catatan penting:
//   - Semua fungsi di sini butuh environment variables buat jalan beneran.
//   - GetValidToken pakai in-memory cache (thread-safe), jadi aman dipanggil
//     dari banyak goroutine sekaligus.
//   - Env vars yang dibutuhkan:
//       API_URL           — URL endpoint auth API (untuk GetTokenFromAuthAPIAWS)
//       API_USERNAME      — username yang diencrypt
//       API_PASSWORD      — password yang diencrypt
//       AWS_REGION        — region AWS (untuk ValidateTokenAWS/GetJWTSecretFromAWS)
//       AWS_SECRET_NAME   — nama secret di AWS Secrets Manager
//       KUNCI_IP_DOMAIN   — (opsional) IP domain service kunci

import (
	"context"
	"fmt"
	"os"

	helper "github.com/boni-fm/go-libsd3/pkg/auth"
)

// ContohGetTokenFromAuthAPIAWS mendemonstrasikan pengambilan token dari API auth.
// API auth dipanggil dengan Basic Auth menggunakan API_USERNAME & API_PASSWORD.
func ContohGetTokenFromAuthAPIAWS() {
	// Set env vars yang dibutuhkan
	os.Setenv("API_URL", "http://auth.internal/api/login-auth")
	os.Setenv("API_USERNAME", "encrypted_username_base64")
	os.Setenv("API_PASSWORD", "encrypted_password_base64")
	defer os.Unsetenv("API_URL")
	defer os.Unsetenv("API_USERNAME")
	defer os.Unsetenv("API_PASSWORD")

	token, err := helper.GetTokenFromAuthAPIAWS()
	if err != nil {
		fmt.Printf("gagal dapetin token: %v\n", err)
		return
	}

	fmt.Printf("token diperoleh (panjang: %d char)\n", len(token))
	// Token adalah JWT string, contoh: eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9...
}

// ContohFetchTokenAuthFromSettingLib mendemonstrasikan pengambilan token
// di mana URL-nya diambil dari service kunci (BaseUrlCloud).
func ContohFetchTokenAuthFromSettingLib() {
	os.Setenv("API_USERNAME", "encrypted_user")
	os.Setenv("API_PASSWORD", "encrypted_pass")
	defer os.Unsetenv("API_USERNAME")
	defer os.Unsetenv("API_PASSWORD")

	// KodeDc dipakai buat nyari kunci yang menyimpan BaseUrlCloud
	// Fungsi akan manggil: http://kunci<KodeDc>/GetVariabel dengan key="BaseUrlCloud"
	// Lalu hit URL tersebut + "/apicloudauth/login-auth" buat dapetin token
	token, err := helper.FetchTokenAuthFromSettingLib("G009SIM")
	if err != nil {
		fmt.Printf("gagal fetch token dari settinglib: %v\n", err)
		fmt.Println("pastikan service kunci jalan dan BaseUrlCloud sudah diset")
		return
	}

	fmt.Printf("token dari settinglib berhasil (panjang: %d)\n", len(token))
}

// ContohFetchTokenAuthFromUrl mendemonstrasikan pengambilan token
// dari URL auth yang sudah diketahui langsung.
func ContohFetchTokenAuthFromUrl() {
	os.Setenv("API_USERNAME", "encrypted_user")
	os.Setenv("API_PASSWORD", "encrypted_pass")
	defer os.Unsetenv("API_USERNAME")
	defer os.Unsetenv("API_PASSWORD")

	// Kalau udah tau URL auth-nya, bisa langsung kasih
	authURL := "http://my-auth-service.internal/api/v2/login-auth"

	token, err := helper.FetchTokenAuthFromUrl(authURL)
	if err != nil {
		fmt.Printf("gagal fetch token dari URL: %v\n", err)
		return
	}

	fmt.Printf("token diperoleh dari URL langsung (panjang: %d)\n", len(token))
}

// ContohGetValidToken mendemonstrasikan GetValidToken yang punya caching otomatis.
//
// GetValidToken:
//   1. Cek apakah token yang di-cache masih valid (belum expired)
//   2. Kalau masih valid → langsung return dari cache (TANPA hit network)
//   3. Kalau expired → ambil token baru dari API, simpan ke cache
//
// Thread-safe: aman dipanggil dari banyak goroutine sekaligus.
func ContohGetValidToken() {
	os.Setenv("API_URL", "http://auth.internal/login")
	os.Setenv("API_USERNAME", "user")
	os.Setenv("API_PASSWORD", "pass")
	defer os.Unsetenv("API_URL")
	defer os.Unsetenv("API_USERNAME")
	defer os.Unsetenv("API_PASSWORD")

	ctx := context.Background()

	// Panggil pertama kali → hit API, simpan ke cache
	token1, err := helper.GetValidToken(ctx)
	if err != nil {
		fmt.Printf("GetValidToken pertama gagal: %v\n", err)
		return
	}
	fmt.Printf("token pertama (dari API): %s...\n", token1[:20])

	// Panggil kedua kalinya → dari cache, ga hit network
	token2, err := helper.GetValidToken(ctx)
	if err != nil {
		fmt.Printf("GetValidToken kedua gagal: %v\n", err)
		return
	}
	fmt.Printf("token kedua (dari cache): %s...\n", token2[:20])

	// token1 == token2 (sama persis karena dari cache)
	fmt.Printf("token sama? %v\n", token1 == token2)
}

// ContohIsTokenExpired mendemonstrasikan pengecekan expiry JWT token.
// Bisa dipakai tanpa koneksi apapun — parsing JWT dilakukan lokal.
func ContohIsTokenExpired() {
	// Contoh JWT token (ini adalah token dummy yang sudah expired)
	// Di produksi, ini adalah token yang kamu dapat dari GetValidToken/GetTokenFromAuthAPIAWS
	expiredToken := "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9." +
		"eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiZXhwIjoxNjAwMDAwMDAwfQ." +
		"dummysignature"

	expired, expTime, err := helper.IsTokenExpired(expiredToken)
	if err != nil {
		// Error bisa terjadi kalau format token tidak valid / bukan JWT
		fmt.Printf("gagal cek expiry: %v\n", err)
		return
	}

	if expired {
		fmt.Printf("token sudah expired di Unix timestamp: %d\n", expTime)
	} else {
		fmt.Printf("token masih valid sampai Unix timestamp: %d\n", expTime)
	}

	// Contoh penggunaan dalam middleware auth:
	// token := r.Header.Get("Authorization")
	// expired, _, err := helper.IsTokenExpired(token)
	// if err != nil || expired {
	//     http.Error(w, "token expired", http.StatusUnauthorized)
	//     return
	// }
}

// ContohValidateTokenAWS mendemonstrasikan validasi JWT token
// menggunakan secret yang disimpan di AWS Secrets Manager.
//
// Fungsi ini:
//  1. Ambil JWT secret dari AWS Secrets Manager (key: AppSettings__Token)
//  2. Validasi signature, expiry, dan algorithm token menggunakan secret tersebut
func ContohValidateTokenAWS() {
	ctx := context.Background()

	// Set env vars AWS
	os.Setenv("AWS_REGION", "ap-southeast-1")
	os.Setenv("AWS_SECRET_NAME", "my-app/jwt-secret")
	defer os.Unsetenv("AWS_REGION")
	defer os.Unsetenv("AWS_SECRET_NAME")

	// Token yang mau divalidasi (biasanya dari Authorization header)
	tokenString := "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIn0.dummysig"

	valid, err := helper.ValidateTokenAWS(ctx, tokenString)
	if err != nil {
		fmt.Printf("ValidateTokenAWS error: %v\n", err)
		// Error bisa karena: secret ga ditemukan di AWS, signature invalid, token expired
		return
	}

	if valid {
		fmt.Println("token valid dan belum expired!")
	} else {
		fmt.Println("token tidak valid")
	}
}

// ContohGetJWTSecretFromAWS mendemonstrasikan pengambilan JWT secret dari AWS.
//
// Secret disimpan sebagai JSON di AWS Secrets Manager:
// {"AppSettings__Token": "secret-key-yang-panjang-banget"}
func ContohGetJWTSecretFromAWS() {
	ctx := context.Background()

	os.Setenv("AWS_REGION", "ap-southeast-1")
	os.Setenv("AWS_SECRET_NAME", "production/my-api/secrets")
	defer os.Unsetenv("AWS_REGION")
	defer os.Unsetenv("AWS_SECRET_NAME")

	secret, err := helper.GetJWTSecretFromAWS(ctx)
	if err != nil {
		fmt.Printf("gagal ambil JWT secret dari AWS: %v\n", err)
		fmt.Println("pastikan AWS credentials sudah diset dan secret ada di Secrets Manager")
		return
	}

	// Jangan print secret di production! Ini cuma buat demo
	fmt.Printf("JWT secret berhasil diambil (panjang: %d char)\n", len(secret))
}

// ContohAuthLengkap mendemonstrasikan flow auth end-to-end.
// Biasanya dipanggil di middleware HTTP buat validasi setiap request.
func ContohAuthLengkap() {
	ctx := context.Background()

	// Step 1: Dapetin token yang valid (dengan caching)
	token, err := helper.GetValidToken(ctx)
	if err != nil {
		fmt.Printf("gagal dapetin token: %v\n", err)
		return
	}

	// Step 2: Cek dulu sebelum dipakai — token masih valid?
	expired, expUnix, err := helper.IsTokenExpired(token)
	if err != nil {
		fmt.Printf("gagal cek token: %v\n", err)
		return
	}
	if expired {
		fmt.Printf("token expired di %d, minta yang baru\n", expUnix)
		return
	}
	fmt.Printf("token valid sampai Unix %d\n", expUnix)

	// Step 3: Validasi signature token
	valid, err := helper.ValidateTokenAWS(ctx, token)
	if err != nil {
		fmt.Printf("validasi AWS gagal: %v\n", err)
		return
	}

	if valid {
		fmt.Println("token tervalidasi, request diizinkan lanjut")
	} else {
		fmt.Println("token tidak valid, request ditolak")
	}
}
