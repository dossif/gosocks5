package app

import (
	"context"
	"github.com/dossif/gosocks5/internal/config"
	"github.com/dossif/gosocks5/internal/services"
	"github.com/dossif/gosocks5/pkg/logger"
	"sync"
)

func Run(ctx context.Context, wg *sync.WaitGroup, cfg *config.Config, log *logger.Logger, appName string, appVersion string) {
	log.Lg.Info().Msgf("start %v ver %v", appName, appVersion)
	svc, err := services.NewService(ctx, log, cfg.Listen, cfg.AuthMethod)
	if err != nil {
		log.Lg.Fatal().Msgf("failed to create service: %v", err)
	}
	wg.Add(1)
	go func() {
		err = svc.Start()
		if err != nil {
			log.Lg.Fatal().Msgf("failed to start service: %v", err)
		}
		defer wg.Done()
	}()
	wg.Wait()
	defer log.Lg.Info().Msgf("%v stopped", appName)
}
