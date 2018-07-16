package gfsweb

import (
	"github.com/AgencyPMG/go-from-scratch/app/internal/data/client"
	"github.com/AgencyPMG/go-from-scratch/app/internal/data/client/clientsql"
	"github.com/AgencyPMG/go-from-scratch/app/internal/data/storage/sqlrepo"
	"github.com/AgencyPMG/go-from-scratch/app/internal/data/user"
	"github.com/AgencyPMG/go-from-scratch/app/internal/data/user/usersql"
)

type Repos struct {
	Users   user.Repo
	Clients client.Repo
}

func NewRepos(sqlRepo *sqlrepo.Repo) *Repos {
	return &Repos{
		Users:   usersql.New(sqlRepo),
		Clients: clientsql.New(sqlRepo),
	}
}
