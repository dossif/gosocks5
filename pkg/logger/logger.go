package logger

import (
	"fmt"
	"github.com/rs/zerolog"
	"os"
	"time"
)

type Logger struct {
	Lg *zerolog.Logger
}

func NewLogger(level string) (*Logger, error) {
	lvl, err := zerolog.ParseLevel(level)
	if err != nil {
		return nil, fmt.Errorf("failed to parse log level %v: %w", level, err)
	}
	logOut := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	var newLg zerolog.Logger
	if lvl == zerolog.DebugLevel {
		newLg = zerolog.New(logOut).Level(zerolog.DebugLevel).With().Timestamp().Logger()
	} else {
		newLg = zerolog.New(logOut).Level(zerolog.InfoLevel).With().Timestamp().Logger()
	}
	lg := Logger{}
	lg.Lg = &newLg
	return &lg, nil
}
