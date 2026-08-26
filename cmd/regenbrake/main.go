package main

import (
	"log"
	"net/http"
)

func main() {
	cfg := parseConfig()
	service := NewService(cfg)
	defer service.Close()
	if err := service.SelfCheck(); err != nil {
		log.Fatalf("selfcheck failed: %v", err)
	}
	mux := newServer(service)
	log.Printf("regenbrake listening on %s", cfg.Addr)
	if err := http.ListenAndServe(cfg.Addr, mux); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
