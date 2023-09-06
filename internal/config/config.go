package config

import (
	"github.com/kelseyhightower/envconfig"
	"os"
)

type Config struct {
	Proto      string `default:"tcp"`
	Listen     string `default:"127.0.0.1:1080"`
	LogLevel   string `default:"info" description:"log level: debug|info|warn|error|fatal"`
	AuthMethod string `default:"none" description:"auth method: none|static|ldap"`
	AuthStatic Static
	AuthLdap   Ldap
}

type Static struct {
	User string
	Pass string
}

type Ldap struct{}

func NewConfig(prefix string) (*Config, error) {
	var cfg Config
	err := envconfig.Process(prefix, &cfg)
	if err != nil {
		return &cfg, err
	}
	err = envconfig.CheckDisallowed(prefix, &cfg)
	if err != nil {
		return &cfg, err
	}
	return &cfg, nil
}

func PrintUsage(prefix string) {
	var cfg Config
	_ = envconfig.Usage(prefix, &cfg)
	os.Exit(128)
}
