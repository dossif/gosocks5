package services

import (
	"fmt"
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

func NewService(ctx context.Context, lg *logger.Logger, listen string, auth string) (*Service, error) {
	conf := &socks5.Config{
		AuthMethods: nil, // TODO: implement
		Credentials: nil, // TODO: implement
		Resolver:    nil,
		Rules:       nil,
		Rewriter:    nil,
		BindIP:      nil,
		Logger:      lg,
		Dial:        nil,
	}
	srv, err := socks5.New(ctx, conf)
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
