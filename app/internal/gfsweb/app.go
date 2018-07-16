package gfsweb

import (
	"context"
	"io"
	"net/http"

	"github.com/AgencyPMG/go-from-scratch/app/internal/gfsweb/server"
)

type App struct {
	//server is our application http server.
	server *server.Server

	//handler is used to handle all incoming requests.
	handler http.Handler

	//dbCloser is a Closer that should be called to close our connection to the database.
	dbCloser io.Closer
}

func (a *App) Run(ctx context.Context) error {
	go func() {
		<-ctx.Done()
		a.server.Shutdown(context.Background())
	}()

	return a.server.Serve(a.handler)
}

func (a *App) Close() error {
	return a.dbCloser.Close()
}
