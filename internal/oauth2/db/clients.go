package db

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/joelewaldo/go-micro-service/pkg/models"
	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

func GetClientByID(db *sql.DB, clientID string) (*models.Client, error) {
	var c models.Client
	row := db.QueryRow(`SELECT id, client_id, client_secret, scopes, is_active
                        FROM oauth_clients
                        WHERE client_id = $1`, clientID)

	if err := row.Scan(&c.ID, &c.ClientID, &c.HashedSecret, pq.Array(&c.Scopes), &c.IsActive); err != nil {
		return nil, err
	}

	if !c.IsActive {
		return nil, errors.New("client is inactive")
	}

	return &c, nil
}

func CreateClient(db *sql.DB, clientID, plainSecret string, scopes []string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plainSecret), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("bcrypt: %w", err)
	}
	_, err = db.Exec(`
      INSERT INTO oauth_clients (client_id, client_secret, scopes, is_active)
      VALUES ($1,$2,$3, TRUE)
    `, clientID, string(hash), pq.Array(scopes))
	if err != nil {
		return fmt.Errorf("insert oauth_clients: %w", err)
	}
	return nil
}

func RegisterClient(db *sql.DB, scopes []string) (string, string, error) {
	clientID := uuid.NewString()

	raw := make([]byte, 32)
	if _, err := rand.Read(raw); err != nil {
		return "", "", fmt.Errorf("failed to generate secret: %w", err)
	}
	clientSecret := base64.RawURLEncoding.EncodeToString(raw)

	if err := CreateClient(db, clientID, clientSecret, scopes); err != nil {
		return "", "", fmt.Errorf("CreateClient: %w", err)
	}

	return clientID, clientSecret, nil
}

func VerifyClientSecret(client *models.Client, providedSecret string) bool {
	return bcrypt.CompareHashAndPassword([]byte(client.HashedSecret), []byte(providedSecret)) == nil
}
