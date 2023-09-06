package services

import (
	"fmt"
	"github.com/dossif/gosocks5/pkg/logger"
	"github.com/dossif/gosocks5/pkg/socks5"
	"golang.org/x/net/context"
)

type Server struct {
	Ctx    context.Context
	Lg     *logger.Logger
	Srv    socks5.Server
	Proto  string
	Listen string
}

func NewServer() {}

func (s *Server) Server() {
	// Create a SOCKS5 server
	conf := &socks5.Config{
		AuthMethods: nil,
		Credentials: nil,
		Resolver:    nil,
		Rules:       nil,
		Rewriter:    nil,
		BindIP:      nil,
		Logger:      s.Lg,
		Dial:        nil,
	}
	server, err := socks5.New(conf)
	if err != nil {
		panic(err)
	}
	fmt.Printf("server listen %s://%s\n", s.Proto, s.Listen)
	if err := server.ListenAndServe(s.Proto, s.Listen); err != nil {
		panic(err)
	}
}
