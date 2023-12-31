package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/carlmjohnson/versioninfo"
	"github.com/dossif/gosocks5/internal/app"
	"github.com/dossif/gosocks5/internal/config"
	"github.com/dossif/gosocks5/pkg/logger"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

const (
	appName      = "gosocks5"
	configPrefix = "gosocks5"
)

func main() {
	var help = flag.Bool("h", false, "print usage and exit")
	flag.Parse()
	if *help == true {
		config.PrintUsage(configPrefix)
	}
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
	version := fmt.Sprintf("%v (%v)", versioninfo.Short(), versioninfo.LastCommit.Format(time.RFC3339))
	app.Run(ctx, wg, cfg, lg, appName, version)
}
