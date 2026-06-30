package config

import (
	"bufio"
	"os"
	"strings"
)

// Config menampung seluruh konfigurasi aplikasi yang diambil dari environment
// variable / file .env. Dengan ini DSN database dan JWT secret tidak lagi
// hardcoded di beberapa tempat berbeda (sebelumnya JWT_SECRET didefinisikan
// dua kali di usecase & middleware, rawan tidak sinkron).
type Config struct {
	DBDSN     string
	JWTSecret []byte
	Port      string
}

// Load membaca file .env (jika ada) lalu mengembalikan Config.
// Environment variable yang sudah di-set di OS akan selalu diprioritaskan
// dibanding isi file .env.
func Load() *Config {
	loadEnvFile(".env")

	return &Config{
		DBDSN:     getEnv("DB_DSN", "root:password@tcp(127.0.0.1:3306)/golang_rbac_db"),
		JWTSecret: []byte(getEnv("JWT_SECRET", "RAHASIA_SUPER_SECURE_123")),
		Port:      getEnv("PORT", "8080"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

// loadEnvFile mem-parsing file .env sederhana (format KEY=VALUE per baris,
// baris kosong / diawali '#' diabaikan) lalu menyetelnya sebagai environment
// variable proses ini. Tidak menimpa env var yang sudah ada.
func loadEnvFile(path string) {
	f, err := os.Open(path)
	if err != nil {
		return // .env opsional, kalau tidak ada ya pakai default / env OS
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.Trim(strings.TrimSpace(parts[1]), `"'`)

		if _, exists := os.LookupEnv(key); !exists {
			_ = os.Setenv(key, value)
		}
	}
}
