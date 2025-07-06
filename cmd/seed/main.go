package main

import (
	"flag"
	"fmt"
	"log"
	"strings"

	_ "github.com/lib/pq"

	"github.com/joelewaldo/go-micro-service/internal/oauth2/config"
	oauthdb "github.com/joelewaldo/go-micro-service/internal/oauth2/db"
)

func main() {
	var scopesCSV string
	flag.StringVar(&scopesCSV, "scopes", "", "Comma-separated scopes for the new client (e.g. \"read:users,write:orders\")")
	flag.Parse()

	cfg, err := config.ExtractConfig()
	if err != nil {
		log.Fatalf("config error: %v", err)
	}

	db, err := oauthdb.Connect(cfg)
	if err != nil {
		log.Fatalf("db connect error: %v", err)
	}
	defer db.Close()

	var scopes []string
	if scopesCSV != "" {
		scopes = strings.Split(scopesCSV, ",")
	}

	clientID, clientSecret, err := oauthdb.RegisterClient(db, scopes)
	if err != nil {
		log.Fatalf("could not register client: %v", err)
	}

	fmt.Println("=== New OAuth2 Client ===")
	fmt.Printf("client_id:     %s\n", clientID)
	fmt.Printf("client_secret: %s\n", clientSecret)
	fmt.Println("\nStore the secret somewhere safe! It will not be shown again.")
}
