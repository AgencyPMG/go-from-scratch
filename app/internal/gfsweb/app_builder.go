package gfsweb

import (
	"io"

	"github.com/AgencyPMG/go-from-scratch/app/internal/data/storage/sqlrepo"
	"github.com/AgencyPMG/go-from-scratch/app/internal/gfsweb/handler/api"
	"github.com/AgencyPMG/go-from-scratch/app/internal/gfsweb/server"
	"github.com/gogolfing/cbus"
	"github.com/gogolfing/config"
)

//SQLRepoFactory is a function type that creates a new sqlrepo.Reop and a related
//io.Closer that should close the connection to the database.
type SQLRepoFactory func(config *config.Config) (*sqlrepo.Repo, io.Closer, error)

//CBusFactory is a function type that creates a new cbus.Bus will all required
//Commands and Handlers correctly registered.
type CBusFactory func(repos *Repos) *cbus.Bus

//APIFactory is a function type that creates a new API ready for use.
type APIFactory func(bus *cbus.Bus, repos *Repos) *api.API

//ServerFactory is a function type that creates a new server.Server ready for use.
type ServerFactory func(config *config.Config) *server.Server

//AppBuilder is a type that knows how to initialize and build everything required
//for an App.
type AppBuilder struct {
	SQLRepoFactory

	CBusFactory

	APIFactory

	ServerFactory
}

//NewAppBuilder returns a new AppBuilder with fields set to default values.
func NewAppBuilder() *AppBuilder {
	return &AppBuilder{
		SQLRepoFactory: CreateSQLRepo,
		CBusFactory:    CreateCBus,
		APIFactory:     CreateAPI,
		ServerFactory:  CreateServer,
	}
}

//Build uses config and its factory functions to build a new App.
//If err is nil, the returned App is ready to Run.
func (ab *AppBuilder) Build(config *config.Config) (app *App, err error) {
	var sqlRepo *sqlrepo.Repo
	var dbCloser io.Closer

	defer func() {
		if err != nil && dbCloser != nil {
			dbCloser.Close()
		}
	}()

	sqlRepo, dbCloser, err = ab.SQLRepoFactory(config)
	if err != nil {
		return nil, err
	}

	repos := NewRepos(sqlRepo)

	bus := ab.CBusFactory(repos)

	api := ab.APIFactory(bus, repos)

	server := ab.ServerFactory(config)

	return &App{
		server:   server,
		handler:  api.Handler(),
		dbCloser: dbCloser,
	}, nil
}
