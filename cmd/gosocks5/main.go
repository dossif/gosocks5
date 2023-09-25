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
)

const (
	appName      = "gosocks5"
	configPrefix = "gosocks5"
)

var appVersion = "0.0.0"

func main() {

	fmt.Println("ShortInfo:", versioninfo.Short())
	fmt.Println("Version:", versioninfo.Version)
	fmt.Println("Revision:", versioninfo.Revision)
	fmt.Println("DirtyBuild:", versioninfo.DirtyBuild)
	fmt.Println("LastCommit:", versioninfo.LastCommit)

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
	app.Run(ctx, wg, cfg, lg, appName, appVersion)
}
