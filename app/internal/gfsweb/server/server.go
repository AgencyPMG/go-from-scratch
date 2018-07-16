package server

import (
	"context"
	"fmt"
	"net"
	"net/http"
)

type ListenerFactory interface {
	Listener() (net.Listener, error)
}

type PortListenerFactory int

func (f PortListenerFactory) Listener() (net.Listener, error) {
	address := fmt.Sprintf("%s:%d", "127.0.0.1", f)

	tcpAddress, err := net.ResolveTCPAddr("tcp4", address)
	if err != nil {
		return nil, err
	}

	return net.ListenTCP("tcp4", tcpAddress)
}

type Server struct {
	//ListenerFactory is the promoted factory used to create a Listener.
	ListenerFactory

	server *http.Server
}

func (s *Server) Serve(h http.Handler) error {
	lis, err := s.Listener()
	if err != nil {
		return err
	}

	defer lis.Close()

	s.server = &http.Server{
		Addr:    lis.Addr().String(),
		Handler: h,
	}

	defer func() {
		s.server = nil
	}()

	err = s.server.Serve(lis)
	if err == http.ErrServerClosed {
		err = nil
	}
	return err
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
