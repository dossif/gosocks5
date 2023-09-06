package main

import (
	"context"
	"github.com/dossif/gosocks5/internal/app"
	"github.com/dossif/gosocks5/internal/config"
	"github.com/dossif/gosocks5/pkg/logger"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

const (
	appName      = "gosocks5"
	configPrefix = "gosocks5"
)

var appVersion = "0.0.0"

func main() {
	wg := new(sync.WaitGroup)
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg, err := config.NewConfig(configPrefix)
	if err != nil {
		log.Printf("failed to create config: %v", err)
		config.PrintUsage(configPrefix)
	}
	lg, err := logger.NewLogger(cfg.LogLevel)
	if err != nil {
		log.Fatalf("failed to create logger: %v", err)
	}
	app.Run(ctx, wg, cfg, lg, appName, appVersion)
}
