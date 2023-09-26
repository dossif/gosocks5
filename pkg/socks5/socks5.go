package socks5

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/dossif/gosocks5/pkg/logger"
	"golang.org/x/net/context"
	"net"
)

const (
	socks5Version = uint8(5)
)

// Config is used to set up and configure a Server
type Config struct {
	// AuthMethods can be provided to implement custom authentication
	// By default, "auth-less" mode is enabled.
	// For password-based auth use UserPassAuthenticator.
	AuthMethods []Authenticator

	// If provided, username/password authentication is enabled,
	// by appending a UserPassAuthenticator to AuthMethods. If not provided,
	// and AUthMethods is nil, then "auth-less" mode is enabled.
	Credentials CredentialStore

	// Resolver can be provided to do custom name resolution.
	// Defaults to DNSResolver if not provided.
	Resolver NameResolver

	// Rules is provided to enable custom logic around permitting
	// various commands. If not provided, PermitAll is used.
	Rules RuleSet

	// Rewriter can be used to transparently rewrite addresses.
	// This is invoked before the RuleSet is invoked.
	// Defaults to NoRewrite.
	Rewriter AddressRewriter

	// BindIP is used for bind or udp associate
	BindIP net.IP

	// Logger can be used to provide a custom log target.
	// Defaults to stdout.
	Logger *logger.Logger

	// Optional function for dialing out
	Dial func(ctx context.Context, network, addr string) (net.Conn, error)
}

// Server is responsible for accepting connections and handling
// the details of the SOCKS5 protocol
type Server struct {
	ctx         context.Context
	config      *Config
	authMethods map[uint8]Authenticator
}

// New creates a new Server and potentially returns an error
func New(ctx context.Context, conf *Config) (*Server, error) {
	// Ensure we have at least one authentication method enabled
	if len(conf.AuthMethods) == 0 {
		if conf.Credentials != nil {
			conf.AuthMethods = []Authenticator{&UserPassAuthenticator{conf.Credentials}}
		} else {
			conf.AuthMethods = []Authenticator{&NoAuthAuthenticator{}}
		}
	}

	// Ensure we have a DNS resolver
	if conf.Resolver == nil {
		conf.Resolver = DNSResolver{}
	}

	// Ensure we have a rule set
	if conf.Rules == nil {
		conf.Rules = PermitAll()
	}

	server := &Server{
		config: conf,
		ctx:    ctx,
	}

	server.authMethods = make(map[uint8]Authenticator)

	for _, a := range conf.AuthMethods {
		server.authMethods[a.GetCode()] = a
	}

	return server, nil
}

// ListenAndServe is used to create a listener and serve on it
func (s *Server) ListenAndServe(network, addr string) error {
	var l net.ListenConfig
	ll, err := l.Listen(s.ctx, network, addr)
	if err != nil {
		return err
	} else {
		s.config.Logger.Lg.Info().Msgf("start new %v listener on %v", network, addr)
	}
	return s.Serve(ll)
}

// Serve is used to serve connections from a listener
func (s *Server) Serve(l net.Listener) error {
	go func() {
		<-s.ctx.Done()
		_ = l.Close()
	}()
	for {
		conn, err := l.Accept()
		if err != nil {
			select {
			case <-s.ctx.Done():
				s.config.Logger.Lg.Info().Msgf("close listener")
				return nil
			default:
				return fmt.Errorf("failed to accept connection: %v", err)
			}
		}
		go func() {
			err := s.ServeConn(conn)
			if err != nil {
				s.config.Logger.Lg.Warn().Msgf("failed to serve connection: %v", err)
			}
		}()
	}
}

// ServeConn is used to serve a single connection.
func (s *Server) ServeConn(conn net.Conn) error {
	defer func() { _ = conn.Close() }()
	bufConn := bufio.NewReader(conn)

	// Read the version byte
	version := []byte{0}
	if _, err := bufConn.Read(version); err != nil {
		return fmt.Errorf("failed to get version byte: %v", err)
	}

	// Ensure we are compatible
	if version[0] != socks5Version {
		return fmt.Errorf("unsupported socks version: %v", version)
	}

	// Authenticate the connection
	authContext, err := s.authenticate(conn, bufConn)
	if err != nil {
		return fmt.Errorf("failed to authenticate: %v", err)
	}

	request, err := NewRequest(bufConn)
	if err != nil {
		if errors.Is(err, unrecognizedAddrType) {
			if err := sendReply(conn, addrTypeNotSupported, nil); err != nil {
				return fmt.Errorf("failed to send reply: %v", err)
			}
		}
		return fmt.Errorf("failed to read destination address: %v", err)
	}
	request.AuthContext = authContext
	if client, ok := conn.RemoteAddr().(*net.TCPAddr); ok {
		request.RemoteAddr = &AddrSpec{IP: client.IP, Port: client.Port}
	}
	s.config.Logger.Lg.Debug().Msgf("%s -> %s", request.RemoteAddr, request.DestAddr)

	// Process the client request
	if err := s.handleRequest(request, conn); err != nil {
		return fmt.Errorf("failed to handle request: %v", err)
	}

	return nil
}
