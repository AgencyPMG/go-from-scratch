package gfsweb

import (
	"github.com/AgencyPMG/go-from-scratch/app/internal/data/client"
	"github.com/AgencyPMG/go-from-scratch/app/internal/data/client/clientsql"
	"github.com/AgencyPMG/go-from-scratch/app/internal/data/storage/sqlrepo"
	"github.com/AgencyPMG/go-from-scratch/app/internal/data/user"
	"github.com/AgencyPMG/go-from-scratch/app/internal/data/user/usersql"
)

//Repos is a collection of domain type repositories used throughout the application.
//Passing this collection around is easier than passing all of them individually.
type Repos struct {
	Users   user.Repo
	Clients client.Repo
}

//NewRepos returns a new Repos with each repository created from each domain
//package's sql repo implementation.
func NewRepos(sqlRepo *sqlrepo.Repo) *Repos {
	return &Repos{
		Users:   usersql.New(sqlRepo),
		Clients: clientsql.New(sqlRepo),
	}
}
