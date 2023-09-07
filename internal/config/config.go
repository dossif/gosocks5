package config

import (
	"github.com/kelseyhightower/envconfig"
	"os"
)

type Config struct {
	Listen     string `default:"0.0.0.0:1080" desc:"socks5 server listen ip:port"`
	LogLevel   string `default:"info" desc:"log level: debug|info|warn|error|fatal"`
	AuthMethod string `default:"none" desc:"auth method: none|static|ldap"`
	AuthStatic Static
	AuthLdap   Ldap
}

type Static struct {
	User string `desc:"static username"`
	Pass string `desc:"static password"`
}

type Ldap struct {
	Url      string `desc:"ldap url. example: ldaps://example.com:636"`
	BindUser string `desc:"ldap bind user. example: uid=bind,cn=users,cn=accounts,dc=example,dc=com"`
	BindPass string `desc:"ldap bind pass"`
	BaseDn   string `desc:"ldap search base dn. example: cn=users,cn=accounts,dc=example,DC=com"`
	Filter   string `desc:"ldap search filter. example: (&(uid=%s)(memberOf=cn=devops,cn=groups,cn=accounts,dc=example,dc=com))"`
}

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
