package main

import (
	"log"

	"es-meili-growi-bridge/internal/config"
	"es-meili-growi-bridge/internal/server"
)

func main() {
	cfg := config.Load()

	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	log.Printf("es-meili-growi-bridge v0.1.0")
	log.Printf("Meilisearch URL: %s", cfg.MeilisearchURL)
	log.Printf("Listen address: %s", cfg.ListenAddr)

	srv := server.New(cfg)
	if err := srv.Start(); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
