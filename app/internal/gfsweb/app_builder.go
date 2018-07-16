package gfsweb

import (
	"io"

	"github.com/AgencyPMG/go-from-scratch/app/internal/data/storage/sqlrepo"
	"github.com/AgencyPMG/go-from-scratch/app/internal/gfsweb/handler/api"
	"github.com/AgencyPMG/go-from-scratch/app/internal/gfsweb/server"
	"github.com/gogolfing/cbus"
	"github.com/gogolfing/config"
)

type SQLRepoFactory func(config *config.Config) (*sqlrepo.Repo, io.Closer, error)

type CBusFactory func(repos *Repos) *cbus.Bus

type APIFactory func(bus *cbus.Bus, repos *Repos) *api.API

type ServerFactory func(config *config.Config) *server.Server

type AppBuilder struct {
	SQLRepoFactory

	CBusFactory

	APIFactory

	ServerFactory
}

func NewAppBuilder() *AppBuilder {
	return &AppBuilder{
		SQLRepoFactory: CreateSQLRepo,
		CBusFactory:    CreateCBus,
		APIFactory:     CreateAPI,
		ServerFactory:  CreateServer,
	}
}

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
