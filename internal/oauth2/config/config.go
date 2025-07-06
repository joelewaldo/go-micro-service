package config

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"os"
	"strings"
)

type Config struct {
	LOG_LEVEL    string
	DATABASE_URL string
	HOST         string
	PORT         string
	ISSUER       string
	RSA_KEY      *rsa.PrivateKey
}

func ExtractConfig() (*Config, error) {
	keySource := os.Getenv("PRIVATE_KEY")
	if file := os.Getenv("PRIVATE_KEY_FILE"); file != "" {
		keySource = file
	}

	rsaKey, err := loadPrivateKey(keySource)
	if err != nil {
		return nil, err
	}
	databaseUrl, err := envOrRequired("DATABASE_URL")
	if err != nil {
		return nil, err
	}

	host := envOrDefault("HOST", "localhost")
	port := envOrDefault("PORT", "7000")
	issuer := host + ":" + port

	return &Config{
		LOG_LEVEL:    envOrDefault("LOG_LEVEL", "info"),
		DATABASE_URL: databaseUrl,
		HOST:         host,
		PORT:         port,
		ISSUER:       issuer,
		RSA_KEY:      rsaKey,
	}, nil
}

func envOrDefault(key, defaultVal string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}

	return defaultVal
}

func envOrRequired(key string) (string, error) {
	if v := os.Getenv(key); v != "" {
		return v, nil
	}

	return "", fmt.Errorf("key, %s, not set", key)
}

func loadPrivateKey(pemOrPath string) (*rsa.PrivateKey, error) {
	if pemOrPath == "" {
		return nil, errors.New("private key not set")
	}

	var data []byte
	if info, err := os.Stat(pemOrPath); err == nil && !info.IsDir() {
		data, err = os.ReadFile(pemOrPath)
		if err != nil {
			return nil, fmt.Errorf("could not read private key file %q: %w", pemOrPath, err)
		}
	} else {
		data = []byte(pemOrPath)
	}

	block, _ := pem.Decode(data)
	if block == nil || !strings.Contains(block.Type, "PRIVATE KEY") {
		return nil, errors.New("failed to parse PEM block containing the private key")
	}

	if parsed, err := x509.ParsePKCS8PrivateKey(block.Bytes); err == nil {
		if key, ok := parsed.(*rsa.PrivateKey); ok {
			return key, nil
		}
		return nil, errors.New("parsed PKCS#8 key is not an RSA private key")
	}

	if pk1, err := x509.ParsePKCS1PrivateKey(block.Bytes); err == nil {
		return pk1, nil
	}

	return nil, errors.New("unable to parse RSA private key (tried PKCS#8 and PKCS#1)")
}
