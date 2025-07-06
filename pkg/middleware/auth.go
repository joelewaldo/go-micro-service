package middleware

import (
	"context"
	"crypto/rsa"
	"errors"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type ctxKey string

const (
	ClaimsContextKey = ctxKey("claims")
)

func Auth(publicKey *rsa.PublicKey) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authz := r.Header.Get("Authorization")
			if authz == "" || !strings.HasPrefix(authz, "Bearer ") {
				http.Error(w, "missing or invalid Authorization header", http.StatusUnauthorized)
				return
			}
			tokenString := strings.TrimSpace(authz[len("Bearer "):])

			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
					return nil, errors.New("unexpected signing method")
				}
				return publicKey, nil
			})
			if err != nil || !token.Valid {
				http.Error(w, "invalid or expired token", http.StatusUnauthorized)
				return
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				http.Error(w, "invalid token claims", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), ClaimsContextKey, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func RequireScope(required string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ci := r.Context().Value(ClaimsContextKey)
			if ci == nil {
				http.Error(w, "missing token", http.StatusUnauthorized)
				return
			}
			claims, ok := ci.(jwt.MapClaims)
			if !ok {
				http.Error(w, "invalid token claims", http.StatusUnauthorized)
				return
			}

			var scopes []string
			if raw, found := claims["scopes"]; found {
				switch arr := raw.(type) {
				case []interface{}:
					for _, v := range arr {
						if s, ok := v.(string); ok {
							scopes = append(scopes, s)
						}
					}
				case []string:
					scopes = arr
				}
			}

			for _, s := range scopes {
				if s == required {
					next.ServeHTTP(w, r)
					return
				}
			}

			http.Error(w, "insufficient_scope", http.StatusForbidden)
		})
	}
}
