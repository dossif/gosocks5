package services

import (
	"fmt"
	"github.com/dossif/gosocks5/internal/auth/ldap"
	"github.com/dossif/gosocks5/internal/auth/static"
	"github.com/dossif/gosocks5/internal/config"
	"github.com/dossif/gosocks5/pkg/logger"
	"github.com/dossif/gosocks5/pkg/socks5"
	"golang.org/x/net/context"
)

const (
	proto = "tcp"
)

type Service struct {
	Ctx    context.Context
	Lg     *logger.Logger
	Srv    socks5.Server
	Listen string
}

func NewService(ctx context.Context, lg *logger.Logger, listen string, auth config.Auth) (*Service, error) {
	var authMethod socks5.Authenticator
	switch auth.Method {
	case "none":
		authMethod = socks5.NoAuthAuthenticator{}
		lg.Lg.Info().Msgf("auth mode: none")
	case "static":
		st, err := static.NewStatic(auth.Static.User, auth.Static.Pass)
		if err != nil {
			return &Service{}, fmt.Errorf("failed to create static auth: %v", err)
		}
		authMethod = socks5.UserPassAuthenticator{Credentials: st}
		lg.Lg.Info().Msgf("auth mode: static")
	case "ldap":
		ld, err := ldap.NewLdap(*lg, auth.Ldap.Url, auth.Ldap.BindUser, auth.Ldap.BindPass, auth.Ldap.BaseDn, auth.Ldap.Filter)
		if err != nil {
			return &Service{}, fmt.Errorf("failed to create ldap auth: %v", err)
		}
		authMethod = socks5.UserPassAuthenticator{Credentials: ld}
		lg.Lg.Info().Msgf("auth mode: ldap")
	}
	conf := &socks5.Config{
		AuthMethods: []socks5.Authenticator{authMethod},
	}
	srv, err := socks5.New(ctx, lg, conf)
	if err != nil {
		return &Service{}, fmt.Errorf("failed to create socks5 server: %v", err)
	}
	return &Service{
		Ctx:    ctx,
		Lg:     lg,
		Srv:    *srv,
		Listen: listen,
	}, nil
}

func (s *Service) Start() error {
	err := s.Srv.ListenAndServe(proto, s.Listen)
	if err != nil {
		return fmt.Errorf("failed to start gosocks5 server: %v", err)
	}
	return nil
}
