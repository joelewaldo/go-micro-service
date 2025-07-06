package handler

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/joelewaldo/go-micro-service/internal/oauth2/config"
	"github.com/joelewaldo/go-micro-service/internal/oauth2/db"

	"github.com/golang-jwt/jwt/v5"
)

type jsonTokenReq struct {
	ClientID     string   `json:"client_id"`
	ClientSecret string   `json:"client_secret"`
	Scopes       []string `json:"scopes"`
}

type TokenRequest struct {
	ClientID     string   `json:"client_id"`
	ClientSecret string   `json:"client_secret"`
	Scopes       []string `json:"scopes"`
}

type TokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

type TokenHandler struct {
	DB     *sql.DB
	Config *config.Config
}

func NewTokenHandler(cfg *config.Config, db *sql.DB) *TokenHandler {
	return &TokenHandler{DB: db, Config: cfg}
}

func (h *TokenHandler) Handle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var clientID, clientSecret string
	var scopes []string

	ct := r.Header.Get("Content-Type")
	switch {
	case strings.HasPrefix(ct, "application/json"):
		var jr jsonTokenReq
		if err := json.NewDecoder(r.Body).Decode(&jr); err != nil {
			http.Error(w, "invalid JSON", http.StatusBadRequest)
			return
		}
		clientID, clientSecret, scopes = jr.ClientID, jr.ClientSecret, jr.Scopes

	case strings.HasPrefix(ct, "application/x-www-form-urlencoded"):
		if err := r.ParseForm(); err != nil {
			http.Error(w, "invalid form", http.StatusBadRequest)
			return
		}
		if r.FormValue("grant_type") != "client_credentials" {
			http.Error(w, "unsupported_grant_type", http.StatusBadRequest)
			return
		}
		if s := strings.TrimSpace(r.FormValue("scope")); s != "" {
			scopes = strings.Fields(s)
		}
		if u, p, ok := r.BasicAuth(); ok {
			clientID, clientSecret = u, p
		} else {
			clientID, clientSecret = r.FormValue("client_id"), r.FormValue("client_secret")
		}

	default:
		http.Error(w, "unsupported content type", http.StatusBadRequest)
		return
	}

	client, err := db.GetClientByID(h.DB, clientID)
	if err != nil {
		http.Error(w, "invalid_client", http.StatusUnauthorized)
		return
	}

	if !db.VerifyClientSecret(client, clientSecret) {
		http.Error(w, "invalid_client", http.StatusUnauthorized)
		return
	}

	for _, reqScope := range scopes {
		ok := false
		for _, allowed := range client.Scopes {
			if reqScope == allowed {
				ok = true
				break
			}
		}
		if !ok {
			http.Error(w, "invalid_scope", http.StatusBadRequest)
			return
		}
	}

	now := time.Now().UTC()
	exp := now.Add(time.Hour)
	claims := jwt.MapClaims{
		"iss":    h.Config.ISSUER,
		"sub":    client.ClientID,
		"iat":    now.Unix(),
		"exp":    exp.Unix(),
		"scopes": scopes,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	signed, err := token.SignedString(h.Config.RSA_KEY)
	if err != nil {
		http.Error(w, "server_error", http.StatusInternalServerError)
		return
	}

	resp := TokenResponse{
		AccessToken: signed,
		TokenType:   "Bearer",
		ExpiresIn:   int(exp.Sub(now).Seconds()),
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}
