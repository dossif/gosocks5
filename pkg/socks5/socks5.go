package socks5

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/dossif/gosocks5/pkg/logger"
	"github.com/google/uuid"
	"golang.org/x/net/context"
	"net"
	"runtime"
	"time"
)

const (
	socks5Version = uint8(5)
	connDeadline  = time.Second * 10
)

var ConnCount int

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

	// Optional function for dialing out
	Dial func(ctx context.Context, network, addr string) (net.Conn, error)
}

// Server is responsible for accepting connections and handling
// the details of the SOCKS5 protocol
type Server struct {
	id          uuid.UUID
	ctx         context.Context
	Lg          *logger.Logger
	config      *Config
	authMethods map[uint8]Authenticator
}

type Connection struct {
	id   uuid.UUID
	Lg   *logger.Logger
	conn net.Conn
}

// New creates a new Server and potentially returns an error
func New(ctx context.Context, lg *logger.Logger, conf *Config) (*Server, error) {
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
		id:     uuid.New(),
		config: conf,
		ctx:    ctx,
		Lg:     lg,
	}

	server.authMethods = make(map[uint8]Authenticator)

	for _, a := range conf.AuthMethods {
		server.authMethods[a.GetCode()] = a
	}

	return server, nil
}

// ListenAndServe is used to create a listener and serve on it
func (s *Server) ListenAndServe(network, addr string) error {
	l := net.ListenConfig{KeepAlive: time.Second * 1}
	ll, err := l.Listen(s.ctx, network, addr)
	if err != nil {
		return err
	} else {
		s.Lg.Lg.Info().Msgf("start new %v listener on %v", network, addr)
	}
	go func() {
		for {
			var m runtime.MemStats
			runtime.ReadMemStats(&m)
			s.Lg.Lg.Trace().
				Int("conn", ConnCount).
				Int("gorutines", runtime.NumGoroutine()).
				Uint64("memTotalAllocMiB", m.TotalAlloc/1024/1024).
				Uint64("memSysMiB", m.Sys/1024/1024).
				Uint32("memNumGc", m.NumGC).
				Msgf("statistics")
			time.Sleep(time.Second * 5)
		}
	}()
	return s.ServeListener(ll)
}

// ServeListener is used to serve connections from a listener
func (s *Server) ServeListener(l net.Listener) error {
	go func() {
		<-s.ctx.Done()
		_ = l.Close()
		s.Lg.Lg.Trace().Msgf("close listener")
	}()
	for {
		conn, err := l.Accept() // blocking
		if err != nil {
			select {
			case <-s.ctx.Done():
				s.Lg.Lg.Info().Msgf("close listener")
				return nil
			default:
				return fmt.Errorf("failed to accept connection: %v", err)
			}
		}
		connId := uuid.New()
		l := *s.Lg
		l.AddField(map[string]string{"connId": connId.String()})
		go func() {
			err := s.ServeConnection(Connection{
				id:   connId,
				Lg:   &l,
				conn: conn,
			})
			if err != nil {
				s.Lg.Lg.Warn().Msgf("failed to serve connection: %v", err)
			}
		}()
	}
}

// ServeConnection is used to serve a single connection.
func (s *Server) ServeConnection(conn Connection) error {
	ConnCount = ConnCount + 1
	conn.Lg.Lg.Trace().Msgf("start connection")
	defer func() {
		err := conn.conn.Close()
		if err != nil {
			conn.Lg.Lg.Warn().Msgf("failed to close connection %v", err)
		} else {
			ConnCount = ConnCount - 1
			conn.Lg.Lg.Trace().Msgf("close connection")
		}
	}()
	bufConn := bufio.NewReader(conn.conn)

	_ = conn.conn.SetReadDeadline(time.Now().Add(connDeadline))

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
	authContext, user, err := s.authenticate(conn.conn, bufConn)
	if err != nil {
		return fmt.Errorf("failed to authenticate: %v", err)
	}
	reqId := uuid.New()
	l := *conn.Lg
	l.AddField(map[string]string{"reqId": reqId.String(), "user": user})
	request, err := NewRequest(reqId, &l, bufConn)
	if err != nil {
		if errors.Is(err, unrecognizedAddrType) {
			if err := sendReply(conn.conn, addrTypeNotSupported, nil); err != nil {
				return fmt.Errorf("failed to send reply: %v", err)
			}
		}
		return fmt.Errorf("failed to read destination address: %v", err)
	}
	request.AuthContext = authContext
	if client, ok := conn.conn.RemoteAddr().(*net.TCPAddr); ok {
		request.RemoteAddr = &AddrSpec{IP: client.IP, Port: client.Port}
	}
	l.Lg.Debug().Msgf("%s -> %s", request.RemoteAddr, request.DestAddr)

	// Process the client request
	if err := s.handleRequest(request, conn.conn); err != nil {
		return fmt.Errorf("failed to handle request: %v", err)
	}

	return nil
}
