package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/YurcheuskiRadzivon/disk-diag/internal/config"
	"github.com/YurcheuskiRadzivon/disk-diag/internal/server"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}
	run(cfg)
}

func run(cfg *config.Config) {

	srv := server.New(cfg.HTTP.PORT)

	srv.Start()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case <-interrupt:
		log.Println("Shutdown")

	case err := <-srv.Notify():
		log.Panicf("Httpserver: %s", err)
	}

	err := srv.Shutdown()
	if err != nil {
		log.Fatalf("Httpserver: %v", err)
	}
}
