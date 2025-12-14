package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/YurcheuskiRadzivon/disk-diag/internal/config"
	"github.com/YurcheuskiRadzivon/disk-diag/internal/server"
	"github.com/YurcheuskiRadzivon/disk-diag/internal/service/base"
	"github.com/YurcheuskiRadzivon/disk-diag/internal/service/benchmark"
	"github.com/YurcheuskiRadzivon/disk-diag/internal/service/diagnostic"
	"github.com/YurcheuskiRadzivon/disk-diag/internal/service/smart"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	run(cfg, ctx)
}

func run(cfg *config.Config, ctx context.Context) {
	base, err := base.NewService(ctx)
	if err != nil {
		log.Fatalf("Base: %v", err)
	}

	smart, err := smart.NewService(ctx)
	if err != nil {
		log.Fatalf("Smart: %v", err)
	}

	benchmark, err := benchmark.NewService(ctx)
	if err != nil {
		log.Fatalf("Benchmark: %v", err)
	}

	diagnostic, err := diagnostic.NewService(ctx, cfg.DIAGNOSTIC.API_KEY)
	if err != nil {
		log.Fatalf("Diagnosis: %v", err)
	}
	defer diagnostic.Close()

	srv := server.New(cfg.HTTP.PORT, base, smart, benchmark, diagnostic)

	srv.RegisterRoutes()

	srv.Start()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case <-interrupt:
		log.Println("Shutdown")

	case err := <-srv.Notify():
		log.Panicf("server: %s", err)
	}

	err = srv.Shutdown()
	if err != nil {
		log.Fatalf("server: %v", err)
	}
}
