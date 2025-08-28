package main

import (
	"log"
	"net/http"

	"ect-borrow-be/internal/config"
	"ect-borrow-be/internal/repository"
	"ect-borrow-be/internal/transport/http"
)

func main() {
	cfg := config.Load()

	repo, err := repository.NewSheetsRepository(cfg)
	if err != nil {
		log.Fatalf("init repository: %v", err)
	}

	r := transport.NewRouter(repo, cfg)

	addr := ":" + cfg.Port
	log.Printf("listening on %s", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatal(err)
	}
}
