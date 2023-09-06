package logger

import (
	"fmt"
	"github.com/rs/zerolog"
	"os"
	"time"
)

type LoggerInterface interface {
	AddField(fields map[string]string)
}

type Logger struct {
	Lg *zerolog.Logger
}

func NewLogger(level string) (*Logger, error) {
	lvl, err := zerolog.ParseLevel(level)
	if err != nil {
		return &Logger{}, fmt.Errorf("failed to parse log level %v: %v", level, err)
	}
	cwr := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	lg := zerolog.New(cwr).Level(lvl).With().Timestamp().Logger()
	return &Logger{Lg: &lg}, nil
}

func (l *Logger) AddField(fields map[string]string) {
	nl := *l.Lg
	for k, v := range fields {
		nl = l.Lg.With().Str(k, v).Logger()
		l.Lg = &nl
	}
}
