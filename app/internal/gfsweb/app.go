package gfsweb

import (
	"context"
	"io"
	"net/http"

	"github.com/AgencyPMG/go-from-scratch/app/internal/gfsweb/server"
)

//App is the central application type for the gfsweb executable.
type App struct {
	//server is our application http server.
	server *server.Server

	//handler is used to handle all incoming requests.
	handler http.Handler

	//dbCloser is a Closer that should be called to close our connection to the database.
	dbCloser io.Closer
}

//Run runs a within ctx.
//It calls Serce on a's internal server.Server.
//It also waits in a sepearate goroutine for ctx to be done and calls Shutdown
//on the internal server.
func (a *App) Run(ctx context.Context) error {
	go func() {
		defer a.server.Shutdown(context.Background())
		<-ctx.Done()
	}()

	return a.server.Serve(a.handler)
}

//Close closes all internal resources held be a.
//This must be called regardless of the result from a.Run and should be called
//after Run returns.
func (a *App) Close() error {
	return a.dbCloser.Close()
}
