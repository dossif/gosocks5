package app

import (
	"context"
	"fmt"
	"github.com/dossif/gosocks5/internal/config"
	"github.com/dossif/gosocks5/pkg/logger"
	"github.com/dossif/gosocks5/pkg/socks5"
	"sync"
)

const network = "tcp"

func Run(ctx context.Context, wg *sync.WaitGroup, cfg *config.Config, log *logger.Logger, appName string, appVersion string) {

	// Create a SOCKS5 server
	conf := &socks5.Config{Logger: log}
	server, err := socks5.New(conf)
	if err != nil {
		panic(err)
	}
	fmt.Printf("server listen %s://%s\n", cfg.Listen)
	if err := server.ListenAndServe(network, cfg.Listen); err != nil {
		panic(err)
	}
}
