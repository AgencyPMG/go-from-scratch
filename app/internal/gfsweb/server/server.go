package server

import (
	"context"
	"fmt"
	"net"
	"net/http"
)

//ListenerFactory provides the Listener method that is used to create a net.Listener.
type ListenerFactory interface {
	//Listener should return a new net.Listener that is open and available for use
	//or an error if there is a failure.
	Listener() (net.Listener, error)
}

//PortListenerFactory is a ListenerFactory that returns a Listener by listening
//via tcp4 on the loopback interface and the port represented by this value.
type PortListenerFactory int

//Listener returns a tcp listener listening on the loopback interface and the port
//f.
//
//Note that a value of 0 is useful for testing and allows the system to select
//an open port.
func (f PortListenerFactory) Listener() (net.Listener, error) {
	address := fmt.Sprintf("%s:%d", "127.0.0.1", f)

	tcpAddress, err := net.ResolveTCPAddr("tcp4", address)
	if err != nil {
		return nil, err
	}

	return net.ListenTCP("tcp4", tcpAddress)
}

//Server is a wrapper for the http.Server to use in the application.
type Server struct {
	//ListenerFactory is the promoted factory used to create a Listener.
	//Must not be nil.
	ListenerFactory

	server *http.Server
}

//Serve attempts to serve h on the Listener returned from s.ListenerFactory.
//If a Listener could not be created, then that error is returned immediately.
//Otherwise, an http.Server is created to use h and its Serve method is called
//with the new Listener.
//
//The new Listener is closed before this method returns.
//
//A non-nil error means serving has failed in some way.
//If s is stoped via a call to s.Shutdown, then the error returned should be nil
//if nothing went wrong during the shutdown process.
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

//Shutdown attempts to stop s.
func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
