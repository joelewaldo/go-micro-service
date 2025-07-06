package config

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"strings"
)

type Config struct {
	LOG_LEVEL string
	HOST      string
	PORT      string
	ISSUER    string
	RSA_KEY   *rsa.PublicKey
}

func ExtractConfig() (*Config, error) {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		return nil, fmt.Errorf("environment variable DATABASE_URL is required")
	}

	host := envOrDefault("HOST", "localhost")
	port := envOrDefault("PORT", "7979")
	issuer := fmt.Sprintf("http://%s:%s", host, port)

	logLevel := envOrDefault("LOG_LEVEL", "info")

	keySrc := os.Getenv("PUBLIC_KEY")
	if file := os.Getenv("PUBLIC_KEY_FILE"); file != "" {
		keySrc = file
	}
	pub, err := loadPublicKey(keySrc)
	if err != nil {
		return nil, err
	}

	return &Config{
		LOG_LEVEL: logLevel,
		HOST:      host,
		PORT:      port,
		ISSUER:    issuer,
		RSA_KEY:   pub,
	}, nil
}

func envOrDefault(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func loadPublicKey(pemOrPath string) (*rsa.PublicKey, error) {
	if pemOrPath == "" {
		return nil, fmt.Errorf("public key not set (via PUBLIC_KEY or PUBLIC_KEY_FILE)")
	}

	var data []byte
	if fi, err := os.Stat(pemOrPath); err == nil && !fi.IsDir() {
		data, err = os.ReadFile(pemOrPath)
		if err != nil {
			return nil, fmt.Errorf("failed to read public key file %q: %w", pemOrPath, err)
		}
	} else {
		data = []byte(pemOrPath)
	}

	block, _ := pem.Decode(data)
	if block == nil || !strings.Contains(block.Type, "PUBLIC KEY") {
		return nil, fmt.Errorf("failed to parse PEM block containing the public key")
	}

	if parsed, err := x509.ParsePKIXPublicKey(block.Bytes); err == nil {
		if rsaPub, ok := parsed.(*rsa.PublicKey); ok {
			return rsaPub, nil
		}
	}

	if rsaPub1, err := x509.ParsePKCS1PublicKey(block.Bytes); err == nil {
		return rsaPub1, nil
	}

	return nil, fmt.Errorf("failed to parse RSA public key in PKIX or PKCS#1 format")
}
